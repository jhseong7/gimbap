package fiber_engine

import (
	"crypto/tls"
	"fmt"
	"reflect"
	"time"

	"github.com/gofiber/fiber/v2"
	r "github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/jhseong7/ecl"
	"github.com/jhseong7/gimbap/controller"
	"github.com/jhseong7/gimbap/engine"
)

type (
	FiberHttpEngine struct {
		engine.IServerEngine

		// The underlying http engine
		engine          *fiber.App
		globalApiPrefix string

		logger ecl.Logger
	}

	FiberHttpEngineOption struct {
		engine.ServerEngineOption
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
		engine.CheckMethodValidity(routeSpec.Method)
		fullPath := engine.MergeRestPath(e.globalApiPrefix, rootPath, routeSpec.Path)

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
		handlerName := engine.RuntimeFuncName(routeSpec.Handler)

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

func (e *FiberHttpEngine) Run(option engine.ServerRuntimeOption) {
	port := option.Port

	if port == 0 {
		e.logger.Warn("Port is not set. Defaulting to 8080")
		port = 8080
	}

	// Load the tls config if it is given
	// Split the case for TLS and non-TLS
	if option.TLSOption != nil {
		e.logger.Logf("Starting the http engine with TLS on port %d", port)

		// If the config is given directly, use it, else load the cert/key files
		var config *tls.Config
		if option.TLSOption.Config != nil {
			// Use the given tls config directly
			config = option.TLSOption.Config
		} else {
			config = &tls.Config{
				MinVersion: tls.VersionTLS12,
			}
		}

		// Only load the cert/key files if the config does not have a certificate
		if option.TLSOption.CertFile != "" && option.TLSOption.KeyFile != "" && config.Certificates == nil {
			var err error
			cert, err := tls.LoadX509KeyPair(option.TLSOption.CertFile, option.TLSOption.KeyFile)
			if err != nil {
				e.logger.Fatalf("Failed to load TLS config: %s", err)
			}

			config = &tls.Config{
				MinVersion:   tls.VersionTLS12,
				Certificates: []tls.Certificate{cert},
			}
		}

		// If the certificates are not loaded, panic
		if config.Certificates == nil {
			e.logger.Fatalf(
				"Failed to load TLS config: At least one of tls.Config.Certificates or 'CertFile and KeyFile' are required",
			)
		}

		// Create a listener with the tls config
		ln, err := tls.Listen("tcp", fmt.Sprintf(":%d", port), config)
		if err != nil {
			e.logger.Fatalf("Failed to create a tls listener: %s", err)
		}

		if err := e.engine.Listener(ln); err != nil {
			e.logger.Fatalf("Failed to start the http engine: %v", err)
		}
		return
	}

	e.logger.Logf("Starting the http engine on port %d", port)

	// TODO: add a way to set cert and key for https
	if err := e.engine.Listen(fmt.Sprintf(":%d", port)); err != nil {
		e.logger.Fatalf("Failed to start the http engine: %v", err)
	}
}

func (e *FiberHttpEngine) Stop() {
	e.logger.Log("Stopping the http engine (Max 5 seconds)")

	if err := e.engine.ShutdownWithTimeout(5 * time.Second); err != nil {
		e.logger.Fatalf("Failed to shutdown the http engine: %v", err)
	}
}

func (e *FiberHttpEngine) AddStatic(prefix, root string, config ...interface{}) {
	// Try and cast the config to fiber.Static
	var fiberStaticConfig fiber.Static
	if len(config) > 0 {
		fiberStaticConfig = config[0].(fiber.Static)
		defer func() {
			if r := recover(); r != nil {
				e.logger.Panicf("Failed to add static files to the engine: %s", r)
			}
		}()
	}

	e.engine.Static(prefix, root, fiberStaticConfig)
}

// Internal function to create a new fiber engine
func createFiberHttpEngine(fiberConfig fiber.Config) (e *fiber.App) {
	// Inject the custom logger to the fiber logger
	// Initialize the custom logger

	e = fiber.New(fiberConfig)

	// Add the recover middleware to the engine
	e.Use(r.New())

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
