package food

import (
	"github.com/jhseong7/gimbap/controller"
	"github.com/jhseong7/gimbap/core"
	"github.com/jhseong7/gimbap/provider"
)

var FoodModuleGin = core.DefineModule(core.ModuleOption{
	Name: "FoodModuleGin",
	Providers: []provider.ProviderDefinition{
		*FoodProvider,
	},
	Controllers: []controller.ControllerDefinition{
		*FoodControllerGinProvider,
	},
})
