package salt

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jhseong7/gimbap/controller"
)

type (
	SaltControllerGin struct {
		controller.IController
		saltSvc SaltService
	}
)

func (c *SaltControllerGin) GetRouteSpecs() []controller.RouteSpec {
	return []controller.RouteSpec{
		{Method: "GET", Path: "/salt", Handler: c.GetSalt},
		{Method: "GET", Path: "/salt/:id", Handler: c.GetSalt},
		{Method: "POST", Path: "/salt", Handler: c.PostSalt},
	}
}

func (c *SaltControllerGin) GetSalt(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		id = "default"
	}

	ctx.String(200, fmt.Sprintf("Get Salt: %s %s", id, c.saltSvc.GetName()))
}

func (c *SaltControllerGin) PostSalt(ctx *gin.Context) {
	ctx.String(200, "Post Salt "+c.saltSvc.GetName())
}

func NewSaltControllerGin(saltSvc *SaltService) *SaltControllerGin {
	return &SaltControllerGin{
		saltSvc: *saltSvc,
	}
}

var SaltControllerGinProvider = controller.DefineController(
	controller.ControllerOption{
		Name:         "SaltControllerGin",
		Instantiator: NewSaltControllerGin,
		RootPath:     "salt",
	})
