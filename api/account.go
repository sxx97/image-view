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

func AccountLogin(ctx iris.Context) {
	account := ctx.FormValue("account")
	password := ctx.FormValue("password")

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		jwt.MapClaims{
			"exp": time.Now().Add(time.Hour * time.Duration(2)).Unix(),
			"iat": time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		fmt.Println("加密token信息错误:", err)
	}
	fmt.Println("加密后的token信息:===>", tokenString)
	fmt.Println("账号信息: ", account, "  密码:", password)
	_, _ = ctx.JSON(ResponseResult{
		Status:  "success",
		Message: "登录成功!",
		Data:    tokenString,
	})
}

func CheckJWTToken(ctx iris.Context) {
	token, isOk := ctx.Values().Get("jwt").(*jwt.Token)
	if !isOk {
		fmt.Println("断言失败:这不是Token, ===", ctx.Values().Get("jwt"))
	}
	ctx.Writef("This is an authenticated request\n")
	ctx.Writef("Claim content:\n")
	//可以了解一下token的数据结构
	ctx.Writef("%s", token.Signature)
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

/**
* 注册账号处理函数
*/
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
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(password + "tp"))
	cipherStr := md5Ctx.Sum(nil)
	encryptedData := hex.EncodeToString(cipherStr)
	createAccount(User{
		Id: accountTotal+1,
		Account: account,
		Email: email,
		Password: encryptedData,
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