package api

import (
	"fmt"
	"github.com/kataras/iris"
	"main/uploadImg"
	"strconv"
)

// 获取图片列表
func ApiGetImgList(ctx iris.Context) {
	fmt.Println("获取图片列表");
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
	fmt.Println("上传接口调用")
	file, handler, err := ctx.FormFile("uploadfile")
	if err != nil {
		fmt.Println("上传文件错误:", err)
		ctx.JSON(map[string]string{"message": "请选择文件再上传", "status": "error"})
		return
	}
	defer file.Close()

	uploadRes := uploadImg.UploadFileStream(file, handler.Filename, ctx.FormValue("alt"))
	ctx.JSON(uploadRes)
}
