package food

import (
	"github.com/jhseong7/gimbap/core"
	"github.com/jhseong7/gimbap/core/controller"
	"github.com/jhseong7/gimbap/core/provider"
)

var FoodModuleGin = core.DefineModule(core.ModuleOption{
	Name: "FoodModule",
	Providers: []provider.ProviderDefinition{
		*FoodProvider,
	},
	Controllers: []controller.ControllerDefinition{
		*FoodControllerGinProvider,
	},
})
