package mail

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridConfig struct {
	Name   string // name of sender
	Email  string // email of sender
	ApiKey string // api key of sendgrid
}

type sendgridSender struct {
	config SendGridConfig
}

type SendGridOptions struct {
	TemplateID string
}

func NewSendGridSender(cfg SendGridConfig) IEmailSender {
	return &sendgridSender{
		config: cfg,
	}
}

func (sender *sendgridSender) Send(request *Request) (*Request, *Response) {

	opts, ok := request.Options.(SendGridOptions)
	if !ok {
		fmt.Println(opts.TemplateID)
	}

	e := mail.NewV3Mail()
	from := mail.NewEmail(sender.config.Name, sender.config.Email)
	e.SetFrom(from)

	e.Subject = request.Subject

	personalization := mail.NewPersonalization()
	for _, recipient := range request.To {
		personalization.AddTos(mail.NewEmail("", recipient))
	}
	for _, ccRecipient := range request.Cc {
		personalization.AddCCs(mail.NewEmail("", ccRecipient))
	}
	for _, bccRecipient := range request.Bcc {
		personalization.AddBCCs(mail.NewEmail("", bccRecipient))
	}
	e.AddPersonalizations(personalization)

	contentType := "text/plain"
	e.AddContent(mail.NewContent(contentType, request.Content))
	if opts.TemplateID != "" {
		e.TemplateID = opts.TemplateID
	}

	response := &Response{}
	for _, f := range request.AttachFiles {
		fileBytes, err := os.ReadFile(f)
		if err != nil {
			response.Message = fmt.Sprintf("failed to attach file %s: %s", f, err.Error())
			return request, response
		}

		encodedFile := base64.StdEncoding.EncodeToString(fileBytes)
		attachment := mail.NewAttachment()
		attachment.SetContent(encodedFile)
		attachment.SetType("application/octet-stream")
		attachment.SetFilename(f)
		attachment.SetDisposition("attachment")
		attachment.SetContentID("Attachment")

		e.AddAttachment(attachment)
	}

	client := sendgrid.NewSendClient(sender.config.ApiKey)
	res, err := client.Send(e)
	if err != nil {
		response.Message = fmt.Sprintf("failed to send email: %s", err.Error())
		return request, response
	}

	data, err := json.Marshal(res)
	if err != nil {
		response.Message = fmt.Sprintf("error marshalling response data: %s", err.Error())
		return request, response
	}
	response.Data = string(data)

	if res.StatusCode != 250 {
		response.Message = fmt.Sprintf("failed to send email, status code: %d", res.StatusCode)
		return request, response
	}

	response.Success = true
	return request, response
}
