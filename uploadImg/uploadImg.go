package uploadImg

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"main/mongoose"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	ossBucket *oss.Bucket
)

type Image struct {
	ID  bsontype.Type "_id,omitempty"               // 简写bson映射口
	Alt string `bson:"alt",json:"alt"` // bson和json映射
	Src string // 属性名 为全小写的key
}

func init() {
	fmt.Println("OSS Go SDK Version: ", oss.Version)
	endpoint := "http://oss-cn-beijing.aliyuncs.com"
	accessKeyId := "LTAI4FhJUZB4WCLjtdcaHZiz"
	accessKeySecret := "7abBcSN7Y7qdgBQCUsnB0D78fs6bjJ"
	bucketName := "tongpaotk"
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		handleError(err)
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		handleError(err)
	}
	ossBucket = bucket
}

/**
* 接受错误信息
 */
func handleError(err error) {
	fmt.Println("Error:", err)
	os.Exit(-1)
}

type ImgResultData struct {
	Status string `json:"status"`
	Message string `json:"msg"`
	Data Image `json:"data"`
}

/**
* 上传文件
 */
func UploadImg(localFileName string) ImgResultData {
	var (
		insertCode interface{}
		visitImgUrl string
		resultData ImgResultData
	)
	imgCollections := mongoose.NewMgo("test", "trainers")
	visitHost := "https://tongpaotk.oss-cn-beijing.aliyuncs.com"
	filePathArr := strings.Split(localFileName, "/")
	objectName := filePathArr[len(filePathArr)-1] + strconv.FormatInt(time.Now().Unix(), 10)
	fmt.Println("存储对象名称: ", objectName)
	err := ossBucket.PutObjectFromFile(objectName, localFileName)
	if err != nil {
		handleError(err)
	} else {
		visitImgUrl = visitHost + objectName
	 	insertCode = imgCollections.InsertDatabase(Image{
			Alt: "",
			Src: visitImgUrl,
		})
	}
	if insertCode == 1 {
		resultData = ImgResultData{
			Status: "success",
			Message: "上传成功",
			Data: Image{
				Src: visitImgUrl,
			},
		}
	} else {
		resultData = ImgResultData{
			Status: "error",
			Message: "上传失败",
			Data: Image{},
		}
	}

	b, err := json.Marshal(resultData)
	if err != nil {
		fmt.Println("error: ", err)
	}

	fmt.Println(b, "上传结果信息")
	os.Stdout.Write(b)
	return resultData
}

//	上传文件(流形式)
//	fd 文件流
//	fileName 文件名称
//  alt 图片介绍
func UploadFileStream(fd io.Reader, fileName string, alt ...string){
	imgCollections := mongoose.NewMgo("test", "trainers")
	visitHost := "https://tongpaotk.oss-cn-beijing.aliyuncs.com"
	visitImgUrl := visitHost + fileName
	imgCollections.InsertDatabase(Image{
		Alt: alt[0],
		Src: visitImgUrl,
	})
	err := ossBucket.PutObject(fileName, fd)
	if err != nil {
		fmt.Println("Error:", err)
		//os.Exit(-1)
	}
}

// 查询图片列表(数据库)
func FindImgForDatabase(startIndex, pageSize int64) {
	fmt.Println("获取数据的分页参数:", startIndex, pageSize)
	imgCollections := mongoose.NewMgo("test", "trainers")
	result := imgCollections.FindDatabase(bson.D{}, options.Find().SetSort(bson.D{{"_id", 1}}).SetSkip(startIndex * pageSize).SetLimit(pageSize))
	for _, item := range result {
		fmt.Println("图片列表:", item)
	}
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
