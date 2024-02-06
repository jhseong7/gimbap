package salt

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jhseong7/nassi-golang/core"
)

type (
	SaltController struct {
		core.IController
		saltSvc SaltService
	}
)

func (c *SaltController) GetRouteSpecs() []core.RouteSpec {
	return []core.RouteSpec{
		{Method: "GET", Path: "/salt", Handler: c.GetSalt},
		{Method: "GET", Path: "/salt/:id", Handler: c.GetSalt},
		{Method: "POST", Path: "/salt", Handler: c.PostSalt},
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
