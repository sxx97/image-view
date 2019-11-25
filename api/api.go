package api

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"mime/multipart"
	"os"
	"strconv"
)

var multiUploadImgList []string

func init() {
	initUploadCollections()
}

// 获取图片列表
func ApiGetImgList(ctx iris.Context) {
	fmt.Println("获取图片列表");
	ctx.ContentType("application/json")
	pageIndex, _ := strconv.ParseInt(ctx.FormValue("page_index"), 10, 64)
	pageSize, _ := strconv.ParseInt(ctx.FormValue("page_size"), 10, 64)
	imgList := FindImgForDatabase(pageIndex, pageSize)
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

// 上传多图片接口
func ApiUploadMultiImg(ctx iris.Context) {
	fmt.Println("上传接口调用");
	var imgList []Image
	_, err := ctx.UploadFormFiles("./upload_img", func(c iris.Context, file *multipart.FileHeader) {
		multiUploadImgList = append(multiUploadImgList, file.Filename)
	})
	if err != nil {
		fmt.Println("上传多图片错误: ==", err)
	}
	for _, fileName := range multiUploadImgList {
		filePath := "./upload_img/"+fileName
		imgList = append(imgList, UploadImg(filePath))
		removeErr := os.Remove(filePath)
		if removeErr != nil {
			fmt.Println(fileName, "删除失败,错误:", removeErr)
		}
	}
	multiUploadImgList = multiUploadImgList[0:0]
	if len(multiUploadImgList) == 0 {
		ctx.JSON(map[string]interface{}{
			"status": "success",
			"message": "上传成功",
			"data": imgList,
		})
	} else {
		ctx.JSON(map[string]string{"message": "上传失败,请重新上传", "status": "error"})
	}
	return
}

func ApiUploadImg(ctx iris.Context) {
	fmt.Println("上传单张图片接口");
	file, handler, err := ctx.FormFile("uploadFile")
	if err != nil {
		fmt.Println("上传文件错误:", err)
		ctx.JSON(map[string]string{"message": "请选择文件再上传", "status": "error"})
		return
	}
	defer file.Close()
	uploadRes := UploadFileStream(file, handler.Filename, ctx.FormValue("alt"))
	ctx.JSON(uploadRes)
}
