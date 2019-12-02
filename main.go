package main

import (
	"main/api"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/middleware/logger"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	recover2 "github.com/kataras/iris/v12/middleware/recover"
)

var (
	app *iris.Application
	NO_CHECK_PATH []string = []string{"/api/upload/img", "/api/upload/multiImg"}
)

func main() {
	initServe()
	indexHtml()
	apiParty()
	_ = app.Run(iris.TLS(":443", "mycreat.pem", "mykey.key"), iris.WithConfiguration(iris.TOML("./config/main.tml")))
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
	api.StaticIndexPage()
	/*app.Get("/", func(ctx iris.Context) {
		_ = ctx.View("index.html")
	})
	app.Get("/:page", func(ctx iris.Context) {
		if ctx.Path() != "/api" {
			_ = ctx.View("index.html")
		}
	})*/
}

func apiParty() {
	// jwt验证
	jwtHandler := jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(api.SecretKey), nil
		},
		SigningMethod: jwt.SigningMethodHS512,
		ErrorHandler: func(ctx context.Context, err error) {
			// 略过不需要验证的接口
			for _, path := range NO_CHECK_PATH {
				if ctx.Path() == path {
					ctx.Next()
					return
				}
			}
			_, _ = ctx.JSON(api.ResponseResult{
				Status:  "error",
				Message: "请登录后再操作",
				Data:    nil,
			})
		},
	})
	apiGroup := app.Party("/api")
	apiGroup.Get("/img", api.ApiGetImgList)
	apiGroup.Post("/upload/img", jwtHandler.Serve, api.ApiUploadImg)
	apiGroup.Post("/upload/multiImg", jwtHandler.Serve, api.ApiUploadMultiImg)
	apiGroup.Post("/register", api.RegisterAccount)
	apiGroup.Post("/login", api.AccountLogin)
	apiGroup.Get("/email", api.GetEmailCode)
	apiGroup.Post("/feedback", jwtHandler.Serve, api.FeedBackAdvise)
}
