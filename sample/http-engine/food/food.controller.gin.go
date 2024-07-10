package food

import (
	"github.com/gin-gonic/gin"
	"github.com/jhseong7/gimbap/controller"
)

type (
	FoodControllerGin struct {
		controller.IController
		Food FoodService
	}
)

func (c *FoodControllerGin) GetRouteSpecs() []controller.RouteSpec {
	return []controller.RouteSpec{
		{Method: "GET", Path: "/", Handler: c.GetFood},
		{Method: "POST", Path: "/", Handler: c.PostFood},
	}
}

func (c *FoodControllerGin) GetFood(ctx *gin.Context) {
	ctx.String(200, "Food: %s, Salt: %s and combined is %s", c.Food.Name, c.Food.Salt.Name, c.Food.GetName())
}

func (c *FoodControllerGin) PostFood(ctx *gin.Context) {
	ctx.String(200, "Post Food "+c.Food.GetName())
}

func NewFoodControllerGin(food *FoodService) *FoodControllerGin {
	return &FoodControllerGin{
		Food: *food,
	}
}

var FoodControllerGinProvider = controller.DefineController(controller.ControllerOption{
	Name:         "FoodControllerGin",
	Instantiator: NewFoodControllerGin,
	RootPath:     "food",
})
