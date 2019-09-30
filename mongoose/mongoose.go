package mongoose

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var (
	databaseCollection *mongo.Collection
	databaseUrl string = "mongodb://root:12138@localhost:21000"
)

func init() {
	client, err := mongo.NewClient(options.Client().ApplyURI(databaseUrl))
	if err != nil {
		fmt.Println("创建mongodb错误: ", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20 * time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		println("连接数据库错误: ", err)
	}

	collection := client.Database("test").Collection("trainers")
	databaseCollection = collection
}

func InsertDatabase(data interface{}) int64 {
	insertResult, err := databaseCollection.InsertOne(context.TODO(), data)
	if err != nil {
		println("插入失败", err)
		return 0
	} else {
		fmt.Println("插入成功: ", insertResult.InsertedID)
		return 1
	}
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
