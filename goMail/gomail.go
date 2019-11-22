package goMail

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"strconv"
)
//
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
	m.SetHeader("text/html", body)

	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])

	err := d.DialAndSend(m)
	if err != nil {
		fmt.Println("发送邮件错误信息:", err)
	}
	return err
}