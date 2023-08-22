package ses

import (
	"email/config"
	"email/ses/template"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	ses "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002"
	"log"
)

var credential *common.Credential

type Ses struct {
	FromEmailAddress string
	ReplyToAddresses string
}

func init() {
	getCredential()
}

func NewSes() *Ses {
	return &Ses{
		FromEmailAddress: config.Conf.GetString("TencentCloudSes.FromEmailAddress"),
		ReplyToAddresses: config.Conf.GetString("TencentCloudSes.ReplyToAddresses"),
	}
}

func getCredential() {
	// 密钥可前往https://console.cloud.tencent.com/cam/capi网站进行获取
	credential = common.NewCredential(
		config.Secret.GetString("TencentCloudApiKey.SecretId"),
		config.Secret.GetString("TencentCloudApiKey.SecretKey"),
	)
}

func getClient(region ...string) (*ses.Client, error) {
	defaultRegion := config.Conf.GetString("TencentCloudSes.DefaultRegion")
	if len(region) > 0 {
		defaultRegion = region[0]
	}
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = config.Conf.GetString("TencentCloudSes.Endpoint")
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, err := ses.NewClient(credential, defaultRegion, cpf)
	return client, err
}

type Attachments struct {
	FileName string
	Content  string
}
type SendEmailReq struct {
	Destination []string
	TempName    string
	TempData    string
	Subject     string
	Attachments []*Attachments
	Unsubscribe string
	TriggerType uint64
}
type SendEmailRes struct {
	*ses.SendEmailResponse
}

func (s *Ses) SendEmail(req *SendEmailReq) (res *SendEmailRes, err error) {
	client, err := getClient()
	if err != nil {
		log.Println(err)
		return
	}
	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := ses.NewSendEmailRequest()
	request.FromEmailAddress = common.StringPtr(s.FromEmailAddress)
	request.ReplyToAddresses = common.StringPtr(s.ReplyToAddresses)
	request.Destination = common.StringPtrs(req.Destination)
	temp, err := template.GetTemplate(req.TempName)
	if err != nil {
		log.Println(err)
		return
	}
	request.Template = &ses.Template{
		TemplateID:   common.Uint64Ptr(temp.ID),
		TemplateData: common.StringPtr(req.TempData),
	}
	request.Subject = common.StringPtr(req.Subject)
	attachments := make([]*ses.Attachment, len(req.Attachments))
	for i, a := range req.Attachments {
		attachments[i] = &ses.Attachment{
			FileName: common.StringPtr(a.FileName),
			Content:  common.StringPtr(a.Content),
		}
	}
	request.Attachments = attachments
	request.Unsubscribe = common.StringPtr(req.Unsubscribe)
	request.TriggerType = common.Uint64Ptr(req.TriggerType)
	// 返回的resp是一个SendEmailResponse的实例，与请求对象对应
	response, err := client.SendEmail(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return
	}
	res = &SendEmailRes{
		SendEmailResponse: response,
	}
	return
}
