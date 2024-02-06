package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jhseong7/nassi-golang/core"
	"github.com/jhseong7/nassi-golang/sample"
)

func main() {
	app := core.CreateApp(core.AppOption{
		AppName:   "SampleApp",
		AppModule: sample.AppModule,
	})

	app.AddMiddleware(func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	})

	app.Run()
}
