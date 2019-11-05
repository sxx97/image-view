package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	recover2 "github.com/kataras/iris/middleware/recover"
)

var app *iris.Application
func main() {
	initServe()
	apiParty()
	app.Run(iris.TLS(":443", "mycreat.pem", "mykey.key"), iris.WithConfiguration(iris.TOML("./config/main.tml")))
}

func initServe() {
	app = iris.New()
	app.Logger().SetLevel("debug")
	app.Use(recover2.New())
	app.Use(logger.New())
}


func apiParty() {
	api := app.Party("/api")
	api.Handle("GET", "/img", apiGetImgList)
	api.Post("/upload", apiUploadImg)
	/*api.Handle("GET", "/root.txt", func(ctx iris.Context) {
		ctx.ServeFile("./root.txt", false)
	})*/
}
