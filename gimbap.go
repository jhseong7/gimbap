package gimbap

import (
	"github.com/jhseong7/gimbap/app"
	"github.com/jhseong7/gimbap/controller"
	"github.com/jhseong7/gimbap/engine"
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
	ProviderDefinition = provider.ProviderDefinition

	// Controller related
	IController          = controller.IController
	ControllerDefinition = controller.ControllerDefinition
	RouteSpec            = controller.RouteSpec

	// Engine related
	IHttpEngine      = engine.IHttpEngine
	HttpEngineOption = engine.HttpEngineOption
)

// Public entry point to create a Gimbap application.
func CreateApp(option app.AppOption) *app.GimbapApp {
	return app.CreateApp(option)
}

// Public to retrive the injected provider.
func GetProvider[T interface{}](a app.GimbapApp, prov T) (ret T, err error) {
	return app.GetProvider(a, prov)
}
