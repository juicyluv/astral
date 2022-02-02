package mail

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/spf13/viper"
)

var (
	MimeHTML = "text/html"
)

func SendEmail(to, subject, mimeType, body string) error {
	from := os.Getenv("EMAIL_ADDRESS")
	password := os.Getenv("EMAIL_PASSWORD")
	host := viper.GetString("mail.host")
	port := viper.GetString("mail.port")

	header := make(map[string]string)
	header["From"] = from

	header["To"] = to
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = fmt.Sprintf("%s; charset=\"utf-8\"", mimeType)
	header["Content-Transfer-Encoding"] = "quoted-printable"
	header["Content-Disposition"] = "inline"

	message := ""

	for key, value := range header {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	message += "\r\n" + body

	receiver := []string{to}
	sender := host + ":" + port

	auth := smtp.PlainAuth("", from, password, host)

	return smtp.SendMail(sender, auth, from, receiver, []byte(message))
}
