package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"reflect"
	"syscall"

	logger "github.com/jhseong7/ecl"
	"github.com/jhseong7/gimbap/controller"
	"github.com/jhseong7/gimbap/engine"
	"github.com/jhseong7/gimbap/microservice"
	"github.com/jhseong7/gimbap/module"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

const (
	DEFAULT_PORT = 8080
)

type (
	GimbapApp struct {
		appModule module.Module      // Root module for the app.
		appOption AppOption          // Options for the app.
		engine    engine.IHttpEngine // The http engine that will handle RESTful requests.

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

		// microservice list
		microservices []microservice.MicroServiceProvider
	}

	AppOption struct {
		AppName    string
		AppModule  *module.Module
		HttpEngine engine.IHttpEngine
	}

	RuntimeOptions struct {
		Port int

		// Option injector with provided values from the app module
		WithProvided interface{}
	}
)

// Function to get a provider from the app.
//
// Provide the app and the provider type to get the provider instance.
// If the provider is not found, it will panic.
func GetProvider[T interface{}](app GimbapApp, prov T) (ret T) {
	defer func() {
		if r := recover(); r != nil {
			app.logger.Panicf("provider not found: %s", reflect.TypeOf(prov).String())
		}
	}()

	ret = app.instanceMap[reflect.TypeOf(prov)].Interface().(T)

	return
}

