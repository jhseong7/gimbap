// file: provider.go
//
// This file defines the provider interface and its implementation.
package core

import logger "github.com/jhseong7/ecl"

type (
	ProviderDefinition struct {
		Name         string
		instantiator interface{}
	}

	ProviderOption struct {
		Name         string
		Instantiator interface{}
	}
)

// Define a provider
func DefineProvider(option ProviderOption) *ProviderDefinition {
	if option.Name == "" {
		logger.NewLogger(logger.LoggerOption{Name: "DefineProvider"}).Panicf("Provider name cannot be empty")
	}

	return &ProviderDefinition{
		Name:         option.Name,
		instantiator: option.Instantiator,
	}
}
