package core

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/jhseong7/nassi-golang/logger"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

const (
	DEFAULT_PORT = 8080
)

type (
	NassiApp struct {
		appModule Module      // Root module for the app.
		appOption AppOption   // Options for the app.
		engine    IHttpEngine // The http engine that will handle RESTful requests.

		fxApp       *fx.App                        // fx.App instance for DI
		instanceMap map[reflect.Type]reflect.Value // Map to save instances of providers.

		logger logger.Logger

		// Lifecycle listeners
		onStartListeners []func()
		onStopListeners  []func()

		// TODO: Guards and Pipes
		// guards []interface{}
		// pipes  []interface{}

		// Function to run with the injection support
		functionsWithInjection []interface{}
	}

	AppOption struct {
		AppName   string
		AppModule *Module
	}

	RuntimeOptions struct {
		Port int

		// Option injector with provided values from the app module
		WithProvided interface{}
	}
)

// Function to get a provider from the app.
func GetProvider[T interface{}](app NassiApp, prov T) (ret T, err error) {
	defer func() {
		if r := recover(); r != nil {
			app.logger.Panicf("provider not found: %s", reflect.TypeOf(prov).String())
		}
	}()

	ret = app.instanceMap[reflect.TypeOf(prov)].Interface().(T)

	return
}

// Prepares the provider and injection functions for the app.
func (app *NassiApp) createInjectionInits() (provider, initInvoker fx.Option) {
	// List to save all providers.
	var opList []fx.Option = []fx.Option{}

	// List to save return types of all providers.
	returnTypeList := []reflect.Type{}

	for _, p := range app.appModule.providerMap {
		// Add the instantiator to the optionList
		opList = append(opList, fx.Provide(p.instantiator))

		funcType := reflect.TypeOf(p.instantiator)

		// Get the output types and save them to returnTypeList
		for i := 0; i < funcType.NumOut(); i++ {
			returnTypeList = append(returnTypeList, funcType.Out(i))
		}
	}

	// Process controller maps
	for _, c := range app.appModule.controllerMap {
		// Add the instantiator to the optionList
		opList = append(opList, fx.Provide(c.instantiator))

		funcType := reflect.TypeOf(c.instantiator)

		// Get the output types and save them to returnTypeList
		for i := 0; i < funcType.NumOut(); i++ {
			returnTypeList = append(returnTypeList, funcType.Out(i))
		}
	}

	// Create function type. Inputs --> all provider outputs (instantiators), returns --> nothing.
	functionType := reflect.FuncOf(returnTypeList, nil, false)

	// Create function to initialize all providers. at runtime
	function := reflect.MakeFunc(functionType, func(args []reflect.Value) []reflect.Value {
		for _, a := range args {
			// Save the instance to instanceMap
			app.instanceMap[reflect.TypeOf(a.Interface())] = a
			app.logger.Debugf("Provider instance created: %s", reflect.TypeOf(a.Interface()).String())
		}

		return nil
	})

	// Create a fx.Module as a provider for the instantiators
	provider = fx.Module("AppModule", opList...)

	// Return the function and the fx.Option
	initInvoker = fx.Invoke(function.Interface())

	return
}

// Create a NassiApp instance.
func CreateApp(option AppOption) *NassiApp {
	// Setup global logger name
	logger.SetAppName(option.AppName)

	// Create a new Gin Engine
	engine := NewEchoHttpEngine(HttpEngineOption{})

	a := &NassiApp{
		appModule: *option.AppModule,
		appOption: option,

		engine: engine,

		instanceMap: make(map[reflect.Type]reflect.Value),

		logger: logger.NewLogger(logger.LoggerOption{
			Name: "NassiApp",
		}),

		onStartListeners: []func(){},
		onStopListeners:  []func(){},
	}

	return a
}

func (app *NassiApp) AddStartListener(listener func()) {
	app.onStartListeners = append(app.onStartListeners, listener)
}

func (app *NassiApp) AddStopListener(listener func()) {
	app.onStopListeners = append(app.onStopListeners, listener)
}

// Retrieve the type of the return value of the instantiator function.
// Also checks if the instantiator is a function with a single return value.
func (app *NassiApp) deriveTypeFromInstantiator(instantiator interface{}) (reflect.Type, bool) {
	funcType := reflect.TypeOf(instantiator)

	if funcType.Kind() != reflect.Func {
		app.logger.Panicf("Instantiator is not a function: %s", funcType.String())
		return nil, false
	}

	// Check if the return type exists
	if funcType.NumOut() == 0 || funcType.NumOut() > 1 {
		app.logger.Panicf("Instantiator must have a single return value: %s", funcType.String())
		return nil, false
	}

	return funcType.Out(0), true
}

