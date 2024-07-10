package salt

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jhseong7/gimbap/controller"
)

type (
	SaltControllerFiber struct {
		controller.IController
		saltSvc SaltService
	}
)

func (c *SaltControllerFiber) GetRouteSpecs() []controller.RouteSpec {
	return []controller.RouteSpec{
		{Method: "GET", Path: "/salt", Handler: c.GetSaltFiber},
		{Method: "GET", Path: "/salt/:id", Handler: c.GetSaltFiber},
		{Method: "POST", Path: "/salt", Handler: c.PostSaltFiber},
	}
}

func (c *SaltControllerFiber) GetSaltFiber(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if id == "" {
		id = "default"
	}

	ctx.Status(200).SendString(fmt.Sprintf("Get Salt: %s %s", id, c.saltSvc.GetName()))

	return nil
}

func (c *SaltControllerFiber) PostSaltFiber(ctx *fiber.Ctx) error {
	ctx.Status(200).SendString("Post Salt " + c.saltSvc.GetName())

	return nil
}

func NewSaltControllerFiber(saltSvc *SaltService) *SaltControllerFiber {
	return &SaltControllerFiber{
		saltSvc: *saltSvc,
	}
}

var SaltControllerFiberProvider = controller.DefineController(
	controller.ControllerOption{
		Name:         "SaltControllerFiber",
		Instantiator: NewSaltControllerFiber,
		RootPath:     "salt",
	})
