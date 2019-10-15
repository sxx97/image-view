package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
)

//import "main/mongoose"

func main() {
	startServe()
}

func startServe() {
	app := iris.New()
	//app.Use(recover2.New())
	app.Use(logger.New())
	app.Handle("GET", "/", func(ctx iris.Context) {
		ctx.HTML("<h1>开发中...</h1>")
	})
	//
	app.Get("/test", testRoute);
	app.Run(iris.TLS("172.16.47.62:443", "mycreat.pem", "mykey.key"), iris.WithConfiguration(iris.TOML("./config/main.tml")))
}

func testRoute(ctx iris.Context) {
	ctx.Writef("Hello from method: %s and path: %s", ctx.Method(), ctx.Path())
}
