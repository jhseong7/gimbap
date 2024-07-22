// file: provider.go
//
// This file defines the provider interface and its implementation.
package provider

import (
	logger "github.com/jhseong7/ecl"
)

type (
	Provider struct {
		Name         string
		Instantiator interface{}

		// The handler string is used to identify the handler in the provider. (e.g. Controller)
		Handler string
	}

	ProviderOption struct {
		Name         string
		Instantiator interface{}
	}
)

const (
	HandlerName string = "default"
)

// Define a provider
func DefineProvider(option ProviderOption) *Provider {
	if option.Name == "" {
		logger.NewLogger(logger.LoggerOption{Name: "DefineProvider"}).Panicf("Provider name cannot be empty")
	}

	return &Provider{
		Name:         option.Name,
		Instantiator: option.Instantiator,
		Handler:      HandlerName,
	}
}
