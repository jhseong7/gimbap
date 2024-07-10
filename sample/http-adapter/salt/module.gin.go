package salt

import (
	"github.com/jhseong7/gimbap/core"
	"github.com/jhseong7/gimbap/core/controller"
	"github.com/jhseong7/gimbap/core/provider"
)

var SaltModuleGin = core.DefineModule(core.ModuleOption{
	Name: "SaltModule",
	Providers: []provider.ProviderDefinition{
		*SaltProvider,
	},
	Controllers: []controller.ControllerDefinition{
		*SaltControllerGinProvider,
	},
})
