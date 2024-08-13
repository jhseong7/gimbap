// File: http-engine.go
//
// This file defines the http engine interface and its implementation. The http engine is responsible for handling RESTful requests.
// This framework uses gin as the default http engine.
package engine

import (
	"crypto/tls"
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/jhseong7/gimbap/controller"
)

type (
	// The engine must implement this interface to be used in the app.
	IServerEngine interface {
		// Registers a controller to the engine
		RegisterController(rootPath string, controller controller.IController)

		// Run the server on the specified port
		Run(option ServerRuntimeOption)

		// Stop the server gracefully
		// Must implement a timeout if there is a possibility of a hanging.
		Stop()

		// Add middleware to the engine.
		// This will be native to the engine's core
		AddMiddleware(middleware ...interface{})
	}

	ServerEngineOption struct {
		GlobalApiPrefix string
	}

	TLSOption struct {
		// The path to the certificate/key files
		CertFile string
		KeyFile  string

		// The tls config can be given directly.
		// If the Certificate is given through the config
		//, In this case, the CertFile and KeyFile will be ignored.
		Config *tls.Config
	}

	ServerRuntimeOption struct {
		Port int

		TLSOption *TLSOption
	}
)

// Common util functions
func mergeRestPath(paths ...string) string {
	processedPaths := make([]string, 0)

	// Remove trailing and leading slashes
	for _, p := range paths {
		if p == "" {
			continue
		}

		processedPaths = append(processedPaths, strings.Trim(p, "/"))
	}

	return "/" + strings.Join(processedPaths, "/")
}

func checkMethodValidity(method string) {
	switch method {
	case "GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD":
		return
	default:
		panic(fmt.Sprintf("Invalid HTTP method: %s. Must be one of (GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD)", method))
	}
}

// Get the name of the function. the split the fi
func runtimeFuncName(f interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	split := strings.Split(name, ".")

	// Join the last 2 elements of the split
	return split[len(split)-2] + "." + strings.Replace(split[len(split)-1], "-fm", "", 1)
}
