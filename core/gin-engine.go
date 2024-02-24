package core

import (
	"bytes"
	"fmt"
	"io"
	"reflect"

	"github.com/gin-gonic/gin"
	logger "github.com/jhseong7/ecl"
)

type (
	GinHttpEngine struct {
		IHttpEngine

		// The underlying http engine
		engine          *gin.Engine
		globalApiPrefix string

		logger logger.Logger
	}

	GinLogger struct {
		io.Writer
		logger logger.Logger
	}
)

// Write the log to the logger
func (g *GinLogger) Write(p []byte) (n int, err error) {
	// Trim the last newline character (if exists)
	p = bytes.TrimRight(p, "\n")
	g.logger.Logf(string(p))
	return len(p), nil
}

// Check if the handler is valid and cast it to gin.HandlerFunc.
//
// This is a helper function to check if the handler is valid and cast it to gin.HandlerFunc before registering it to the engine.
func (e *GinHttpEngine) checkAndCastToGinHandler(handler interface{}) gin.HandlerFunc {
	// Check if the input 0 is *gin.Context
	handlerType := reflect.TypeOf(handler)
	if handlerType.NumIn() == 0 || handlerType.In(0) != reflect.TypeOf(&gin.Context{}) {
		e.logger.Panicf("Handler's first parameter must be *gin.Context: %s", handlerType.String())
	}

	// Cast the handler value's interface to gin.HandlerFunc
	return handler.(func(*gin.Context))
}

func (e *GinHttpEngine) RegisterController(rootPath string, instance IController) {
	defer func() {
		if r := recover(); r != nil {
			e.logger.Panicf("Failed to register controller to path: %s", rootPath)
		}
	}()

	routeSpecs := instance.GetRouteSpecs()

	for _, routeSpec := range routeSpecs {
		checkMethodValidity(routeSpec.Method)
		fullPath := mergeRestPath(e.globalApiPrefix, rootPath, routeSpec.Path)

		e.engine.Handle(routeSpec.Method, fullPath, e.checkAndCastToGinHandler(routeSpec.Handler))

		// Get the name of the Handler function
		handlerName := runtimeFuncName(routeSpec.Handler)

		e.logger.Logf("Registered route: %-8s %-20s --> %s", routeSpec.Method, fullPath, handlerName)
	}
}

// Add middleware to the engine
func (e *GinHttpEngine) AddMiddleware(middleware ...interface{}) {
	for _, m := range middleware {
		casted := e.checkAndCastToGinHandler(m)
		e.engine.Use(casted)
	}
}

func (e *GinHttpEngine) Run(port int) {
	if port == 0 {
		e.logger.Warn("Port is not set. Defaulting to 8080")
		port = 8080
	}

	e.logger.Logf("Starting the http engine on port %d\n", port)
	e.engine.Run(fmt.Sprintf(":%d", port))
}

func CreateGinHttpEngine(logger logger.Logger) (e *gin.Engine) {
	// Set gin to release mode (suppresses debug messages)
	gin.SetMode(gin.ReleaseMode)

	e = gin.New()
	e.Use(gin.Recovery())
	e.Use(gin.LoggerWithWriter(&GinLogger{logger: logger}))

	return
}

// Create a new http engine (for now, gin is the only supported engine)
func NewGinHttpEngine(options ...HttpEngineOption) IHttpEngine {
	// Get the options
	var option HttpEngineOption
	if len(options) > 0 {
		option = options[0]
	}

	// Create logger
	l := logger.NewLogger(logger.LoggerOption{
		Name: "GinHttpEngine",
	})

	// Create gin engine with the logger
	e := CreateGinHttpEngine(l)

	return &GinHttpEngine{
		engine:          e,
		logger:          l,
		globalApiPrefix: option.GlobalApiPrefix,
	}
}
