package main

import (
	"crypto/md5"
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	recover2 "github.com/kataras/iris/middleware/recover"
	"io"
	"main/uploadImg"

	//"main/uploadImg"
	//"os"
	"strconv"
	"time"
)

//import "main/mongoose"

func main() {
	startServe()
}


func startServe() {
	app := iris.New()

	//app.Use(recover2.New())
	app.Logger().SetLevel("debug")
	app.Use(recover2.New())
	app.Use(logger.New())
	app.RegisterView(iris.HTML("./webapp", ".html"))
	app.Handle("GET", "/", func(ctx iris.Context) {
		now := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(now, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))
		ctx.View("index.html", token)
	})
	app.Handle("GET", "/root.txt", func(ctx iris.Context) {
		ctx.ServeFile("./root.txt", false)
	})
	//
	app.Post("/upload", func(ctx iris.Context) {
		file, handler, _ := ctx.FormFile("uploadfile")
		defer file.Close()
		uploadImg.UploadFileStream(file, handler.Filename)
	})
	app.Get("/markdown", writeMarkdown)
	//app.Run(iris.Addr(":80"), iris.WithConfiguration(iris.TOML("./config/main.tml")))
	app.Run(iris.TLS(":443", "mycreat.pem", "mykey.key"), iris.WithConfiguration(iris.TOML("./config/main.tml")))
}



var markdownContents = []byte(`## Hello Markdown

This is a sample of Markdown contents

Features
--------

All features of Sundown are supported, including:

*   **Compatibility**. The Markdown v1.0.3 test suite passes with
    the --tidy option.  Without --tidy, the differences are
    mostly in whitespace and entity escaping, where blackfriday is
    more consistent and cleaner.

*   **Common extensions**, including table support, fenced code
    blocks, autolinks, strikethroughs, non-strict emphasis, etc.

*   **Safety**. Blackfriday is paranoid when parsing, making it safe
    to feed untrusted user input without fear of bad things
    happening. The test suite stress tests this and there are no
    known inputs that make it crash.  If you find one, please let me
    know and send me the input that does it.

    NOTE: "safety" in this context means *runtime safety only*. In order to
    protect yourself against JavaScript injection in untrusted content, see
    [this example](https://github.com/russross/blackfriday#sanitize-untrusted-content).

*   **Fast processing**. It is fast enough to render on-demand in
    most web applications without having to cache the output.

*   **Routine safety**. You can run multiple parsers in different
    goroutines without ill effect. There is no dependence on global
    shared state.

*   **Minimal dependencies**. Blackfriday only depends on standard
    library packages in Go. The source code is pretty
    self-contained, so it is easy to add to any project, including
    Google App Engine projects.

*   **Standards compliant**. Output successfully validates using the
    W3C validation tool for HTML 4.01 and XHTML 1.0 Transitional.

    [this is a link](https://github.com/kataras/iris)`)

func writeMarkdown(ctx iris.Context) {
	hashKey  := []byte("the-big-and-secret-fash-key-here")
	blockKey := []byte("lot-secret-of-characters-big-too")
	sc := securecookie.New(hashKey, blockKey)
	ctx.SetCookieKV("cookie", "dadaadadada", iris.CookieEncode(sc.Encode))
	ctx.Request().Cookie("cookie")
	ctx.Markdown(markdownContents)
}