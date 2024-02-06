package sample

import (
	"github.com/jhseong7/nassi-golang/core"
	"github.com/jhseong7/nassi-golang/sample/food"
	"github.com/jhseong7/nassi-golang/sample/salt"
)

var AppModule = core.DefineModule(core.ModuleOption{
	Name: "AppModule",
	SubModules: []core.Module{
		*salt.SaltModule,
		*food.FoodModule,
	},
})