// Prepares the provider and injection functions for the app.
// Fx. does not initialize any providers unless they are explicitly called through fx.Invoke
// This method creates a list of providers to use in the fx.Invoke function from the given root Module
func (app *GimbapApp) createInjectionInits() (provider, initInvoker fx.Option) {
	// List to save all providers.
	var opList []fx.Option = []fx.Option{}

	// List to save return types of all providers.
	returnTypeList := []reflect.Type{}

	for _, p := range app.appModule.GetProviderMap() {
		// Add the instantiator to the optionList
		opList = append(opList, fx.Provide(p.Instantiator))

		funcType := reflect.TypeOf(p.Instantiator)

		// Get the output types and save them to returnTypeList
		for i := 0; i < funcType.NumOut(); i++ {
			returnTypeList = append(returnTypeList, funcType.Out(i))
		}
	}

	// Process controller maps
	for _, c := range app.appModule.GetControllerMap() {
		// Add the instantiator to the optionList
		opList = append(opList, fx.Provide(c.Instantiator))

		funcType := reflect.TypeOf(c.Instantiator)

		// Get the output types and save them to returnTypeList
		for i := 0; i < funcType.NumOut(); i++ {
			returnTypeList = append(returnTypeList, funcType.Out(i))
		}
	}

	// Process Microservices
	for _, m := range app.microservices {
		// Add the instantiator to the optionList
		opList = append(opList, fx.Provide(m.Instantiator))

		funcType := reflect.TypeOf(m.Instantiator)

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

// Create a Gimbap instance.
//
// This is the entry point to create a Gimbap application.
func CreateApp(option AppOption) *GimbapApp {
	// Setup global logger name
	logger.SetAppName(option.AppName)

	l := logger.NewLogger(logger.LoggerOption{
		Name: "Gimbap",
	})

	if option.AppModule == nil {
		l.Panic("AppModule is not set")
	}

	var e engine.IHttpEngine
	if option.HttpEngine == nil {
		l.Warn("HttpEngine is not set. Using default engine: GinHttpEngine")
		e = engine.NewGinHttpEngine()
	} else {
		e = option.HttpEngine
	}

	a := &GimbapApp{
		appModule: *option.AppModule,
		appOption: option,

		engine: e,

		instanceMap: make(map[reflect.Type]reflect.Value),

		logger: l,

		onStartListeners: []func(){},
		onStopListeners:  []func(){},
	}

	return a
}

func (app *GimbapApp) AddStartListener(listener func()) {
	app.onStartListeners = append(app.onStartListeners, listener)
}

func (app *GimbapApp) AddStopListener(listener func()) {
	app.onStopListeners = append(app.onStopListeners, listener)
}

// Retrieve the type of the return value of the instantiator function.
// Also checks if the instantiator is a function with a single return value.
func (app *GimbapApp) deriveTypeFromInstantiator(instantiator interface{}) (reflect.Type, bool) {
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
func (app *GimbapApp) registerControllerInstances() {
	// For all controllers
	for _, c := range app.appModule.GetControllerMap() {
		// Get the return type of the instantiator (this will be the controller's type)
		instanceType, ok := app.deriveTypeFromInstantiator(c.Instantiator)
		if !ok {
			app.logger.Panicf("Failed to derive type from instantiator: %s", c.Instantiator)
		}

		// Get the instance from the instance map
		instVal, ok := app.instanceMap[instanceType]
		if !ok {
			app.logger.Panicf("Controller instance not found in instance map: %s", instanceType.String())
		}

		// Bind the controller instance to the controller
		inst, ok := instVal.Interface().(controller.IController)
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
func (app *GimbapApp) run(lc fx.Lifecycle, runtimeOpts RuntimeOptions) {
	// Initialize the engine
	// Bind the controller instances to the engine and register the routes.
	// This will automatically call the GetRouteSpecs function of each controller. (if it is implemented)
	app.registerControllerInstances()

	// Register a SIGTEM, SIGINT listener to stop the app gracefully.
	// This will trigger the engine to stop --> calling an end to the app's lifecycle.
	// This will also call all the onStopListeners.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		app.logger.Logf("Received signal: %s", sig)

		// Stop the main engine to break the loop
		app.engine.Stop()
	}()

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			// NOTE: change this to go routine if there is a risk for deadlock.
			for _, listener := range app.onStartListeners {
				listener()
			}

			// Start the microservices if exists
			if len(app.microservices) > 0 {
				app.startMicroServices()
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

			// Stop the microservices if exists
			if len(app.microservices) > 0 {
				app.stopMicroServices()
			}

			app.logger.Log("App stopped")

			return nil
		},
	})
}

// Set a custom logger for the app.
//
// This will override the default logger.
func (app *GimbapApp) SetCustomLogger(logger logger.Logger) {
	app.logger = logger
}

// UseInjection is a function to add functions that will be called with the injection support.
//
// This is useful for initializing functions that need to use providers.
func (app *GimbapApp) UseInjection(functions ...interface{}) {
	// Initialize the functions list if it is nil.
	if app.functionsWithInjection == nil {
		app.functionsWithInjection = []interface{}{}
	}

	app.functionsWithInjection = append(app.functionsWithInjection, functions...)
}

// Add middleware to the engine.
//
// This will be added as a global middleware to the engine.
func (app *GimbapApp) AddMiddleware(middleware ...interface{}) {
	// Check if the engine is set.
	if middleware == nil {
		app.logger.Warn("Middleware is nil. Skipping")
		return
	}

	if app.engine == nil {
		app.logger.Panic("HttpEngine is not set. Cannot add middleware")
	}

	app.engine.AddMiddleware(middleware...)
}

// Add a microservice to the app.
//
// This will add a microservice to the app.
// The microservice will start with the app.
// Multiple microservices can be added.
func (app *GimbapApp) AddMicroServices(microservices ...microservice.MicroServiceProvider) {
	if microservices == nil {
		app.logger.Warn("(AddMicroServices) At least 1 microservice must be added to use this API. Skipping.")
		return
	}

	app.logger.Logf("Adding %d microservices", len(microservices))
	for _, m := range microservices {
		app.logger.Logf("Registering microservice: %s", m.Name)

		// Register the microservice as a provider and add it to the microservice list.
		app.microservices = append(app.microservices, m)
	}
}

// Check if the microservices are a valid type (implements IMicroService),
// then start the microservices using go routines. (non-blocking)
func (app *GimbapApp) startMicroServices() {
	for _, m := range app.microservices {
		// Get the return type
		instanceType, ok := app.deriveTypeFromInstantiator(m.Instantiator)
		if !ok {
			app.logger.Panicf("Failed to derive type from instantiator: %s", m.Instantiator)
		}

		// Get the instance from the instance map
		instVal, ok := app.instanceMap[instanceType]
		if !ok {
			app.logger.Panicf("Microservice instance not found in instance map: %s", instanceType.String())
		}

		// Check if the instance implements IMicroService
		if !microservice.IsMicroService(instVal.Interface()) {
			app.logger.Panicf("Microservice instance does not implement IMicroService: %s", instanceType.String())
		}

		// Start the microservice
		go instVal.Interface().(microservice.IMicroService).Start()

		app.logger.Logf("Microservice started: %s", m.Name)
	}
}

// Stop the microservices gracefully.
func (app *GimbapApp) stopMicroServices() {
	app.logger.Log("Stopping microservices")

	for _, m := range app.microservices {
		// Get the return type
		instanceType, ok := app.deriveTypeFromInstantiator(m.Instantiator)
		if !ok {
			app.logger.Panicf("Failed to derive type from instantiator: %s", m.Instantiator)
		}

		// Get the instance from the instance map
		instVal, ok := app.instanceMap[instanceType]
		if !ok {
			app.logger.Panicf("Microservice instance not found in instance map: %s", instanceType.String())
		}

		// Check if the instance implements IMicroService
		if !microservice.IsMicroService(instVal.Interface()) {
			app.logger.Panicf("Microservice instance does not implement IMicroService: %s", instanceType.String())
		}

		// Stop the microservice (don't use a go routine as we want to stop the app gracefully)
		instVal.Interface().(microservice.IMicroService).Stop()

		app.logger.Logf("Microservice stopped: %s", m.Name)
	}
}

// Stop the app
func (app *GimbapApp) Stop() {
	stopCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := app.fxApp.Stop(stopCtx); err != nil {
		log.Fatal(err)
	}
}

// The public start function that will start the app.
func (app *GimbapApp) Run(options ...RuntimeOptions) {
	// Catch any panic and log it.
	defer func() {
		if r := recover(); r != nil {
			app.logger.Fatalf("Failed to start the app. %s", r)
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

	stopCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := app.fxApp.Stop(stopCtx); err != nil {
		log.Fatal(err)
	}
}
