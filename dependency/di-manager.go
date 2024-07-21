package dependency

import (
	"reflect"

	"github.com/jhseong7/gimbap/provider"
)

type (
	IDependencyManager interface {
		// ResolveDependencies the dependencies to the target
		ResolveDependencies(result map[reflect.Type]reflect.Value, providers []*provider.Provider)

		// Lifecycle methods
		OnStart()
		OnStop()
	}
)
