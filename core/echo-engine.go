package core

import (
	"fmt"
	"reflect"

	"github.com/jhseong7/nassi-golang/logger"
	echo "github.com/labstack/echo/v4"
)

type (
	EchoHttpEngine struct {
		IHttpEngine

		// The underlying http engine
		engine          *echo.Echo
		globalApiPrefix string

		logger logger.Logger
	}
)

// Check if the handler is valid and cast it to gin.HandlerFunc.
//
// This is a helper function to check if the handler is valid and cast it to gin.HandlerFunc before registering it to the engine.
func (e *EchoHttpEngine) checkAndCastToEchoHandler(handler interface{}) echo.HandlerFunc {
	// Check if the input 0 is echo.Context
	handlerType := reflect.TypeOf(handler)

	if handlerType.NumIn() == 0 || handlerType.In(0) != reflect.TypeOf((*echo.Context)(nil)).Elem() {
		e.logger.Panicf("Handler's first parameter must be echo.Context: got %s %s", handlerType.String())
	}

	// Cast the handler value's interface to func(echo.Context) error
	return handler.(func(echo.Context) error)
}

func (e *EchoHttpEngine) checkAndCastToEchoMiddlewareHandler(handler interface{}) echo.MiddlewareFunc {
	// Check if the input 0 is echo.Context
	handlerType := reflect.TypeOf(handler)

	if handlerType.NumIn() == 0 || handlerType.In(0) != reflect.TypeOf((echo.HandlerFunc)(nil)) {
		e.logger.Panicf("Middleware Handler's first parameter must be echo.HandlerFunc:  got %s", handlerType.String())
	}

	// Also check if the output is echo.HandlerFunc which is a type func(echo.Context) error
	if handlerType.NumOut() == 0 || handlerType.Out(0) != reflect.TypeOf((echo.HandlerFunc)(nil)) {
		e.logger.Panicf("Middleware Handler's return type must be echo.HandlerFunc:  got %s", handlerType.String())
	}

	// Cast the handler value's interface to gin.HandlerFunc
	return handler.(func(echo.HandlerFunc) echo.HandlerFunc)
}

func (e *EchoHttpEngine) RegisterController(rootPath string, instance IController) {
	defer func() {
		if r := recover(); r != nil {
			e.logger.Panicf("Failed to register controller to path: %s", rootPath)
		}
	}()

	routeSpecs := instance.GetRouteSpecs()

	for _, routeSpec := range routeSpecs {
		checkMethodValidity(routeSpec.Method)
		fullPath := mergeRestPath(e.globalApiPrefix, rootPath, routeSpec.Path)

		e.engine.Add(routeSpec.Method, fullPath, e.checkAndCastToEchoHandler(routeSpec.Handler))

		// Get the name of the Handler function
		handlerName := runtimeFuncName(routeSpec.Handler)

		e.logger.Logf("Registered route: %-8s %-20s --> %s", routeSpec.Method, fullPath, handlerName)
	}
}

// Add middleware to the engine
func (e *EchoHttpEngine) AddMiddleware(middleware ...interface{}) {
	for _, m := range middleware {
		casted := e.checkAndCastToEchoMiddlewareHandler(m)
		e.engine.Use(casted)
	}
}

func (e *EchoHttpEngine) Run(port int) {
	if port == 0 {
		e.logger.Warn("Port is not set. Defaulting to 8080")
		port = 8080
	}

	e.logger.Logf("Starting the http engine on port %d\n", port)
	e.logger.Fatal(e.engine.Start(fmt.Sprintf(":%d", port)).Error())
}

func CreateEchoHttpEngine(logger logger.Logger) (e *echo.Echo) {

	e = echo.New()

	return
}

// Create a new http engine (for now, gin is the only supported engine)
func NewEchoHttpEngine(options ...HttpEngineOption) IHttpEngine {
	// Get the options
	var option HttpEngineOption
	if len(options) > 0 {
		option = options[0]
	}

	// Create logger
	l := logger.NewLogger(logger.LoggerOption{
		Name: "EchoHttpEngine",
	})

	// Create gin engine with the logger
	e := CreateEchoHttpEngine(l)

	return &EchoHttpEngine{
		engine:          e,
		logger:          l,
		globalApiPrefix: option.GlobalApiPrefix,
	}
}
