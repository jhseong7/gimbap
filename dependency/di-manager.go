package dependency

import (
	"reflect"

	"github.com/jhseong7/gimbap/provider"
)

type (
	IDependencyManager interface {
		// ResolveDependencies the dependencies to the target
		//
		// The first parameter is the result map which contains the resolved dependencies.
		// The second parameter is the list of providers to resolve.
		ResolveDependencies(instanceMap map[reflect.Type]reflect.Value, providers []*provider.Provider)

		// Lifecycle methods
		OnStart()
		OnStop()
	}
)
