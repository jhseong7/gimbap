package main

import (
	"github.com/jhseong7/gimbap/core"
	"github.com/jhseong7/gimbap/engine"
	sample "github.com/jhseong7/gimbap/sample/http-engine"
	"github.com/labstack/echo/v4"
)

func main() {
	app := core.CreateApp(core.AppOption{
		AppName:    "SampleAppEcho",
		AppModule:  sample.AppModuleEcho,
		HttpEngine: engine.NewEchoHttpEngine(),
	})

	app.AddMiddleware(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
			ctx.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			ctx.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

			return next(ctx)
		}
	})

	app.Run()
}
