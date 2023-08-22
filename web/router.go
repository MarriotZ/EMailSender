package web

import "github.com/gin-gonic/gin"

func InitRouter(r *gin.Engine) {
	sesGroup := r.Group("/ses")
	sesGroup.POST("/send", SesSendEmail)

	smtpGroup := r.Group("/smtp")
	smtpGroup.POST("/send", SmtpSendEmail)

}
