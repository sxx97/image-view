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
	app.Run(iris.TLS("116.62.213.108:443", "mycreat.pem", "mykey.key"))
}
