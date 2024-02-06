package salt

import "github.com/jhseong7/nassi-golang/core"

var SaltModule = core.DefineModule(core.ModuleOption{
	Name: "SaltModule",
	Providers: []core.ProviderDefinition{
		*SaltProvider,
	},
	Controllers: []core.ControllerDefinition{
		*SaltControllerProvider,
	},
})
