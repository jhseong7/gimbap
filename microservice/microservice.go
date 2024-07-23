package microservice

import "github.com/jhseong7/gimbap/provider"

type (
	// Interface to define a microservice
	//
	// This interface defines the basic methods to manage a microservice.
	// All microservices must implement this interface.
	IMicroService interface {
		// Start the microservice
		Start()

		// Gracefully stop the microservice
		Stop()
	}

	// Structure to define a microservice. This is used to define a microservice in the app.
	// NOTE: since Microservice provider does not have any extra fields for now, just alias provider
	// however for future extensibility, it is better to keep it separate.
	MicroServiceProvider struct {
		provider.Provider
	}

	MicroServiceProviderOption struct {
		Name         string
		Instantiator interface{}
	}
)

const (
	HandlerName provider.ProviderHandlerName = "microservice"
)

// Check if the input implements the IMicroService interface.
func IsMicroService(i interface{}) bool {
	_, ok := i.(IMicroService)
	return ok
}

// Define a microservice
func DefineMicroService(option MicroServiceProviderOption) *MicroServiceProvider {
	return &MicroServiceProvider{
		Provider: provider.Provider{
			Name:         option.Name,
			Instantiator: option.Instantiator,
			Handler:      HandlerName,
		},
	}
}
