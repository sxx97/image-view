package goMail

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"math/rand"
	"strconv"
	"time"
)

var captchaMap  = make(map[string]string)



// 生成随机验证码
func GenValidateCode(width int) string {
	numeric := [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())
	var codeStr string
	for i := 0; i < width; i++ {
		codeStr += strconv.Itoa(numeric[rand.Intn(r)])
	}
	return codeStr
}

// 发送验证码邮件
//
// mailTo 接收邮件的邮箱地址数组
// subject 邮件主题
func SendCaptchaEmail(mailTo []string, subject string) {
	for _, emailAddress := range mailTo {
		captchaMap[emailAddress] = GenValidateCode(6)
		_ = SendMail([]string{emailAddress}, subject, "<p>验证码: <b style='color: #ff6161;font-size: 20px;'>"+captchaMap[emailAddress]+"</b></p>")
		go deleteCaptcha(emailAddress)
	}
}

// 删除邮件验证码
//
// email 邮件名称
func deleteCaptcha(email string) {
	time.Sleep(time.Minute)
	delete(captchaMap, email)
}

// 验证邮件验证码是否符合
//
// email 邮件地址
// code 验证码
func CheckCaptchaCode(email string, code string) bool {
	if captchaMap[email] != code {
		return false
	}
	return true
}

// 发送验证码
//
// mailTo 接收邮件的邮箱地址数组
// subject 邮件主题
// body 邮件内容
func SendMail(mailTo []string, subject string, body string) error {
	mailConn := map[string]string {
		"user": "postmaster@tongpaotk.cn",
		"pass": "13184234719Ali",
		"host": "smtp.mxhichina.com",
		"port": "465",
	}

	port, _ := strconv.Atoi(mailConn["port"])
	m := gomail.NewMessage()
	m.SetHeader("From", mailConn["user"])
	m.SetHeader("To", mailTo...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])

	err := d.DialAndSend(m)
	if err != nil {
		fmt.Println("发送邮件错误信息:", err)
	}
	return err
}