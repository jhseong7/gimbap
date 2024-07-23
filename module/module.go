// File: module.go
//
// This file defines the module interface.
package module

import (
	"reflect"

	"github.com/jhseong7/ecl"
	"github.com/jhseong7/gimbap/controller"
	"github.com/jhseong7/gimbap/provider"
	"github.com/jhseong7/gimbap/util"
)

type (
	Module struct {
		// Name of module
		Name string

		providerList []*provider.Provider

		providerMapWithHandler map[string]map[ProviderKey]interface{}
	}

	ModuleOption struct {
		// Name of the module
		Name string

		// List of modules that this module depends on.
		SubModules []*Module

		// List of providers that this module provides.
		Providers []*provider.Provider

		// Rest controllers
		Controllers []*controller.Controller
	}

	// Key to uniquely identify a provider in a module
	ProviderKey struct {
		Type reflect.Type
		Name string
	}
)

// Logger for the module initialization
var log = ecl.NewLogger(ecl.LoggerOption{Name: "DefineModule", AppName: "GIMBAP"})

// Extract the embedded provider from the interface. If the interface is a provider, return the provider.
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

// Get the key of the provider.
func getKeyFromProvider(p provider.Provider) ProviderKey {
	// Get the return type of the Instantiator
	pType := reflect.TypeOf(p.Instantiator).Out(0)

	return ProviderKey{
		Type: pType,
		Name: util.GetFullNameOfType(pType),
	}
}

// Define a module and return a module struct.
func DefineModule(option ModuleOption) *Module {
	if option.Name == "" {
		log.Panicf("Module name cannot be empty")
	}

	providerList := []*provider.Provider{}
	providerMapWithHandler := map[string]map[ProviderKey]interface{}{}

	// For all the Submodules
	for _, m := range option.SubModules {
		// For all values of the provider with the handler
		for handlerName, providerMap := range m.providerMapWithHandler {
			if _, ok := providerMapWithHandler[handlerName]; !ok {
				providerMapWithHandler[handlerName] = map[ProviderKey]interface{}{}
			}

			for _, p := range providerMap {
				casted, ok := extractEmbeddedProvider(p)
				if !ok {
					log.Panicf("Provider %v is not a provider", p)
				}

				// Get the key of the provider
				pKey := getKeyFromProvider(*casted)

				// If the provider is already defined in the handler --> show warning, then skip
				if _, ok := providerMapWithHandler[handlerName][pKey]; ok {
					log.Warnf("Duplicate Provider warning: %s from module %s is already defined in handler [%s]. Skipping", pKey.Name, m.Name, handlerName)
					continue
				}

				// Add to the list
				providerMapWithHandler[handlerName][pKey] = p
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
			providerMapWithHandler[p.Handler] = map[ProviderKey]interface{}{}
		}

		pKey := getKeyFromProvider(*p)

		// If the provider is already defined in the handler --> show warning, then skip
		if _, ok := providerMapWithHandler[p.Handler][pKey]; ok {
			log.Warnf("Duplicate Provider warning: %s is already defined in handler [%s]. Skipping", pKey.Name, p.Handler)
			continue
		}

		// Add to the list
		providerMapWithHandler[p.Handler][pKey] = p
		providerList = append(providerList, p)
	}

	for _, c := range option.Controllers {
		// If the p.Handler is controller --> show warning
		if c.Handler != "controller" {
			log.Panicf("Controller %s is not a controller. Only controllers must be given to the controllers option", c.Name)
		}

		// Add to the "default" handler
		if _, ok := providerMapWithHandler[c.Handler]; !ok {
			providerMapWithHandler[c.Handler] = map[ProviderKey]interface{}{}
		}

		pKey := getKeyFromProvider(c.Provider)

		// If the controller is already defined in the handler --> show warning, then skip
		if _, ok := providerMapWithHandler[c.Handler][pKey]; ok {
			log.Warnf("Duplicate Provider warning: %s is already defined in handler [%s]. Skipping", pKey.Name, c.Handler)
			continue
		}

		// Add to the list
		providerMapWithHandler[c.Handler][pKey] = c
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

func (m *Module) GetProviderMapOfHandler(handler string) map[ProviderKey]interface{} {
	return m.providerMapWithHandler[handler]
}
