package api

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"main/mongoose"
	"os"
	"strings"
)

var (
	ossBucket *oss.Bucket
	ossClient *oss.Client
	bucketName string = "tongpaotk"
)




func initUploadCollections() {
	endpoint := "http://oss-cn-beijing.aliyuncs.com"
	accessKeyId := "LTAI4FhJUZB4WCLjtdcaHZiz"
	accessKeySecret := "7abBcSN7Y7qdgBQCUsnB0D78fs6bjJ"

	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		handleError(err)
	}
	ossClient = client
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		handleError(err)
	}
	ossBucket = bucket
}


// 静态网站托管
//
//
func StaticIndexPage() {
	err := ossClient.SetBucketWebsite(bucketName, "./website/index.html", "")
	if err != nil {
		fmt.Println("Error:", err)
	}
}

/**
* 接受错误信息
 */
func handleError(err error) {
	fmt.Println("Error:", err)
	os.Exit(-1)
}

type ImgResultData struct {
	Status  string `json:"status"`
	Message string `json:"msg"`
	Data    interface{}  `json:"data"`
}

// 上传文件(文件名称)
// localFileName 本地文件名称
// userId 上传图片用户的id
// 用于服务器上传
func UploadImg(localFileName string, userId int) (resultData Image) {
	var (
		visitImgUrl string
	)

	imgCollections := mongoose.NewMgo("tongpao", "imgs")
	visitHost := "https://tongpaotk.oss-cn-beijing.aliyuncs.com/"
	filePathArr := strings.Split(localFileName, "/")
	objectName := filePathArr[len(filePathArr)-1]
	fmt.Println("存储对象名称: ", objectName)
	err := ossBucket.PutObjectFromFile(objectName, localFileName)
	if err != nil {
		handleError(err)
		return
	}
	visitImgUrl = visitHost + objectName
	insertResult := imgCollections.InsertDatabase(Image{
		Alt: "",
		FullSrc: visitImgUrl,
		Src: objectName,
		UserId: userId,
	})
	if insertResult != nil {
		resultData = Image{
			ID: insertResult.InsertedID,
			Src: visitImgUrl,
		}
	}
	return
}

//	上传文件(流形式)
//	fd 文件流
//  userId 上传的用户的id
//	fileName 文件名称
//  alt 图片介绍
// 用于客户端上传
func UploadFileStream(fd io.Reader, fileName string, userId int, alt ...string) ImgResultData {
	var (
		visitImgUrl string
		resultData  ImgResultData
	)

	resultData = ImgResultData{
		Message: "上传失败",
		Status: "error",
		Data: nil,
	}

	err := ossBucket.PutObject(fileName, fd)
	if err != nil {
		fmt.Println("流文件上传结果:", err)
		return resultData
	}

	imgCollections := mongoose.NewMgo("tongpao", "imgs")
	visitHost := "https://tongpaotk.oss-cn-beijing.aliyuncs.com/"
	visitImgUrl = visitHost + fileName
	insertResult := imgCollections.InsertDatabase(Image{
		Alt: alt[0],
		Src: fileName,
		FullSrc: visitImgUrl,
		UserId: userId,
	})
	if insertResult != nil {
		resultData = ImgResultData{
			Message: "上传成功",
			Status: "success",
			Data: Image{
				ID: insertResult.InsertedID,
				Src: visitImgUrl,
			},
		}
	}
	return resultData
}

// 查询图片列表(数据库)
// pageIndex 分页页数
// pageSize 分页数量
func FindImgForDatabase(pageIndex, pageSize int64) (imgList []map[string]interface{}) {
	if pageSize == 0 {
		pageSize = 10
	}
	fmt.Println("获取数据的分页参数:", pageIndex, pageSize)
	imgCollections := mongoose.NewMgo("tongpao", "imgs")
	result := imgCollections.FindDatabase(bson.D{}, options.Find().SetSort(bson.D{{"_id", 1}}).SetSkip(pageIndex*pageSize).SetLimit(pageSize))
	for _, item := range result {
		imgList = append(imgList, item)
	}
	return
}

/**
* 获取同一空间下的图片列表(阿里oss)
 */
func GetImgList() {
	marker := ""
	for {
		lsRes, err := ossBucket.ListObjects(oss.Marker(marker))
		if err != nil {
			handleError(err)
		}
		for _, object := range lsRes.Objects {
			fmt.Println("Bucket: ", object.Key)
		}
		if lsRes.IsTruncated {
			marker = lsRes.NextMarker
		} else {
			break
		}
	}
}

/**
* 删除图片
 */
func DeleteImg(objectName string) {
	err := ossBucket.DeleteObject(objectName)
	if err != nil {
		handleError(err)
	}
}