// Register the controller instances to the engine.
func (app *NassiApp) registerControllerInstances() {
	// For all controllers
	for _, c := range app.appModule.controllerMap {
		// Get the return type of the instantiator (this will be the controller's type)
		instanceType, ok := app.deriveTypeFromInstantiator(c.instantiator)
		if !ok {
			app.logger.Panicf("Failed to derive type from instantiator: %s", c.instantiator)
		}

		// Get the instance from the instance map
		instVal, ok := app.instanceMap[instanceType]
		if !ok {
			app.logger.Panicf("Controller instance not found in instance map: %s", instanceType.String())
		}

		// Bind the controller instance to the controller
		inst, ok := instVal.Interface().(IController)
		if !ok {
			app.logger.Panicf("Controller instance does not implement IController: %s", instanceType.String())
		}

		// Register the controller
		app.engine.RegisterController(c.RootPath, inst)
	}
}

// The internal run function that will be called by fx.Invoke.
//
// This function will start the engine and call all the onStartListeners.
func (app *NassiApp) run(lc fx.Lifecycle, runtimeOpts RuntimeOptions) {
	// Initialize the engine
	// Bind the controller instances to the engine and register the routes.
	// This will automatically call the GetRouteSpecs function of each controller. (if it is implemented)
	app.registerControllerInstances()

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			// NOTE: change this to go routine if there is a risk for deadlock.
			for _, listener := range app.onStartListeners {
				listener()
			}

			// Start the engine
			app.logger.Logf("App started on port: %d", runtimeOpts.Port)
			app.engine.Run(runtimeOpts.Port)

			return nil
		},
		OnStop: func(context.Context) error {
			// NOTE: change this to go routine if there is a risk for deadlock.
			for _, listener := range app.onStopListeners {
				listener()
			}

			app.logger.Log("App stopped")

			return nil
		},
	})
}

// Set a custom logger for the app.
//
// This will override the default logger.
func (app *NassiApp) SetCustomLogger(logger logger.Logger) {
	app.logger = logger
}

// UseInjection is a function to add functions that will be called with the injection support.
//
// This is useful for initializing functions that need to use providers.
func (app *NassiApp) UseInjection(functions ...interface{}) {
	if app.functionsWithInjection == nil {
		app.functionsWithInjection = []interface{}{}
	}

	app.functionsWithInjection = append(app.functionsWithInjection, functions...)
}

// Add middleware to the engine.
func (app *NassiApp) AddMiddleware(middleware ...interface{}) {
	app.engine.AddMiddleware(middleware...)
}

// The public start function that will start the app.
func (app *NassiApp) Run(options ...RuntimeOptions) {
	// Catch any panic and log it.
	defer func() {
		if r := recover(); r != nil {
			app.logger.Fatal("Failed to start the app")
		}
	}()

	var option RuntimeOptions
	if len(options) > 0 {
		option = options[0]
	} else {
		option = RuntimeOptions{Port: DEFAULT_PORT}
	}

	// Runtime option provider function
	var optionProvider fx.Option
	if option.WithProvided != nil {
		// TODO: Add a check for the input types of provided and see if it is in our provider map.
		optionProvider = fx.Provide(option.WithProvided)
	} else {
		optionProvider = fx.Provide(func() RuntimeOptions { return option })
	}

	// Create the injection function and the fx.Option
	providers, initInvoker := app.createInjectionInits()

	app.fxApp = fx.New(
		// Logger settings (TODO: make this configurable)
		fx.WithLogger(func() fxevent.Logger { return fxevent.NopLogger }),

		// PROVIDERS
		providers,      // Get auto generated fx.Option from the module
		optionProvider, // Inject the runtime options

		// INVOKEs
		initInvoker,

		// Run any custom run function for injections
		fx.Invoke(app.functionsWithInjection...),
		fx.Invoke(app.run),
	)

	startCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := app.fxApp.Start(startCtx); err != nil {
		log.Fatal(err)
	}

	if _, err := http.Get(fmt.Sprintf("http://localhost:%d/", option.Port)); err != nil {
		log.Fatal(err)
	}

	stopCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := app.fxApp.Stop(stopCtx); err != nil {
		log.Fatal(err)
	}
}