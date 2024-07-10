package food

import (
	"github.com/jhseong7/gimbap/core"
	"github.com/jhseong7/gimbap/core/controller"
	"github.com/jhseong7/gimbap/core/provider"
)

var FoodModuleFiber = core.DefineModule(core.ModuleOption{
	Name: "FoodModuleFiber",
	Providers: []provider.ProviderDefinition{
		*FoodProvider,
	},
	Controllers: []controller.ControllerDefinition{
		*FoodControllerFiberProvider,
	},
})
