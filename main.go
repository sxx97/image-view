package main

import (
	"crypto/md5"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	recover2 "github.com/kataras/iris/middleware/recover"
	"io"
	"main/uploadImg"
	"strconv"
	"time"
)

func main() {
	startServe()
}

func startServe() {
	var startIndex, pageSize int64 = 0, 10
	app := iris.New()
	app.Logger().SetLevel("debug")
	app.Use(recover2.New())
	app.Use(logger.New())
	app.RegisterView(iris.HTML("./webapp", ".html"))
	app.Handle("GET", "/", func(ctx iris.Context) {
		now := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(now, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))
		uploadImg.FindImgForDatabase(startIndex, pageSize)
		startIndex++
		ctx.View("index.html", token)
	})
	app.Handle("GET", "/root.txt", func(ctx iris.Context) {
		ctx.ServeFile("./root.txt", false)
	})
	app.Post("/upload", func(ctx iris.Context) {
		file, handler, _ := ctx.FormFile("uploadfile")
		defer file.Close()
		uploadImg.UploadFileStream(file, handler.Filename, ctx.PostValue("alt"))
	})
	app.Run(iris.TLS(":443", "mycreat.pem", "mykey.key"), iris.WithConfiguration(iris.TOML("./config/main.tml")))
}