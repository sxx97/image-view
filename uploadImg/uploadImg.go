package uploadImg

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"go.mongodb.org/mongo-driver/bson/bsontype"
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
	Alt string `bson:"dbalt",json:"jsonalt"` // bson和json映射
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

type imgResultData struct {
	Status string `json:"status"`
	Message string `json:"msg"`
	Data Image `json:"data"`
}

/**
* 上传文件
 */
func UploadImg(localFileName string) imgResultData {
	var (
		insertCode interface{}
		visitImgUrl string
		resultData imgResultData
	)
	visitHost := "https://tongpaotk.oss-cn-beijing.aliyuncs.com"
	filePathArr := strings.Split(localFileName, "/")
	objectName := filePathArr[len(filePathArr)-1] + strconv.FormatInt(time.Now().Unix(), 10)
	fmt.Println("存储对象名称: ", objectName)
	err := ossBucket.PutObjectFromFile(objectName, localFileName)
	if err != nil {
		handleError(err)
	} else {
		visitImgUrl = visitHost + objectName
	 	insertCode = mongoose.InsertDatabase(Image{
			Alt: "",
			Src: visitImgUrl,
		})
	}
	if insertCode == 1 {
		resultData = imgResultData{
			Status: "success",
			Message: "上传成功",
			Data: Image{
				Src: visitImgUrl,
			},
		}
	} else {
		resultData = imgResultData{
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

/**
* 获取同一空间下的图片列表
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
