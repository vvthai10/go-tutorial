package mail

type IEmailSender interface {
	Send(request *Request) (*Request, *Response)
}

type EmailSender struct {
	IEmailSender
}

func (s *EmailSender) SetupSendGrid(cfg SendGridConfig) error {
	s.IEmailSender = NewSendGridSender(cfg)
	return nil
}

func (s *EmailSender) SetupSMTP(cfg SMTPConfig) error {
	s.IEmailSender = NewSMTPSender(cfg)
	return nil
}

func NewEmailSender(emailSender IEmailSender) EmailSender {
	return EmailSender{
		IEmailSender: emailSender,
	}
}
