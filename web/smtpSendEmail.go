package web

import (
	"email/smtp"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

type smtpSendEmailReq struct {
	Destination string `form:"destination" binding:"required,json"`
	Subject     string `form:"subject" binding:"required"`
	Format      int    `form:"format" binding:"eq=0|eq=1"`
	Body        string `form:"body" binding:"required"`
}

func SmtpSendEmail(c *gin.Context) {
	request := &smtpSendEmailReq{}
	err := c.ShouldBind(request)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	to := make([]string, 0)
	err = json.Unmarshal([]byte(request.Destination), &to)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	files := form.File["attachments"]
	attachments := make([]*smtp.Attachment, len(files))
	if len(files) > 0 {
		for i, file := range files {
			f, _ := file.Open()
			bytes, err := ioutil.ReadAll(f)
			f.Close()
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}
			content := base64.StdEncoding.EncodeToString(bytes)
			attachment := &smtp.Attachment{
				FileName:      file.Filename,
				Base64Content: content,
			}
			attachments[i] = attachment
		}
	}
	err = smtp.Send(request.Subject, to, request.Body, request.Format, attachments)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "发送成功")

}
