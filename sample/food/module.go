package food

import "github.com/jhseong7/nassi-golang/core"

var FoodModule = core.DefineModule(core.ModuleOption{
	Name: "FoodModule",
	Providers: []core.ProviderDefinition{
		*FoodProvider,
	},
	Controllers: []core.ControllerDefinition{
		*FoodControllerProvider,
	},
})
