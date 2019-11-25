package main

import (
	"github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	recover2 "github.com/kataras/iris/v12/middleware/recover"
	"main/api"
	/*"net/http"
	"strings"*/
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
	app.HandleDir("/css", "./webapp/css")
	app.HandleDir("/js", "./webapp/js")
	app.HandleDir("/img", "./webapp/img")
	/*app.HandleDir("/root.txt", "./webapp")
	app.HandleDir("/jd_root.txt", "./webapp")*/
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
	jwtHandler := jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(api.SecretKey), nil
		},
		SigningMethod: jwt.SigningMethodHS512,
	})
	apiGroup := app.Party("/api")
	apiGroup.Get("/img", api.ApiGetImgList)
	apiGroup.Post("/upload/img", api.ApiUploadImg)
	apiGroup.Post("/upload/multiImg", api.ApiUploadMultiImg)
	apiGroup.Post("/register", api.RegisterAccount)
	apiGroup.Post("/login", api.AccountLogin)
	apiGroup.Get("/email", api.GetEmailCode)
	apiGroup.Get("/checkJWT", api.CheckJWTToken, jwtHandler.Serve)
	/*api.Handle("GET", "/root.txt", func(ctx iris.Context) {
		ctx.ServeFile("./root.txt", false)
	})*/
}
