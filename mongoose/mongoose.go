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
	database   string
	collection string
}

var (
	account  string
	password string
)

var (
	client      *mongo.Client
	databaseUrl string
)

func NewMgo(database, collection string) *mgo {
	return &mgo{
		database,
		collection,
	}
}

func init() {
	/*fmt.Println("请输入数据库账号:")
	fmt.Scanln(&account)
	fmt.Println("请输入数据库密码:")
	fmt.Scanln(&password)*/
	//TODO: 临时使用,提交时删除
	account = "root"
	password = "12138"
	databaseUrl = "mongodb://" + account + ":" + password + "@116.62.213.108:21000"
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

func (m *mgo) InsertDatabase(data interface{}) *mongo.InsertOneResult {
	collection := client.Database(m.database).Collection(m.collection)
	insertResult, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		println("插入失败", err)
	} else {
		fmt.Println("插入成功: ", insertResult.InsertedID)
	}
	return insertResult
}

func (m *mgo) FindDatabase(filter bson.D, findOptions *options.FindOptions) (tempArr []bson.M) {
	collection := client.Database(m.database).Collection(m.collection)
	cur, err := collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		fmt.Println("查询数据库的错误信息", err)
		return
	}
	for cur.Next(context.Background()) {
		var tempData bson.M
		cur.Decode(&tempData)
		tempArr = append(tempArr, tempData)
	}
	return
}

func (m *mgo) FindDatabaseTotal() int {
	collection := client.Database(m.database).Collection(m.collection)
	total, _ := collection.CountDocuments(nil, bson.D{})
	return int(total)
}
