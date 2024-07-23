// Alias public types and functions for the package gimbap.
package gimbap

import (
	"github.com/jhseong7/gimbap/app"
	"github.com/jhseong7/gimbap/controller"
	"github.com/jhseong7/gimbap/engine"
	"github.com/jhseong7/gimbap/microservice"
	"github.com/jhseong7/gimbap/module"
	"github.com/jhseong7/gimbap/provider"
)

// type aliases for public apis
type (
	// App related
	AppOption      = app.AppOption
	GimbapApp      = app.GimbapApp
	RuntimeOptions = app.RuntimeOptions

	// Module related
	ModuleOption = module.ModuleOption
	Module       = module.Module

	// Provider related
	Provider       = provider.Provider
	ProviderOption = provider.ProviderOption

	// Controller related
	IController      = controller.IController
	ControllerOption = controller.ControllerOption
	Controller       = controller.Controller
	RouteSpec        = controller.RouteSpec

	// Engine related
	IServerEngine      = engine.IServerEngine
	ServerEngineOption = engine.ServerEngineOption

	// Microservice related
	IMicroService              = microservice.IMicroService
	MicroServiceProvider       = microservice.MicroServiceProvider
	MicroServiceProviderOption = microservice.MicroServiceProviderOption
)

// Create a Gimbap instance.
//
// This is the entry point to create a Gimbap application.
func CreateApp(option AppOption) *GimbapApp {
	return app.CreateApp(option)
}

// Function to get a provider from the app.
//
// Provide the app and the provider type to get the provider instance.
// If the provider is not found, it will panic.
func GetProvider[T interface{}](a GimbapApp, prov T) (ret T) {
	return app.GetProvider(a, prov)
}

// Define a module.
//
// This defines a module with the given option.
// The module is used to determine the dependencies of the providers.
func DefineModule(option ModuleOption) *Module {
	return module.DefineModule(option)
}

// Define a provider.
//
// This defines a provider with the given option.
// The provider will be registered to the app and can be injected to the controllers.
func DefineProvider(option ProviderOption) *Provider {
	return provider.DefineProvider(option)
}

// Define a controller.
//
// Defines a special provider that is used to handle RESTful requests.
// The controller will be registered to the app and can be injected to the other controllers.
func DefineController(option ControllerOption) *Controller {
	return controller.DefineController(option)
}

// Define a microservice
//
// Define a special provider for microservices
func DefineMicroService(option MicroServiceProviderOption) *MicroServiceProvider {
	return microservice.DefineMicroService(option)
}
