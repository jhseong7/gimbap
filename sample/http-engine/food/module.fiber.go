package food

import (
	"github.com/jhseong7/gimbap/controller"
	"github.com/jhseong7/gimbap/core"
	"github.com/jhseong7/gimbap/provider"
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
