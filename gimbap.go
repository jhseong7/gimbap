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

	// AppOption
	//
	// This is the option provided when creating a Gimbap App
	//
	// Original type:
	// type AppOption struct {
	// 	AppName      string
	// 	AppModule    *module.Module
	// 	ServerEngine engine.IServerEngine
	// 	DepManager   dependency.IDependencyManager
	// }
	AppOption = app.AppOption

	// Original type:
	// type GimbapApp struct {
	// 	appName      string
	// 	appModule    *module.Module
	// 	serverEngine engine.IServerEngine
	// 	depManager   dependency.IDependencyManager
	// }
	GimbapApp = app.GimbapApp

	// Original type:
	// type RuntimeOptions struct {
	// 	Port int
	// 	TLS  *TLSOption
	// }
	RuntimeOptions = app.RuntimeOptions

	// Original type:
	// type TLSOption struct {
	// 	CertFile string
	// 	KeyFile  string
	// }
	TLSOption = engine.TLSOption

	// Module related
	// Original type:
	// type ModuleOption struct {
	// 	Name       string
	// 	Providers  []*Provider
	// 	DependsOn  []*Module
	// }
	ModuleOption = module.ModuleOption

	// Original type:
	// type Module struct {
	// 	name      string
	// 	providers []*Provider
	// 	dependsOn []*Module
	// }
	Module = module.Module

	// Provider related
	// Original type:
	// type Provider struct {
	// 	name     string
	// 	instance interface{}
	// 	module   *Module
	// }
	Provider = provider.Provider

	// Original type:
	// type ProviderOption struct {
	// 	Name     string
	// 	Instance interface{}
	// 	Module   *Module
	// }
	ProviderOption = provider.ProviderOption

	// Controller related
	// Original type:
	// type IController interface {
	// 	RegisterRoutes(router *gin.Engine)
	// }
	IController = controller.IController

	// Original type:
	// type ControllerOption struct {
	// 	Name     string
	// 	Instance IController
	// 	Module   *Module
	// }
	ControllerOption = controller.ControllerOption

	// Original type:
	// type Controller struct {
	// 	name     string
	// 	instance IController
	// 	module   *Module
	// }
	Controller = controller.Controller

	// Original type:
	// type RouteSpec struct {
	// 	Method  string
	// 	Path    string
	// 	Handler gin.HandlerFunc
	// }
	RouteSpec = controller.RouteSpec

	// Engine related
	// Original type:
	// type IServerEngine interface {
	// 	Start(port int, tls *TLSOption) error
	// 	Stop() error
	// }
	IServerEngine = engine.IServerEngine

	// Original type:
	// type ServerEngineOption struct {
	// 	EngineType string
	// }
	ServerEngineOption = engine.ServerEngineOption

	// Microservice related
	// Original type:
	// type IMicroService interface {
	// 	Start() error
	// 	Stop() error
	// }
	IMicroService = microservice.IMicroService

	// Original type:
	// type MicroServiceProvider struct {
	// 	name     string
	// 	instance IMicroService
	// 	module   *Module
	// }
	MicroServiceProvider = microservice.MicroServiceProvider

	// Original type:
	// type MicroServiceProviderOption struct {
	// 	Name     string
	// 	Instance IMicroService
	// 	Module   *Module
	// }
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
