package mail

import (
	"net/smtp"
	"os"

	"github.com/spf13/viper"
)

func SendEmail(to, message string) error {
	from := os.Getenv("EMAIL_ADDRESS")
	password := os.Getenv("EMAIL_PASSWORD")
	host := viper.GetString("mail.host")
	port := viper.GetString("mail.port")

	subject := viper.GetString("mail.subject")

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		message

	receiver := []string{to}
	sender := host + ":" + port

	auth := smtp.PlainAuth("", from, password, host)

	return smtp.SendMail(sender, auth, from, receiver, []byte(msg))
}
