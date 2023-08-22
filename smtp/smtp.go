package smtp

import (
	"bytes"
	"crypto/tls"
	"email/config"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"strings"
	"time"
)

const (
	TXT  = 0
	HTML = 1
)

type Attachment struct {
	FileName      string
	Base64Content string
}

func getHeader(subject string, to []string, format int, boundary string) map[string]string {
	header := make(map[string]string)
	header["From"] = config.Conf.GetString("SmtpConf.EmailName") + " <" + config.Conf.GetString("SmtpConf.Email") + ">"
	header["To"] = strings.Join(to, ";")
	header["Subject"] = subject
	header["Content-Type"] = "text/plain; charset=UTF-8"
	if format == HTML {
		header["Content-Type"] = "text/html; charset=UTF-8"
	}
	if boundary != "" {
		header["Content-Type"] = "multipart/mixed;boundary=" + boundary
		//该字段暂时没有用到 ，默认传1.0
		header["Mime-Version"] = "1.0"
	}
	//该字段暂时没有用到
	header["Date"] = time.Now().String()
	return header
}

func getMessage(header map[string]string, body string, format int, boundary string) string {
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	if boundary == "" {
		message += "\r\n" + body
		return message
	}
	contentType := "text/plain; charset=UTF-8"
	if format == HTML {
		contentType = "text/html; charset=UTF-8"
	}
	message += "\r\n--" + boundary + "\r\n"
	message += "Content-Type:" + contentType + "\r\n"
	message += "\r\n" + body + "\r\n"
	return message
}

func Send(subject string, to []string, body string, format int, attachments []*Attachment) error {
	boundary := ""
	if len(attachments) > 0 {
		boundary = "-------------GoBoundary---------------"
	}
	header := getHeader(subject, to, format, boundary)
	message := getMessage(header, body, format, boundary)
	buffer := bytes.NewBuffer(nil)
	buffer.WriteString(message)
	if boundary != "" {
		//将附件添加到正文
		getAttachmentsMessage(buffer, attachments, boundary)
	}
	auth := getAuth()
	email := config.Conf.GetString("SmtpConf.Email")
	err := SendMailWithTLS(getAddr(), auth, email, to, buffer.Bytes())
	if err != nil {
		log.Println(err)
		return err
	}
	return err
}

func getAddr() string {
	host := config.Conf.GetString("SmtpConf.Host")
	port := config.Conf.GetString("SmtpConf.Port")
	addr := fmt.Sprintf("%s:%s", host, port)
	return addr
}

func getAuth() smtp.Auth {
	email := config.Conf.GetString("SmtpConf.Email")
	password := config.Secret.GetString("Smtp.Password")
	host := config.Conf.GetString("SmtpConf.Host")
	auth := smtp.PlainAuth("", email, password, host)
	return auth
}

func getAttachmentsMessage(buffer *bytes.Buffer, attachments []*Attachment, boundary string) {
	for _, attachment := range attachments {
		message := "\r\n--" + boundary + "\r\n"
		message += "Content-Transfer-Encoding:base64\r\n"
		message += "Content-Disposition:attachment\r\n"
		name := "./" + attachment.FileName
		//name = mime.BEncoding.Encode("UTF-8", name)
		message += "Content-Type:" + "application/octet-stream" + ";name=\"" + name + "\"\r\n"
		buffer.WriteString(message)
		writeBase64(buffer, attachment.Base64Content)
	}
}

// Dial return a smtp client
func Dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		log.Println("tls.Dial Error:", err)
		return nil, err
	}

	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

// SendMailWithTLS send email with tls
func SendMailWithTLS(addr string, auth smtp.Auth, from string,
	to []string, msg []byte) (err error) {
	//create smtp client
	c, err := Dial(addr)
	if err != nil {
		log.Println("Create smtp client error:", err)
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

// writeFile read file to buffer
func writeBase64(buffer *bytes.Buffer, base64Content string) {
	payload := []byte(base64Content)
	buffer.WriteString("\r\n")
	for index, line := 0, len(payload); index < line; index++ {
		buffer.WriteByte(payload[index])
		if (index+1)%76 == 0 {
			buffer.WriteString("\r\n")
		}
	}
}
