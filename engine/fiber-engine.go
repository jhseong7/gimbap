package engine

import (
	"fmt"
	"reflect"

	"github.com/gofiber/fiber/v2"

	"github.com/jhseong7/ecl"
	"github.com/jhseong7/gimbap/controller"
)

type (
	FiberHttpEngine struct {
		IServerEngine

		// The underlying http engine
		engine          *fiber.App
		globalApiPrefix string

		logger ecl.Logger
	}

	FiberHttpEngineOption struct {
		ServerEngineOption
		FiberConfig fiber.Config
	}
)

// Check if the handler is valid and cast it to gin.HandlerFunc.
//
// This is a helper function to check if the handler is valid and cast it to gin.HandlerFunc before registering it to the engine.
func (e *FiberHttpEngine) checkAndCastToFiberHandler(handler interface{}) fiber.Handler {
	// Check if the input 0 is *gin.Context
	handlerType := reflect.TypeOf(handler)
	if handlerType.NumIn() == 0 || handlerType.In(0) != reflect.TypeOf(&fiber.Ctx{}) {
		e.logger.Panicf("Handler's first parameter must be *fiber.Context --> got %s", handlerType.String())
	}

	// Cast the handler value's interface to Fiber's Handler type (func(*fiber.Ctx) error)
	return handler.(func(*fiber.Ctx) error)
}

func (e *FiberHttpEngine) RegisterController(rootPath string, instance controller.IController) {
	defer func() {
		if r := recover(); r != nil {
			e.logger.Panicf("Failed to register controller to path: %s", rootPath)
		}
	}()

	routeSpecs := instance.GetRouteSpecs()

	for _, routeSpec := range routeSpecs {
		checkMethodValidity(routeSpec.Method)
		fullPath := mergeRestPath(e.globalApiPrefix, rootPath, routeSpec.Path)

		// Register the route
		// Unlike gin, fiber does not have a method to register a route with a handler function.
		// so use a switch statement to register the route based on the method.
		switch routeSpec.Method {
		case "GET":
			e.engine.Get(fullPath, e.checkAndCastToFiberHandler(routeSpec.Handler))
		case "POST":
			e.engine.Post(fullPath, e.checkAndCastToFiberHandler(routeSpec.Handler))
		case "PUT":
			e.engine.Put(fullPath, e.checkAndCastToFiberHandler(routeSpec.Handler))
		case "DELETE":
			e.engine.Delete(fullPath, e.checkAndCastToFiberHandler(routeSpec.Handler))
		case "PATCH":
			e.engine.Patch(fullPath, e.checkAndCastToFiberHandler(routeSpec.Handler))
		case "OPTIONS":
			e.engine.Options(fullPath, e.checkAndCastToFiberHandler(routeSpec.Handler))
		case "HEAD":
			e.engine.Head(fullPath, e.checkAndCastToFiberHandler(routeSpec.Handler))
		default:
			e.logger.Panicf("Invalid HTTP method: %s. Must be one of (GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD)", routeSpec.Method)
		}

		// Get the name of the Handler function
		handlerName := runtimeFuncName(routeSpec.Handler)

		e.logger.Logf("Registered route: %-8s %-20s --> %s", routeSpec.Method, fullPath, handlerName)
	}
}

// Add middleware to the engine
func (e *FiberHttpEngine) AddMiddleware(middleware ...interface{}) {
	for _, m := range middleware {
		casted := e.checkAndCastToFiberHandler(m)
		e.engine.Use(casted)
	}
}

func (e *FiberHttpEngine) Run(port int) {
	if port == 0 {
		e.logger.Warn("Port is not set. Defaulting to 8080")
		port = 8080
	}

	e.logger.Logf("Starting the http engine on port %d", port)

	// TODO: add a way to set cert and key for https
	e.engine.Listen(fmt.Sprintf(":%d", port))
}

func (e *FiberHttpEngine) Stop() {
	e.logger.Log("Stopping the http engine")

	if err := e.engine.Shutdown(); err != nil {
		e.logger.Fatalf("Failed to shutdown the http engine: %v", err)
	}
}

// Internal function to create a new fiber engine
func createFiberHttpEngine(fiberConfig fiber.Config) (e *fiber.App) {
	// Inject the custom logger to the fiber logger
	// Initialize the custom logger

	// TODO: maybe add a way to configure the fiber engine configs
	e = fiber.New(fiberConfig)

	// TODO: fiber loggers?

	return
}

// Create a new http engine (for now, gin is the only supported engine)
func NewFiberHttpEngine(options ...FiberHttpEngineOption) *FiberHttpEngine {
	// Get the options
	var option FiberHttpEngineOption
	if len(options) > 0 {
		option = options[0]
	}

	// Create logger
	l := ecl.NewLogger(ecl.LoggerOption{
		Name: "FiberHttpEngine",
	})

	// Create gin engine with the logger
	e := createFiberHttpEngine(option.FiberConfig)

	return &FiberHttpEngine{
		engine:          e,
		logger:          l,
		globalApiPrefix: option.GlobalApiPrefix,
	}
}
