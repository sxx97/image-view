package api

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"main/mongoose"
	"mime/multipart"
	"os"
	"strconv"
)

var multiUploadImgList []string

func init() {
	initUploadCollections()
}

// 反馈建议
func FeedBackAdvise(ctx iris.Context) {
	// 获取用户id
	userId := JWTParse(ctx)
	note := ctx.FormValue("note")
	if note == "" {
		_, _ = ctx.JSON(ResponseResult{
			Status: "error",
			Message: "请填写反馈内容",
		})
		return
	}
	feedBackNoteCollections := mongoose.NewMgo("test", "feedbackNote")
	feedBackNoteCollections.InsertDatabase(FeedbackNote{
		Note: note,
		UserId: userId,
	})
	fmt.Println("note内容:", note)
	_, _ = ctx.JSON(ResponseResult{
		"success",
		"感谢您的反馈与建议",
		nil,
	})
}

// 获取图片列表
func ApiGetImgList(ctx iris.Context) {
	/*ctx.ContentType("application/json")*/
	pageIndex, _ := strconv.ParseInt(ctx.FormValue("page_index"), 10, 64)
	pageSize, _ := strconv.ParseInt(ctx.FormValue("page_size"), 10, 64)
	imgList := FindImgForDatabase(pageIndex, pageSize)
	if len(imgList) > 0 {
		_, _ = ctx.JSON(map[string]interface{}{
			"status": "success",
			"data":   imgList,
			"msg":    "",
		})
	} else {
		_, _ = ctx.JSON(map[string]interface{}{
			"status": "error",
			"data":   nil,
			"msg":    "数据为空",
		})
	}
}

// 上传多图片接口
func ApiUploadMultiImg(ctx iris.Context) {
	// 获取用户id
	userId := JWTParse(ctx)
	var imgList []Image
	_, err := ctx.UploadFormFiles("./upload_img", func(c iris.Context, file *multipart.FileHeader) {
		multiUploadImgList = append(multiUploadImgList, file.Filename)
	})
	if err != nil {
		_, _ =ctx.JSON(ResponseResult{"上传失败,请重新上传", "error", nil})
		return
	}
	for _, fileName := range multiUploadImgList {
		filePath := "./upload_img/"+fileName
		imgList = append(imgList, UploadImg(filePath, userId))
		removeErr := os.Remove(filePath)
		if removeErr != nil {
			fmt.Println(fileName, "删除失败,错误:", removeErr)
			return
		}
	}
	multiUploadImgList = multiUploadImgList[0:0]
	if len(multiUploadImgList) == 0 {
		_, _ =ctx.JSON(ResponseResult{
			"success",
			"上传成功",
			imgList,
		})
	} else {
		_, _ =ctx.JSON(ResponseResult{"上传失败,请重新上传", "error", nil})
	}
	return
}

// 上传单张图片接口
func ApiUploadImg(ctx iris.Context) {
	// 获取用户id
	userId := JWTParse(ctx)
	file, handler, err := ctx.FormFile("uploadFile")
	if err != nil {
		fmt.Println("上传文件错误:", err)
		ctx.JSON(map[string]string{"message": "请选择文件再上传", "status": "error"})
		return
	}
	defer file.Close()
	uploadRes := UploadFileStream(file, handler.Filename, userId, ctx.FormValue("alt"))
	ctx.JSON(uploadRes)
}
