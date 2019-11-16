package api

import (
	"fmt"
	"github.com/kataras/iris"
	"main/uploadImg"
	"strconv"
)

// 获取图片列表
func ApiGetImgList(ctx iris.Context) {
	ctx.ContentType("application/json")
	pageIndex, _ := strconv.ParseInt(ctx.FormValue("page_index"), 10, 64)
	pageSize, _ := strconv.ParseInt(ctx.FormValue("page_size"), 10, 64)
	imgList := uploadImg.FindImgForDatabase(pageIndex, pageSize)
	if len(imgList) > 0 {
		ctx.JSON(map[string]interface{}{
			"status": "success",
			"data":   imgList,
			"msg":    "",
		})
	} else {
		ctx.JSON(map[string]interface{}{
			"status": "error",
			"data":   nil,
			"msg":    "数据为空",
		})
	}
}

// 上传图片接口
func ApiUploadImg(ctx iris.Context) {
	file, handler, _ := ctx.FormFile("uploadfile")
	defer file.Close()
	fmt.Println(ctx.FormValues())
	return
	uploadRes := uploadImg.UploadFileStream(file, handler.Filename, ctx.PostValue("alt"))
	ctx.JSON(uploadRes)
}
