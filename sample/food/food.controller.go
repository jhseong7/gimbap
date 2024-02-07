package food

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jhseong7/nassi-golang/core"
	"github.com/labstack/echo/v4"
)

type (
	FoodController struct {
		core.IController
		Food FoodService
	}
)

func (c *FoodController) GetRouteSpecs() []core.RouteSpec {
	return []core.RouteSpec{
		// {Method: "GET", Path: "/", Handler: c.GetFood},
		// {Method: "POST", Path: "/", Handler: c.PostFood},
		{Method: "GET", Path: "/", Handler: c.GetFoodEcho},
		{Method: "POST", Path: "/", Handler: c.PostFoodEcho},
	}
}

func (c *FoodController) GetFood(ctx *gin.Context) {
	ctx.String(200, "Food: %s, Salt: %s and combined is %s", c.Food.Name, c.Food.Salt.Name, c.Food.GetName())
}

func (c *FoodController) PostFood(ctx *gin.Context) {
	ctx.String(200, "Post Food "+c.Food.GetName())
}

func (c *FoodController) GetFoodEcho(ctx echo.Context) error {
	ctx.String(200, fmt.Sprintf("Food: %s, Salt: %s and combined is %s", c.Food.Name, c.Food.Salt.Name, c.Food.GetName()))

	return nil
}

func (c *FoodController) PostFoodEcho(ctx echo.Context) error {
	ctx.String(200, "Post Food "+c.Food.GetName())

	return nil
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
