package food

import (
	"fmt"

	"github.com/jhseong7/gimbap/controller"
	"github.com/labstack/echo/v4"
)

type (
	FoodControllerEcho struct {
		controller.IController
		Food FoodService
	}
)

func (c *FoodControllerEcho) GetRouteSpecs() []controller.RouteSpec {
	return []controller.RouteSpec{
		{Method: "GET", Path: "/", Handler: c.GetFoodEcho},
		{Method: "POST", Path: "/", Handler: c.PostFoodEcho},
	}
}

func (c *FoodControllerEcho) GetFoodEcho(ctx echo.Context) error {
	ctx.String(200, fmt.Sprintf("Food: %s, Salt: %s and combined is %s", c.Food.Name, c.Food.Salt.Name, c.Food.GetName()))

	return nil
}

func (c *FoodControllerEcho) PostFoodEcho(ctx echo.Context) error {
	ctx.String(200, "Post Food "+c.Food.GetName())

	return nil
}

func NewFoodControllerEcho(food *FoodService) *FoodControllerEcho {
	return &FoodControllerEcho{
		Food: *food,
	}
}

var FoodControllerEchoProvider = controller.DefineController(controller.ControllerOption{
	Name:         "FoodControllerEcho",
	Instantiator: NewFoodControllerEcho,
	RootPath:     "food",
})
