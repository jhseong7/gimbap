package controller

import (
	"fmt"
	"reflect"

	logger "github.com/jhseong7/ecl"
)

type (
	// Interface the controller must implement.
	IController interface {
		// The key method to get the route specs.
		GetRouteSpecs() []RouteSpec
	}

	// RESTful controller.
	ControllerDefinition struct {
		Name         string
		RootPath     string
		Instantiator interface{}
	}

	RouteSpec struct {
		Path    string // Route path to the handler. The full path will be RootPath + Path.
		Method  string // HTTP method (GET, POST, PUT, DELETE, etc.)
		Handler interface{}
	}

	// Redefine ProviderOption as ControllerOption.
	ControllerOption struct {
		// Name of the controller.
		Name string

		// Instantiation function of the controller.
		Instantiator interface{}

		// Root path of the controller. The full path of each route will be RootPath + Path.
		RootPath string
	}
)

func checkInstantiatorInterface(instantiator interface{}) {
	// Check if the instantiator's result type implements IController.
	instantiatorType := reflect.TypeOf(instantiator)
	returnType := instantiatorType.Out(0)

	controllerInterfaceType := reflect.TypeOf((*IController)(nil)).Elem()

	if !returnType.Implements(controllerInterfaceType) {
		panic(fmt.Sprintf("Controller %s's instantiator's result type does not implement IController", returnType))
	}

}

// Define a controller
func DefineController(option ControllerOption) *ControllerDefinition {
	checkInstantiatorInterface(option.Instantiator)

	if option.Name == "" {
		logger.NewLogger(logger.LoggerOption{Name: "DefineController"}).Panicf("Controller name cannot be empty")
	}

	return &ControllerDefinition{
		Name:         option.Name,
		RootPath:     option.RootPath,
		Instantiator: option.Instantiator,
	}
}
