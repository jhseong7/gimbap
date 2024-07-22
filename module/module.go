// File: module.go
//
// This file defines the module interface.
package module

import (
	"reflect"

	logger "github.com/jhseong7/ecl"
	"github.com/jhseong7/gimbap/controller"
	"github.com/jhseong7/gimbap/provider"
)

type (
	Module struct {
		// Name of module
		Name string

		providerList []*provider.Provider

		providerMapWithHandler map[string]map[string]interface{}
	}

	ModuleOption struct {
		// Name of the module
		Name string

		// List of modules that this module depends on.
		SubModules []Module

		// List of providers that this module provides.
		Providers []provider.Provider

		// Rest controllers
		Controllers []controller.Controller
	}
)

func extractEmbeddedProvider(p interface{}) (*provider.Provider, bool) {
	// Return the provider if it is a provider
	if p, ok := p.(*provider.Provider); ok {
		return p, true
	}

	v := reflect.ValueOf(p)

	if v.Kind() != reflect.Ptr {
		return nil, false
	}

	// Dereference the pointer
	v = v.Elem()

	// For the fields of the struct, check if it is a provider
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := v.Type().Field(i)

		// If the type is an embedded struct --> check if it is a provider
		if fieldType.Anonymous {
			if field.Type() == reflect.TypeOf(provider.Provider{}) {
				prov := field.Interface().(provider.Provider)
				return &prov, true
			}
		}
	}

	return nil, false
}

func DefineModule(option ModuleOption) *Module {
	log := logger.NewLogger(logger.LoggerOption{Name: "DefineModule"})

	if option.Name == "" {
		log.Panicf("Module name cannot be empty")
	}

	providerList := []*provider.Provider{}
	providerMapWithHandler := map[string]map[string]interface{}{}

	// For all the Submodules
	for _, m := range option.SubModules {
		// For all values of the provider with the handler
		for handlerName, providerMap := range m.providerMapWithHandler {
			if _, ok := providerMapWithHandler[handlerName]; !ok {
				providerMapWithHandler[handlerName] = map[string]interface{}{}
			}

			for name, p := range providerMap {
				// If the name of the provider is already defined, panic.
				if _, ok := providerMapWithHandler[handlerName][name]; ok {
					log.Panicf("Provider name conflict: %s is already defined in handler %s", name, handlerName)
				}

				providerMapWithHandler[handlerName][name] = p

				casted, ok := extractEmbeddedProvider(p)

				if !ok {
					log.Panicf("Provider %s is not a provider", name)
				}

				providerList = append(providerList, casted)
			}
		}
	}

	// Handle providers
	for _, p := range option.Providers {
		// If the p.Handler is controller --> show warning
		if p.Handler == "controller" {
			log.Warnf("Provider %s is defined with handler 'controller'. Use 'Controller' options instead", p.Name)
		}

		if _, ok := providerMapWithHandler[p.Handler]; !ok {
			providerMapWithHandler[p.Handler] = map[string]interface{}{}
		}

		if _, ok := providerMapWithHandler[p.Handler][p.Name]; ok {
			log.Panicf("Provider name conflict: %s is already defined in handler %s", p.Name, "default")
		}

		providerMapWithHandler[p.Handler][p.Name] = &p
		providerList = append(providerList, &p)
	}

	for _, c := range option.Controllers {
		// If the p.Handler is controller --> show warning
		if c.Handler != "controller" {
			log.Warnf("Controller %s is not a controller. Please check the type", c.Name)
		}

		// Add to the "default" handler
		if _, ok := providerMapWithHandler[c.Handler]; !ok {
			providerMapWithHandler[c.Handler] = map[string]interface{}{}
		}

		if _, ok := providerMapWithHandler[c.Handler][c.Name]; ok {
			log.Panicf("Provider name conflict: %s is already defined in handler %s", c.Name, "default")
		}

		providerMapWithHandler[c.Handler][c.Name] = &c
		providerList = append(providerList, &c.Provider)
	}

	// Return the module with fx.Option
	return &Module{
		Name:                   option.Name,
		providerMapWithHandler: providerMapWithHandler,
		providerList:           providerList,
	}
}

func (m *Module) GetProviderList() []*provider.Provider {
	return m.providerList
}

func (m *Module) GetProviderMapOfHandler(handler string) map[string]interface{} {
	return m.providerMapWithHandler[handler]
}
