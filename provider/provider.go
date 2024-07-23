// file: provider.go
//
// This file defines the provider interface and its implementation.
package provider

import "github.com/jhseong7/ecl"

type (
	Provider struct {
		Name         string
		Instantiator interface{}

		// The handler string is used to identify the handler in the provider. (e.g. Controller)
		Handler ProviderHandlerName
	}

	ProviderOption struct {
		Name         string
		Instantiator interface{}
	}

	ProviderHandlerName string
)

const (
	HandlerName ProviderHandlerName = "default"
)

// Define a provider
func DefineProvider(option ProviderOption) *Provider {
	if option.Name == "" {
		ecl.NewLogger(ecl.LoggerOption{Name: "DefineProvider"}).Panicf("Provider name cannot be empty")
	}

	return &Provider{
		Name:         option.Name,
		Instantiator: option.Instantiator,
		Handler:      HandlerName,
	}
}
