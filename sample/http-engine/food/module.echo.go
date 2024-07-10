package food

import (
	"github.com/jhseong7/gimbap/controller"
	"github.com/jhseong7/gimbap/core"
	"github.com/jhseong7/gimbap/provider"
)

var FoodModuleEcho = core.DefineModule(core.ModuleOption{
	Name: "FoodModuleEcho",
	Providers: []provider.ProviderDefinition{
		*FoodProvider,
	},
	Controllers: []controller.ControllerDefinition{
		*FoodControllerEchoProvider,
	},
})
