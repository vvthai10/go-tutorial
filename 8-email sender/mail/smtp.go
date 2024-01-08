package mail

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

type SMTPConfig struct {
	Name          string // name of sender
	Email         string // email of sender
	Password      string // password of app email
	AuthAddress   string // smtp auth address
	ServerAddress string // smtp server address
}

type smtpSender struct {
	config SMTPConfig
}

func NewSMTPSender(cfg SMTPConfig) IEmailSender {
	return &smtpSender{
		config: cfg,
	}
}

func (sender *smtpSender) Send(request *Request) (*Request, *Response) {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.config.Name, sender.config.Email)
	e.Subject = request.Subject
	e.HTML = []byte(request.Content)
	e.To = request.To
	e.Cc = request.Cc
	e.Bcc = request.Bcc

	response := &Response{}

	for _, f := range request.AttachFiles {
		_, err := e.AttachFile(f)
		if err != nil {
			response.Message = fmt.Sprintf("failed to attach file %s: %s", f, err.Error())
			return request, response
		}
	}

	smtpAuth := smtp.PlainAuth("", sender.config.Email, sender.config.Password, sender.config.AuthAddress)
	err := e.Send(sender.config.ServerAddress, smtpAuth)
	if err != nil {
		response.Message = err.Error()
		return request, response
	}
	response.Success = true
	return request, response
}
