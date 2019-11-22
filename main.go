package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	recover2 "github.com/kataras/iris/middleware/recover"
	"main/api"
	"main/goMail"
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
		app.Logger().Print("请求连接:", path)
		if !strings.Contains(path, ".") {
			router(w, r)
			return
		}
		ctx := app.ContextPool.Acquire(w, r)
		fileServer(ctx)
		app.ContextPool.Release(ctx)
	})
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
	apiGroup := app.Party("/api")
	apiGroup.Handle("GET", "/img", api.ApiGetImgList)
	apiGroup.Post("/upload/img", api.ApiUploadImg)
	apiGroup.Post("/upload/multiImg", api.ApiUploadMultiImg)
	apiGroup.Post("/register", api.RegisterAccount)
	apiGroup.Post("/login", api.AccountLogin)
	apiGroup.Get("/email", func(ctx iris.Context) {
		goMail.SendMail(
			[]string{"1978417547@qq.com"},
		"tp测试发送邮件",
			"<h1>tp测试邮件内容</h1>")
	})
	/*api.Handle("GET", "/root.txt", func(ctx iris.Context) {
		ctx.ServeFile("./root.txt", false)
	})*/
}
