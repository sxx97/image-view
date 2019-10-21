package mongoose

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type mgo struct {
	database string
	collection string
}

var (
	client *mongo.Client
	databaseUrl        string = "mongodb://root:12138@http://116.62.213.108:21000"
)

func NewMgo(database, collection string) *mgo {
	return &mgo{
		database,
		collection,
	}
}

func init() {
	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI(databaseUrl))
	if err != nil {
		fmt.Println("创建mongodb错误: ", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		println("连接数据库错误: ", err)
	}

}

func (m *mgo) InsertDatabase(data interface{}) int64 {
	collection := client.Database(m.database).Collection(m.collection)
	insertResult, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		println("插入失败", err)
		return 0
	} else {
		fmt.Println("插入成功: ", insertResult.InsertedID)
		return 1
	}
}


func (m *mgo) FindDatabase(filter bson.D, findOptions *options.FindOptions) (tempArr []*mongo.Cursor) {
	collection := client.Database(m.database).Collection(m.collection)
	temp, _ := collection.Find(context.Background(), filter, findOptions)
	tempArr = append(tempArr, temp)
	temp.Next(context.Background())
	return
}
