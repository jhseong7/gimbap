package main

import (
	"github.com/jhseong7/gimbap/core"
	"github.com/jhseong7/gimbap/core/engine"
	"github.com/jhseong7/gimbap/sample"
)

func main() {
	app := core.CreateApp(core.AppOption{
		AppName:    "SampleAppFiber",
		AppModule:  sample.AppModuleFiber,
		HttpEngine: engine.NewFiberHttpEngine(),
	})

	// Example of global middleware
	// app.AddMiddleware(func(ctx *fiber.Ctx) {
	// 	ctx.Set("Access-Control-Allow-Origin", "*")
	// 	ctx.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	// 	ctx.Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// })

	app.Run()
}
