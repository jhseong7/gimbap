package main

import (
	"github.com/jhseong7/nassi-golang/core"
	"github.com/jhseong7/nassi-golang/sample"
	"github.com/labstack/echo/v4"
)

func main() {
	// app := core.CreateApp(core.AppOption{
	// 	AppName:   "SampleApp",
	// 	AppModule: sample.AppModule,
	// })

	// app.AddMiddleware(func(ctx *gin.Context) {
	// 	ctx.Header("Access-Control-Allow-Origin", "*")
	// 	ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	// 	ctx.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// })

	app := core.CreateApp(core.AppOption{
		AppName:    "SampleAppEcho",
		AppModule:  sample.AppModule,
		HttpEngine: core.NewEchoHttpEngine(),
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
