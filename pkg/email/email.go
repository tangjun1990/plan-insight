package email

import (
	"crypto/tls"
	"git.4321.sh/feige/flygo/core/kcfg"
	"gopkg.in/gomail.v2"
)

func SendEmail(subject string, message string, toUser []string) {

	host := kcfg.GetString("feige.email.host")
	port := kcfg.GetInt("feige.email.port")
	username := kcfg.GetString("feige.email.username")
	password := kcfg.GetString("feige.email.password")
	enableSSL := kcfg.GetBool("feige.email.enableSSL")

	m := gomail.NewMessage()
	m.SetHeader("From", username) // 发件人
	// m.SetHeader("From", "alias"+"<"+userName+">") // 增加发件人别名

	m.SetHeader("To", toUser...) // 收件人，可以多个收件人，但必须使用相同的 SMTP 连接
	//m.SetHeader("Cc", "******@qq.com")  // 抄送，可以多个
	//m.SetHeader("Bcc", "******@qq.com") // 暗送，可以多个
	m.SetHeader("Subject", subject) // 邮件主题

	// text/html 的意思是将文件的 content-type 设置为 text/html 的形式，浏览器在获取到这种文件时会自动调用html的解析器对文件进行相应的处理。
	// 可以通过 text/html 处理文本格式进行特殊处理，如换行、缩进、加粗等等
	body := message //fmt.Sprintf(message, "testUser")
	m.SetBody("text/html", body)

	d := gomail.NewDialer(host, port, username, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: enableSSL}

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
