package api

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/v12"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"main/goMail"
	"main/mongoose"
	"time"
)

const (
	SecretKey string = "My Secret"
)

type Token struct {
	Token string `json:"token"`
}

type User struct {
	Id int `bson:"id", json:"id"`
	Account string `bson:"account",json:"account"`
	Email string `bson:"email",json:"email"`
	Password string `bson:"password",json:"password"`
}

type ResponseResult struct  {
	Status string `json:"status"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

// 账号登录接口
func AccountLogin(ctx iris.Context) {
	account := ctx.FormValue("account")
	password := ctx.FormValue("password")
	userCollections := mongoose.NewMgo("tongpao", "users")
	accountList := userCollections.FindDatabase(bson.D{{"account", account}, {"password", EncryptAccount(password)}}, options.Find())
	if len(accountList) == 0 {
		_, _ = ctx.JSON(ResponseResult{
			"error",
			"账号或密码错误",
			nil,
		})
		return
	}
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		jwt.MapClaims{
			"id": accountList[0]["id"],
			"exp": time.Now().Add(time.Hour * time.Duration(2)).Unix(),
			"iat": time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		fmt.Println("加密token信息错误:", err)
	}
	_, _ = ctx.JSON(ResponseResult{
		Status:  "success",
		Message: "登录成功!",
		Data:    tokenString,
	})
}

// 检验token是否正确
func CheckJWTToken(ctx iris.Context) {
	token, isOk := ctx.Values().Get("jwt").(*jwt.Token)
	if !isOk {
		fmt.Println("断言失败:这不是Token, ===", ctx.Values().Get("jwt"))
	}
	//userId = token.Claims.(jwt.MapClaims)["id"]
	ctx.Writef("经过身份验证的请求\n")
	ctx.Writef("Header:%v\n", token.Header)
	ctx.Writef("Raw:%v\n", token.Raw)
	ctx.Writef("Valid:%v\n", token.Valid)

	ctx.Writef("claims:%v\n", token.Claims)
	ctx.Writef("id%v\n", token.Claims.(jwt.MapClaims)["id"])
	ctx.Writef("Method:%v\n", token.Method)
	//可以了解一下token的数据结构
	ctx.Writef("Signature:%v\n", token.Signature)
}

// 获取验证码
func GetEmailCode(ctx iris.Context) {
	email := ctx.FormValue("email")
	if email == "" {
		_, _ = ctx.JSON(ResponseResult{
			"error",
			"请输入邮箱",
			nil,
		})
	} else {
		goMail.SendCaptchaEmail([]string{email}, "tp验证码")
		_, _ = ctx.JSON(ResponseResult{
			"success",
			"验证码已发送至您的邮箱,请注意查收",
			nil,
		})
	}
}

// 加密账号
func EncryptAccount(password string) (encryptedData string) {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(password + "tp"))
	cipherStr := md5Ctx.Sum(nil)
	encryptedData = hex.EncodeToString(cipherStr)
	return
}


// 注册处理函数
func RegisterAccount(ctx iris.Context) {
	account := ctx.FormValue("account")
	password := ctx.FormValue("password")
	email := ctx.FormValue("email")
	code := ctx.FormValue("code")
	if account == "" {
		_, _ = ctx.JSON(ResponseResult{
			"error",
			"账号不能为空",
			nil,
		})
		return
	}
	if password == "" {
		_, _ = ctx.JSON(ResponseResult{
			"error",
			"密码不能为空",
			nil,
		})
		return
	}
	if email == "" {
		_, _ = ctx.JSON(ResponseResult{
			"error",
			"邮箱不能为空",
			nil,
		})
		return
	}
	if code == "" {
		_, _ = ctx.JSON(ResponseResult{
			"error",
			"验证码不能为空",
			nil,
		})
	}
	if !goMail.CheckCaptchaCode(email, code) {
		_, _ = ctx.JSON(ResponseResult{
			"error",
			"验证码不匹配,请确认",
			nil,
		})
		return
	}
	isExists := isExistsAccount(account)
	if isExists {
		_, _ = ctx.JSON(ResponseResult{
			"error",
			"账号已存在",
			nil,
		})
		return
	}
	accountTotal := getAccountTotal()
	createAccount(User{
		Id: accountTotal+1,
		Account: account,
		Email: email,
		Password: EncryptAccount(password),
	})
	_, _ = ctx.JSON(ResponseResult{
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