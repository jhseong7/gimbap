package engine

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/jhseong7/ecl"
	"github.com/jhseong7/gimbap/controller"

	echo "github.com/labstack/echo/v4"

	"github.com/labstack/echo/v4/middleware"
)

type (
	EchoHttpEngine struct {
		IServerEngine

		// The underlying http engine
		engine          *echo.Echo
		globalApiPrefix string

		server *http.Server

		logger ecl.Logger
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

func (e *EchoHttpEngine) RegisterController(rootPath string, instance controller.IController) {
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

func (e *EchoHttpEngine) Run(option ServerRuntimeOption) {
	port := option.Port

	if port == 0 {
		e.logger.Warn("Port is not set. Defaulting to 8080")
		port = 8080
	}

	e.logger.Logf("Starting the http engine on port %d", port)

	// Create an http server
	e.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: e.engine,
	}

	// Split the case for TLS and non-TLS
	if option.TLSOption != nil {
		e.logger.Logf("Starting the http engine with TLS on port %d", port)

		// If the config is given directly, use it, else load the cert/key files
		var config *tls.Config
		if option.TLSOption.tlsConfig != nil {
			config = option.TLSOption.tlsConfig
		} else {
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

		// Create a listener with the tls config
		tlsListener, err := tls.Listen("tcp", e.server.Addr, config)
		if err != nil {
			e.logger.Fatalf("Failed to create a tls listener: %s", err)
		}

		// Run the server with the tls listener
		if err := e.server.Serve(tlsListener); err != nil && err != http.ErrServerClosed {
			e.logger.Fatalf("Failed to start the http engine: %s", err)
		}

		return
	}

	// Start the server. Http mode with no TLS
	if err := e.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		e.logger.Fatalf("Failed to start the http engine: %s", err)
	}
}

func (e *EchoHttpEngine) Stop() {
	e.logger.Log("Stopping the http engine (Max 5 seconds)")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// NOTE: shutdown through the server, not the engine as it was started with the server
	if err := e.server.Shutdown(ctx); err != nil {
		e.logger.Fatalf("Failed to shutdown the http engine: %v", err)
	}

	select {
	case <-ctx.Done():
		e.logger.Warn("Server failed to shutdown gracefully with in 5 seconds")
	}
}

// Internal function to create a new echo engine with the logger middleware
func createEchoHttpEngine(logger ecl.Logger) (e *echo.Echo) {

	e = echo.New()

	// Don't show the banner of echo and the port number info (it's redundant)
	e.HideBanner = true
	e.HidePort = true

	// Add the logger middleware to the engine
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		BeforeNextFunc: func(c echo.Context) {
			c.Set("customValueFromContext", 42)
		},
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			value, _ := c.Get("customValueFromContext").(int)
			logger.Logf("REQUEST: uri: %v, status: %v, custom-value: %v\n", v.URI, v.Status, value)
			return nil
		},
	}))

	return
}

// Create a new http engine (for now, gin is the only supported engine)
func NewEchoHttpEngine(options ...ServerEngineOption) *EchoHttpEngine {
	// Get the options
	var option ServerEngineOption
	if len(options) > 0 {
		option = options[0]
	}

	// Create logger
	l := ecl.NewLogger(ecl.LoggerOption{
		Name: "EchoHttpEngine",
	})

	// Create gin engine with the logger
	e := createEchoHttpEngine(l)

	return &EchoHttpEngine{
		engine:          e,
		logger:          l,
		globalApiPrefix: option.GlobalApiPrefix,
	}
}
