package salt

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jhseong7/nassi-golang/core"
	"github.com/labstack/echo/v4"
)

type (
	SaltController struct {
		core.IController
		saltSvc SaltService
	}
)

func (c *SaltController) GetRouteSpecs() []core.RouteSpec {
	return []core.RouteSpec{
		// {Method: "GET", Path: "/salt", Handler: c.GetSalt},
		// {Method: "GET", Path: "/salt/:id", Handler: c.GetSalt},
		// {Method: "POST", Path: "/salt", Handler: c.PostSalt},
		{Method: "GET", Path: "/salt", Handler: c.GetSaltEcho},
		{Method: "GET", Path: "/salt/:id", Handler: c.GetSaltEcho},
		{Method: "POST", Path: "/salt", Handler: c.PostSaltEcho},
	}
}

func (c *SaltController) GetSalt(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		id = "default"
	}

	ctx.String(200, fmt.Sprintf("Get Salt: %s %s", id, c.saltSvc.GetName()))
}

func (c *SaltController) PostSalt(ctx *gin.Context) {
	ctx.String(200, "Post Salt "+c.saltSvc.GetName())
}

func (c *SaltController) GetSaltEcho(ctx echo.Context) error {
	id := ctx.Param("id")

	if id == "" {
		id = "default"
	}

	ctx.String(200, fmt.Sprintf("Get Salt: %s %s", id, c.saltSvc.GetName()))

	return nil
}

func (c *SaltController) PostSaltEcho(ctx echo.Context) error {
	ctx.String(200, "Post Salt "+c.saltSvc.GetName())

	return nil
}

func NewSaltController(saltSvc *SaltService) *SaltController {
	return &SaltController{
		saltSvc: *saltSvc,
	}
}

var SaltControllerProvider = core.DefineController(
	core.ControllerOption{
		Name:         "SaltController",
		Instantiator: NewSaltController,
		RootPath:     "salt",
	})
