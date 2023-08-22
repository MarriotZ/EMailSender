package main

import (
	"email/web"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	web.InitRouter(r)
	r.Run(":8080")
	/*
		s := ses.NewSes()
		req := &ses.SendEmailReq{
			Destination: []string{"TEST@163.com"},
			TempName:    "邮件验证码",
			TempData:    "{\"code\":\"xxx\",\"m\":\"xx\"}",
			Subject:     "邮件主题",
			Unsubscribe: "1",
			TriggerType: 1,
		}
		s.SendEmail(req)
			// 实例化一个认证对象，入参需要传入腾讯云账户secretId，secretKey,此处还需注意密钥对的保密
			// 密钥可前往https://console.cloud.tencent.com/cam/capi网站进行获取
			credential := common.NewCredential(
				"AKIDPpMzRp1xf8MMIMQPAi293lEeebLsEn8E",
				"yhbR7bsvxjP4mKj02sTpxrZ3II2Yrxoc",
			)
			// 实例化一个client选项，可选的，没有特殊需求可以跳过
			cpf := profile.NewClientProfile()
			cpf.HttpProfile.Endpoint = "ses.tencentcloudapi.com"
			// 实例化要请求产品的client对象,clientProfile是可选的
			client, _ := ses.NewClient(credential, "ap-hongkong", cpf)

			// 实例化一个请求对象,每个接口都会对应一个request对象
			request := ses.NewSendEmailRequest()

			request.FromEmailAddress = common.StringPtr("zihantsang@hotmail.com")
			request.ReplyToAddresses = common.StringPtr("TEST@163.com")
			request.Destination = common.StringPtrs([]string{"TEST@163.com"})
			request.Template = &ses.Template{
				TemplateID:   common.Uint64Ptr(32856),
				TemplateData: common.StringPtr("{\"code\":\"xxx\",\"m\":\"xx\"}"),
			}
			request.Subject = common.StringPtr("邮件主题")
			request.Attachments = []*ses.Attachment{
				&ses.Attachment{
					FileName: common.StringPtr("file1"),
					Content:  common.StringPtr("content"),
				},
				&ses.Attachment{
					FileName: common.StringPtr("fil2"),
					Content:  common.StringPtr("content2"),
				},
			}
			request.Unsubscribe = common.StringPtr("1")
			request.TriggerType = common.Uint64Ptr(1)

			// 返回的resp是一个SendEmailResponse的实例，与请求对象对应
			response, err := client.SendEmail(request)
			if _, ok := err.(*errors.TencentCloudSDKError); ok {
				fmt.Printf("An API error has returned: %s", err)
				return
			}
			if err != nil {
				panic(err)
			}
			// 输出json格式的字符串回包
			fmt.Printf("%s", response.ToJsonString())
	*/
}
