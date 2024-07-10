package salt

import (
	"github.com/jhseong7/gimbap/controller"
	"github.com/jhseong7/gimbap/core"
	"github.com/jhseong7/gimbap/provider"
)

var SaltModuleFiber = core.DefineModule(core.ModuleOption{
	Name: "SaltModule",
	Providers: []provider.ProviderDefinition{
		*SaltProvider,
	},
	Controllers: []controller.ControllerDefinition{
		*SaltControllerFiberProvider,
	},
})
