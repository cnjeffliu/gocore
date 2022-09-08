/*
 * @Author: Jeffrey Liu <zhifeng172@163.com>
 * @Date: 2021-10-21 15:36:56
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-09-08 13:49:41
 * @Description:
 */
package message

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/smtp"
)

const (
	HOST   = "smtp.163.com"
	PORT   = 465
	FROM   = "motech_cs@163.com"
	PASSWD = "FAKSLZQCQEHKGCLE"
)

func dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		log.Panicln("Dialing Error:", err)
		return nil, err
	}

	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

// 参考net/smtp的func SendMail()
// 使用net.Dial连接tls(ssl)端口时,smtp.NewClient()会卡住且不提示err
// len(to)>1时,to[1]开始提示是密送
func send(addr string, auth smtp.Auth, from string,
	to []string, msg []byte) (err error) {

	//create smtp client
	c, err := dial(addr)
	if err != nil {
		log.Println("Create smpt client error:", err)
		return err
	}
	defer c.Close()

	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				log.Println("Error during AUTH", err)
				return err
			}
		}
	}

	if err = c.Mail(from); err != nil {
		return err
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}

func SendEmailTLS(to []string, title string, body string) {
	header := make(map[string]string)
	header["From"] = "motech" + "<" + FROM + ">"
	header["To"] = to[0]
	header["Subject"] = title
	// header["Content-Type"] = "text/html;chartset=UTF-8" // 换行符不生效
	header["Content-Type"] = "text/plain;chartset=UTF-8"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s:%s\r\n", k, v)
	}

	message += "\r\n" + body

	auth := smtp.PlainAuth("", FROM, PASSWD, HOST)
	err := send(fmt.Sprintf("%s:%d", HOST, PORT), auth, FROM, to, []byte(message))
	if err != nil {
		fmt.Println(err)
	}
}
