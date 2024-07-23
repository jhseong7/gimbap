package app

import (
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"github.com/jhseong7/ecl"

	"github.com/jhseong7/gimbap/controller"
	"github.com/jhseong7/gimbap/dependency"
	"github.com/jhseong7/gimbap/engine"
	"github.com/jhseong7/gimbap/microservice"
	"github.com/jhseong7/gimbap/module"
	"github.com/jhseong7/gimbap/provider"
	"github.com/jhseong7/gimbap/util"
)

const (
	DefaultPort int = 8080
)

type (
	GimbapApp struct {
		appModule    module.Module                 // Root module for the app.
		appOption    AppOption                     // Options for the app.
		serverEngine engine.IServerEngine          // The http engine that will handle RESTful requests.
		depManager   dependency.IDependencyManager // Engine that handles dependency injection.

		instanceMap map[reflect.Type]reflect.Value // Map to save instances of providers.

		logger ecl.Logger

		// Lifecycle listeners
		onStartListeners []func()
		onStopListeners  []func()

		// TODO: Guards and Pipes
		// guards []interface{}
		// pipes  []interface{}

		// Function to run with the injection support
		functionsWithInjection []*provider.Provider

		// microservice list
		microservices []*microservice.MicroServiceProvider

		// flag to hold the shutdown signal until all the components stop
		shutdownFlag chan string
		stopFlag     chan bool // Signal to trigger the stop of the app
	}

	AppOption struct {
		AppName      string
		AppModule    *module.Module
		ServerEngine engine.IServerEngine
		DepManager   dependency.IDependencyManager
	}

	RuntimeOptions struct {
		Port int

		// Option injector with provided values from the app module
		WithProvided interface{}
	}
)

