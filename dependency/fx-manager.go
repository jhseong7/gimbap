package dependency

import (
	"context"
	"reflect"
	"time"

	"github.com/jhseong7/ecl"
	"github.com/jhseong7/gimbap/provider"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

type (
	FxDependencyManager struct {
		IDependencyManager
		logger ecl.Logger
		fxApp  *fx.App
	}
)

func (f *FxDependencyManager) ResolveDependencies(instanceMap map[reflect.Type]reflect.Value, providerList []*provider.Provider) {
	start := time.Now()

	// List to save all providers.
	opList := []fx.Option{}

	// List to save return types of all providers.
	returnTypeList := []reflect.Type{}

	for _, p := range providerList {
		// Add the instantiator to the optionList
		opList = append(opList, fx.Provide(p.Instantiator))

		funcType := reflect.TypeOf(p.Instantiator)

		// Get the output types and save them to returnTypeList
		for i := 0; i < funcType.NumOut(); i++ {
			returnTypeList = append(returnTypeList, funcType.Out(i))
		}
	}

	// Create function type. Inputs --> all provider outputs (instantiators), returns --> nothing.
	functionType := reflect.FuncOf(returnTypeList, nil, false)

	// Create function to initialize all providers. at runtime
	function := reflect.MakeFunc(functionType, func(args []reflect.Value) []reflect.Value {
		for _, a := range args {
			// Save the instance to instanceMap
			instanceMap[reflect.TypeOf(a.Interface())] = a
			f.logger.Debugf("Provider instance created: %s", reflect.TypeOf(a.Interface()).String())
		}

		return nil
	})

	// Create a fx.Module as a provider for the instantiators
	fxProviders := fx.Module("AppModule", opList...)

	// Initializer function to inject all providers to the instanceMap
	initInvoker := fx.Invoke(function.Interface())

	f.fxApp = fx.New(
		// Logger settings (TODO: make this configurable)
		fx.WithLogger(func() fxevent.Logger { return fxevent.NopLogger }),

		// PROVIDERS
		fxProviders, // Get auto generated fx.Option from the module

		// INVOKEs
		initInvoker,
	)

	// Running the app will initialize all providers
	startCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := f.fxApp.Start(startCtx); err != nil {
		f.logger.Fatalf("Failed to start the fx app: %v", err)
	}

	f.logger.Debugf("Dependency resolution took %v", time.Since(start))
}

func (f *FxDependencyManager) OnStart() {
	// Do nothing. fx is already started in Inject
}

func (f *FxDependencyManager) OnStop() {
	// Stop the app
	stopCtx, cancel := context.WithTimeout(context.Background(), 5)
	defer cancel()
	if err := f.fxApp.Stop(stopCtx); err != nil {
		f.logger.Fatalf("Failed to stop the fx app: %v", err)
	}
}

func NewFxManager() *FxDependencyManager {
	return &FxDependencyManager{
		logger: ecl.NewLogger(ecl.LoggerOption{Name: "FxDepManager"}),
	}
}
