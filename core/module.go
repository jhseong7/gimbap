// File: module.go
//
// This file defines the module interface.
package core

import (
	logger "github.com/jhseong7/ecl"
	"github.com/jhseong7/gimbap/core/controller"
	"github.com/jhseong7/gimbap/core/provider"
)

type (
	Module struct {
		// Name of module
		Name string

		// Provider map of all providers in this module (that is exported)
		// key: provider name, value: provider.
		providerMap map[string]*provider.ProviderDefinition

		controllerMap map[string]*controller.ControllerDefinition
	}

	ModuleOption struct {
		// Name of the module
		Name string

		// List of modules that this module depends on.
		SubModules []Module

		// List of providers that this module provides.
		Providers []provider.ProviderDefinition

		// Rest controllers
		Controllers []controller.ControllerDefinition
	}
)

func DefineModule(option ModuleOption) *Module {
	if option.Name == "" {
		logger.NewLogger(logger.LoggerOption{Name: "DefineModule"}).Panicf("Controller name cannot be empty")
	}

	providerMap := map[string]*provider.ProviderDefinition{}
	controllerMap := map[string]*controller.ControllerDefinition{}

	// For all given imports and providers, create fx.Option
	for _, m := range option.SubModules {
		// Merge provider map and controller map
		for k, v := range m.providerMap {
			if _, ok := providerMap[k]; ok {
				panic("Provider name conflict: " + k + " is already defined.")
			}

			providerMap[k] = v
		}

		for k, v := range m.controllerMap {
			if _, ok := controllerMap[k]; ok {
				panic("Controller name conflict: " + k + " is already defined.")
			}

			controllerMap[k] = v
		}
	}

	// Handle providers
	for _, p := range option.Providers {
		if _, ok := providerMap[p.Name]; ok {
			panic("Provider name conflict: " + p.Name + " is already defined.")
		}

		providerMap[p.Name] = &p
	}

	// Handle controllers
	for _, c := range option.Controllers {
		// Add to the controller map regardless of whether it is exported or not.
		if _, ok := controllerMap[c.Name]; ok {
			panic("Controller name conflict: " + c.Name + " is already defined.")
		}

		controllerMap[c.Name] = &c
	}

	// Return the module with fx.Option
	return &Module{
		Name:          option.Name,
		providerMap:   providerMap,
		controllerMap: controllerMap,
	}
}
