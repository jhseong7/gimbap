// File: http-engine.go
//
// This file defines the http engine interface and its implementation. The http engine is responsible for handling RESTful requests.
// This framework uses gin as the default http engine.
package core

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

type (
	// TODO: make gin as a separate package
	IHttpEngine interface {
		RegisterController(rootPath string, controller IController)
		Run(port int)
		AddMiddleware(middleware ...interface{})
	}

	HttpEngineOption struct {
		GlobalApiPrefix string
	}
)

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
