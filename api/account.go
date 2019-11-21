package api

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/kataras/iris"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"main/mongoose"
)

type User struct {
	Id int `bson:"id", json:"id"`
	Account string `bson:"account",json:"account"`
	Password string `bson:"password",json:"password"`
}

type ResponseResult struct  {
	Status string `json:"status"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

func AccountLogin(ctx iris.Context) {

}

/**
* 注册账号处理函数
*/
func RegisterAccount(ctx iris.Context) {
	account := ctx.FormValue("account")
	password := ctx.FormValue("password")
	if account == "" {
		ctx.JSON(ResponseResult{
			"error",
			"账号不能为空",
			nil,
		})
		return
	}

	if password == "" {
		ctx.JSON(ResponseResult{
			"error",
			"密码不能为空",
			nil,
		})
		return
	}
	isExists := isExistsAccount(account)
	if isExists {
		ctx.JSON(ResponseResult{
			"error",
			"账号已存在",
			nil,
		})
		return
	}
	accountTotal := getAccountTotal()
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(password + "tp"))
	cipherStr := md5Ctx.Sum(nil)
	encryptedData := hex.EncodeToString(cipherStr)
	createAccount(User{
		Id: accountTotal+1,
		Account: account,
		Password: encryptedData,
	})
	ctx.JSON(ResponseResult{
		"success",
		"注册成功!",
		nil,
	})
}

/**
* 创建账号(操作数据库)
*/
func createAccount(account User) {
	userCollections := mongoose.NewMgo("tongpao", "users")
	userCollections.InsertDatabase(account)
}

/**
* 查询账号是否已经存在
*/
func isExistsAccount(account string) bool {
	userCollections := mongoose.NewMgo("tongpao", "users")
	finResultArr := userCollections.FindDatabase(bson.D{{"account", account}}, options.Find())
	return len(finResultArr) > 0
}


/**
* 获取账号总数
*/
func getAccountTotal() (accountTotal int){
	userCollections := mongoose.NewMgo("tongpao", "users")
	return userCollections.FindDatabaseTotal()
}