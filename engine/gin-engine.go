package engine

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jhseong7/ecl"
	"github.com/jhseong7/gimbap/controller"
)

type (
	GinHttpEngine struct {
		IServerEngine

		// The underlying http engine
		engine          *gin.Engine
		globalApiPrefix string

		server *http.Server

		logger ecl.Logger
	}

	GinLogger struct {
		io.Writer
		logger ecl.Logger
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

func (e *GinHttpEngine) RegisterController(rootPath string, instance controller.IController) {
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
		// Check if the handler is compatible with gin.HandlerFunc. else, panic so the user can fix it.
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

func (e *GinHttpEngine) Run(option ServerRuntimeOption) {
	port := option.Port

	if port == 0 {
		e.logger.Warn("Port is not set. Defaulting to 8080")
		port = 8080
	}

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
			// Use the given tls config directly
			config = option.TLSOption.tlsConfig
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

	e.logger.Logf("Starting the http engine on port %d", port)

	// Start the server. Http mode with no TLS
	if err := e.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		e.logger.Fatalf("Failed to start the http engine: %s", err)
	}

}

func (e *GinHttpEngine) Stop() {
	e.logger.Log("Stopping the http engine (Max 5 seconds)")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.server.Shutdown(ctx); err != nil {
		e.logger.Fatalf("Failed to shutdown the http engine: %s", err)
	}

	select {
	case <-ctx.Done():
		e.logger.Warn("Server failed to shutdown gracefully with in 5 seconds")
	}
}

// Internal function to create a new gin engine
func createGinHttpEngine(logger ecl.Logger) (e *gin.Engine) {
	// Set gin to release mode (suppresses debug messages)
	gin.SetMode(gin.ReleaseMode)

	e = gin.New()
	e.Use(gin.Recovery())
	e.Use(gin.LoggerWithWriter(&GinLogger{logger: logger}))

	return
}

// Create a new http engine (for now, gin is the only supported engine)
func NewGinHttpEngine(options ...ServerEngineOption) *GinHttpEngine {
	// Get the options
	var option ServerEngineOption
	if len(options) > 0 {
		option = options[0]
	}

	// Create logger
	l := ecl.NewLogger(ecl.LoggerOption{
		Name: "GinHttpEngine",
	})

	// Create gin engine with the logger
	e := createGinHttpEngine(l)

	return &GinHttpEngine{
		engine:          e,
		logger:          l,
		globalApiPrefix: option.GlobalApiPrefix,
	}
}
