package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	recover2 "github.com/kataras/iris/middleware/recover"
	api2 "main/api"
	"net/http"
	"strings"
)

var app *iris.Application

func main() {
	initServe()
	indexHtml()
	apiParty()
	app.Run(iris.TLS(":443", "mycreat.pem", "mykey.key"), iris.WithConfiguration(iris.TOML("./config/main.tml")))
}

func initServe() {
	app = iris.New()
	app.Logger().SetLevel("debug")
	fileServer := app.StaticHandler("./webapp", false, false)

	app.WrapRouter(func(w http.ResponseWriter, r *http.Request, router http.HandlerFunc) {
		path := r.URL.Path

		if !strings.Contains(path, ".") {
			router(w, r)
			return
		}
		ctx := app.ContextPool.Acquire(w, r)
		fileServer(ctx)
		app.ContextPool.Release(ctx)
	})
	/*
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("/webapp"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("/webapp"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("/webapp"))))*/
	app.Use(recover2.New())
	app.Use(logger.New())
}

func indexHtml() {
	app.RegisterView(iris.HTML("./webapp", ".html"))
	app.Get("/", func(ctx iris.Context) {
		ctx.View("index.html")
	})
	app.Get("/:page", func(ctx iris.Context) {
		if ctx.Path() != "/api" {
			ctx.View("index.html")
		}
	})
}

func apiParty() {
	api := app.Party("/api")
	api.Handle("GET", "/img", api2.ApiGetImgList)
	api.Post("/upload/img", api2.ApiUploadImg)
	/*api.Handle("GET", "/root.txt", func(ctx iris.Context) {
		ctx.ServeFile("./root.txt", false)
	})*/
}
