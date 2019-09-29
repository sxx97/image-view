package mongoose

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	//"time"
)

type Image struct {
	ID  string "_id,omitempty"               // 简写bson映射口
	Alt string `bson:"dbalt",json:"jsonalt"` // bson和json映射
	Src string // 属性名 为全小写的key
}

var databaseUrl string = "mongodb://root:12138@localhost:21000"

func init() {
	/*client, err := mongo.NewClient(options.Client().ApplyURI(databaseUrl))
	if err != nil {
		println("错误", err)
	} else {
		println("连接结果", client)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		println("错误信息2:", err)
	}
	collection := client.Database("test").Collection("trainers")
	imgExample := Image{
		ID:"adada",
		Alt: "alt内容",
		Src: "图片路径",
	}
	insertResult, err := collection.InsertOne(context.TODO(), imgExample)
	if err != nil {
		println("插入失败", err)
	} else {
		fmt.Println("插入成功: ", insertResult.InsertedID)
	}*/
}

func ConnectTestDatabase() {
	clientOptions := options.Client().ApplyURI(databaseUrl)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
}
