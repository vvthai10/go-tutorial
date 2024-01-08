package main

import (
	"fmt"

	mailSender "github.com/vvthai10/email/mail"
)

func main() {
	// sender := EmailSender{}
	var sender mailSender.EmailSender
	// sender.SetupSMTP(mailSender.SMTPConfig{
	// 	Name:          "Demo Sender",
	// 	Email:         "vuvanthai1410@gmail.com",
	// 	Password:      "evswoithrfimicln",
	// 	AuthAddress:   "smtp.gmail.com",
	// 	ServerAddress: "smtp.gmail.com:587",
	// })
	sender.SetupSendGrid(mailSender.SendGridConfig{
		Name:   "Demo Sender",
		Email:  "vuvanthai1410@gmail.com",
		ApiKey: "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",
	})

	subject := "A test email"
	content := `
		<h1>Hello world</h1>
		<p>This email to test sender</p>
	`
	to := []string{"vuvanthai1410@gmail.com"}
	_, res := sender.Send(&mailSender.Request{
		Subject: subject,
		Content: content,
		To:      to,
	})
	fmt.Println(res)
}
