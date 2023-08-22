package web

import (
	"email/ses"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

type sesSendEmailReq struct {
	Destination string `form:"destination" binding:"required,json"`
	TempName    string `form:"temp_name" binding:"required"`
	TempData    string `form:"temp_data" binding:"omitempty,json"`
	Subject     string `form:"subject" binding:"required"`
	Unsubscribe string `form:"unsubscribe" binding:"omitempty,eq=0|eq=1"`
	TriggerType uint64 `form:"trigger_type" binding:"omitempty,eq=0|eq=1"`
}

func SesSendEmail(c *gin.Context) {
	request := &sesSendEmailReq{}
	err := c.ShouldBind(request)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	destination := make([]string, 0)
	err = json.Unmarshal([]byte(request.Destination), &destination)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	req := &ses.SendEmailReq{
		Destination: destination,
		TempName:    request.TempName,
		TempData:    request.TempData,
		Subject:     request.Subject,
		Unsubscribe: request.Unsubscribe,
		TriggerType: request.TriggerType,
	}
	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	files := form.File["attachments"]
	if len(files) > 0 {
		attachments := make([]*ses.Attachments, len(files))
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
			attachment := &ses.Attachments{
				FileName: file.Filename,
				Content:  content,
			}
			attachments[i] = attachment
		}
		req.Attachments = attachments
	}

	s := ses.NewSes()
	res, err := s.SendEmail(req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, res.Response)
}
