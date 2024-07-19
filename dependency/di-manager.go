package dependency

import (
	"reflect"

	"github.com/jhseong7/gimbap/provider"
)

type (
	IDependencyManager interface {
		// Inject the dependencies to the target
		Inject(result map[reflect.Type]reflect.Value, providers []*provider.Provider)

		// Lifecycle methods
		OnStart()
		OnStop()
	}
)
