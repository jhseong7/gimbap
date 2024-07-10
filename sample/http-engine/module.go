package sample

import (
	"github.com/jhseong7/gimbap/core"
	"github.com/jhseong7/gimbap/sample/http-engine/food"
	"github.com/jhseong7/gimbap/sample/http-engine/salt"
)

var AppModuleGin = core.DefineModule(core.ModuleOption{
	Name: "AppModuleGin",
	SubModules: []core.Module{
		*salt.SaltModuleGin,
		*food.FoodModuleGin,
	},
})

var AppModuleEcho = core.DefineModule(core.ModuleOption{
	Name: "AppModuleEcho",
	SubModules: []core.Module{
		*salt.SaltModuleEcho,
		*food.FoodModuleEcho,
	},
})

var AppModuleFiber = core.DefineModule(core.ModuleOption{
	Name: "AppModuleFiber",
	SubModules: []core.Module{
		*salt.SaltModuleFiber,
		*food.FoodModuleFiber,
	},
})