const (
	MicroServiceMaxStartTime time.Duration = 5 * time.Second
	MicroServiceMaxStopTime  time.Duration = 5 * time.Second
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

// Create a Gimbap instance.
//
// This is the entry point to create a Gimbap application.
func CreateApp(option AppOption) *GimbapApp {
	// Setup global logger name
	ecl.SetAppName(option.AppName)

	l := ecl.NewLogger(ecl.LoggerOption{
		Name: "GIMBAP",
	})

	if option.AppModule == nil {
		l.Panic("AppModule is not set. Cannot create app.")
	}

	// Http engine
	var e engine.IServerEngine
	if option.ServerEngine == nil {
		l.Debug("HttpEngine is not set. Using default engine: GinHttpEngine")
		e = engine.NewGinHttpEngine()
	} else {
		e = option.ServerEngine
	}

	// Dependency manager
	var d dependency.IDependencyManager
	if option.DepManager == nil {
		l.Debug("DependencyManager is not set. Using default manager: FxManager")
		d = dependency.NewFxManager()
	} else {
		d = option.DepManager
	}

	a := &GimbapApp{
		appModule: *option.AppModule,
		appOption: option,

		serverEngine: e,
		depManager:   d,

		instanceMap: make(map[reflect.Type]reflect.Value),

		logger: l,

		onStartListeners: []func(){},
		onStopListeners:  []func(){},

		shutdownFlag: make(chan string),
		stopFlag:     make(chan bool),
	}

	return a
}

// Register the controller instances to the engine.
func (app *GimbapApp) registerControllerInstances() {
	// For all controllers
	for _, rc := range app.appModule.GetProviderMapOfHandler(controller.HandlerName) {
		c, ok := rc.(*controller.Controller)
		if !ok {
			app.logger.Panicf("Failed to cast controller: %s", reflect.TypeOf(rc).String())
		}

		// Get the return type of the instantiator (this will be the controller's type)
		instanceType, ok := util.DeriveTypeFromInstantiator(c.Instantiator)
		if !ok {
			app.logger.Panicf("Failed to derive type from instantiator: %v", c.Instantiator)
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
		app.serverEngine.RegisterController(c.RootPath, inst)
	}
}

// Internal function to get each active microservice, and handler with the given handler function
func (app *GimbapApp) forAllMicroservices(loopHandler func(microservice.IMicroService, *microservice.MicroServiceProvider)) {
	// Get the instances with the microservice handler
	for _, m := range app.microservices {
		// Get the instance from the instantiator
		instanceType, ok := util.DeriveTypeFromInstantiator(m.Instantiator)
		if !ok {
			app.logger.Panicf("Failed to derive type from instantiator. %v", m.Instantiator)
		}

		// Get the instance value from the instance map
		instVal, ok := app.instanceMap[instanceType]
		if !ok {
			app.logger.Panicf("Microservice instance not found in the instance map: %s", instanceType.String())
		}

		// Bind the microservice instance to the microservice interface
		inst, ok := instVal.Interface().(microservice.IMicroService)
		if !ok {
			app.logger.Panicf("Microservice instance does not implement the IMicroservice: %s", instanceType.String())
		}

		// Call start of the inst
		loopHandler(inst, m)
	}
}

// Check if the microservices are a valid type (implements IMicroService),
// then start the microservices using go routines. (non-blocking)
func (app *GimbapApp) startMicroServices() {
	app.forAllMicroservices(func(micro microservice.IMicroService, p *microservice.MicroServiceProvider) {
		app.logger.Logf("Starting microservice %s", p.Name)

		// Each micro service has 5 seconds to stop gracefully
		success := util.TimeoutJob(func() {
			micro.Start()
		},
			MicroServiceMaxStartTime,
		)

		if !success {
			app.logger.Warnf("Microservice %s failed to start on time. (within %s)", p.Name, MicroServiceMaxStartTime.String())
		}
	})
}

// Stop the microservices gracefully.
func (app *GimbapApp) stopMicroServices() {
	app.forAllMicroservices(func(micro microservice.IMicroService, p *microservice.MicroServiceProvider) {
		app.logger.Logf("Stopping microservice %s", p.Name)

		// Each micro service has 5 seconds to stop gracefully
		success := util.TimeoutJob(
			func() {
				micro.Stop()
			},
			MicroServiceMaxStopTime,
		)

		if !success {
			app.logger.Warnf("Microservice %s failed to stop on time. (within %s)", p.Name, MicroServiceMaxStopTime.String())
		}
	})
}

// Start routine other than the engine
func (app *GimbapApp) onStart() {
	app.logger.Log("Running on start routine")

	// NOTE: change this to go routine if there is a risk for deadlock.
	for _, listener := range app.onStartListeners {
		listener()
	}

	// Start the microservices if exists
	if len(app.microservices) > 0 {
		app.startMicroServices()
	}
}

// Stop routine other than the engine
func (app *GimbapApp) onStop() {
	app.logger.Log("Running on stop routine")

	// Call the onStopListeners
	for _, listener := range app.onStopListeners {
		listener()
	}

	// Stop the microservices if exists
	if len(app.microservices) > 0 {
		app.stopMicroServices()
	}
}

// The internal run function
//
// This function will start the engine and call all the onStartListeners.
func (app *GimbapApp) run() {
	// Get the runtime options from the instance map
	runtimeOpts := GetProvider(*app, RuntimeOptions{})

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
		select {
		case sig := <-sigs:
			app.logger.Logf("Received signal: %s", sig)
		case <-app.stopFlag:
			app.logger.Log("Received graceful stop request")
		}

		// Stop the main engine to break the loop
		app.serverEngine.Stop()

		app.onStop()

		// Send the shutdown signal
		app.shutdownFlag <- "shutdown"
	}()

	// Run the onStart lifecycle
	app.onStart()

	// NOTE, TODO: it may be a good idea to clear the arrays in the internal function for memory management.
	// if any will be added, do it here.

	// Start the engine
	app.logger.Log("App started")
	app.serverEngine.Run(runtimeOpts.Port) // Blocking from here

	// Defer function that blocks until the stop signal is received.
	defer func() {
		// Wait for the stop signal (blocking)
		<-app.shutdownFlag
		app.logger.Log("App has gracefully stopped")
	}()
}

// Add an onStart lifecycle listener to the app.
func (app *GimbapApp) AddOnStartListener(listener func()) {
	app.onStartListeners = append(app.onStartListeners, listener)
}

// Add an onStop lifecycle listener to the app.
func (app *GimbapApp) AddOnStopListener(listener func()) {
	app.onStopListeners = append(app.onStopListeners, listener)
}

// Set a custom logger for the app.
//
// This will override the default logger.
func (app *GimbapApp) SetCustomLogger(logger ecl.Logger) {
	app.logger = logger
}

// UseInjection is a function to add functions that will be called with the injection support.
//
// This is useful for initializing functions that need to use providers.
// A function that provides the value can be given, or the value itself can be given.
// If a function is given, that function can also benefit from the injection support.
func (app *GimbapApp) UseInjection(injectionValue interface{}) {
	// Initialize the functions list if it is nil.
	if app.functionsWithInjection == nil {
		app.functionsWithInjection = []*provider.Provider{}
	}

	var instantiator interface{}

	// Check if the function is a function type
	funcType := reflect.TypeOf(injectionValue)
	if funcType.Kind() == reflect.Func {
		instantiator = injectionValue
	} else {
		// Create a function that returns the function
		funcType := reflect.FuncOf([]reflect.Type{}, []reflect.Type{funcType}, false)
		instantiator = reflect.MakeFunc(funcType, func(args []reflect.Value) []reflect.Value {
			return []reflect.Value{reflect.ValueOf(injectionValue)}
		}).Interface()
	}

	app.functionsWithInjection = append(
		app.functionsWithInjection,
		provider.DefineProvider(provider.ProviderOption{
			Name:         "InjectionFunction",
			Instantiator: instantiator,
		}),
	)
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

	if app.serverEngine == nil {
		app.logger.Panic("HttpEngine is not set. Cannot add middleware")
	}

	app.serverEngine.AddMiddleware(middleware...)
}

// Add a microservice to the app.
//
// This will add a microservice to the app.
// The microservice will start with the app.
// Multiple microservices can be added.
func (app *GimbapApp) AddMicroServices(microservices ...*microservice.MicroServiceProvider) {
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

// Stop the app
func (app *GimbapApp) Stop() {
	// Send the stop signal
	app.stopFlag <- true
}

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
		option = RuntimeOptions{Port: DefaultPort}
	}

	// Runtime option provider function
	var optionProvider provider.Provider
	if option.WithProvided != nil {
		// TODO: Add a check for the input types of provided and see if it is in our provider map.
		optionProvider = *provider.DefineProvider(provider.ProviderOption{
			Name:         "RuntimeOptions",
			Instantiator: option.WithProvided,
		})
	} else {
		optionProvider = *provider.DefineProvider(provider.ProviderOption{
			Name:         "RuntimeOptions",
			Instantiator: func() RuntimeOptions { return option },
		})
	}

	// Call the dependency manager to inject the providers
	providers := []*provider.Provider{}

	// Collect all providers from the module
	for _, p := range app.appModule.GetProviderList() {
		providers = append(providers, p)
	}

	// Collect all microservices
	for _, m := range app.microservices {
		providers = append(providers, &m.Provider)
	}

	// Add the runtime options provider
	providers = append(providers, &optionProvider)

	// Add the functions with injection support
	if app.functionsWithInjection != nil {
		providers = append(providers, app.functionsWithInjection...)
	}

	// Inject the providers
	app.depManager.ResolveDependencies(app.instanceMap, providers)

	// Run the app (blocking from here)
	app.run()
}
