package food

import (
	"github.com/gin-gonic/gin"
	"github.com/jhseong7/nassi-golang/core"
)

type (
	FoodController struct {
		core.IController
		Food FoodService
	}
)

func (c *FoodController) GetRouteSpecs() []core.RouteSpec {
	return []core.RouteSpec{
		{Method: "GET", Path: "/", Handler: c.GetFood},
	}
}

func (c *FoodController) GetFood(ctx *gin.Context) {
	ctx.String(200, "Food: %s, Salt: %s and combined is %s", c.Food.Name, c.Food.Salt.Name, c.Food.GetName())
}

func (c *FoodController) PostFood(ctx *gin.Context) {
	ctx.String(200, "Post Food "+c.Food.GetName())
}

func NewFoodController(food *FoodService) *FoodController {
	return &FoodController{
		Food: *food,
	}
}

var FoodControllerProvider = core.DefineController(core.ControllerOption{
	Name:         "FoodController",
	Instantiator: NewFoodController,
	RootPath:     "food",
})
