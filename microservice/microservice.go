package microservice

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
	MicroServiceProvider struct {
		Name         string
		Instantiator interface{}
	}
)

// Check if the input implements the IMicroService interface.
func IsMicroService(i interface{}) bool {
	_, ok := i.(IMicroService)
	return ok
}
