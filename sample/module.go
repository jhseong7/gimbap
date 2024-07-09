package sample

import (
	"github.com/jhseong7/gimbap/core"
	"github.com/jhseong7/gimbap/sample/food"
	"github.com/jhseong7/gimbap/sample/salt"
)

var AppModuleGin = core.DefineModule(core.ModuleOption{
	Name: "AppModule",
	SubModules: []core.Module{
		*salt.SaltModuleGin,
		*food.FoodModuleGin,
	},
})

var AppModuleEcho = core.DefineModule(core.ModuleOption{
	Name: "AppModule",
	SubModules: []core.Module{
		*salt.SaltModuleEcho,
		*food.FoodModuleEcho,
	},
})
