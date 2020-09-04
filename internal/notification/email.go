package notification

import (
	"bytes"
	"fmt"
	"mime/quotedprintable"
	"net/smtp"
)

// SMTPEmailSender allows to send email message easitly
type SMTPEmailSender struct {
	Host     string
	Username string
	Password string
	Secure   string
	Port     string
}

// Send method used to send message actually.
func (e *SMTPEmailSender) Send(to, from, subject, message string) (bool, error) {
	header := make(map[string]string)
	header["From"] = from

	header["To"] = to
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = fmt.Sprintf("%s; charset=\"utf-8\"", "text/html")
	header["Content-Transfer-Encoding"] = "quoted-printable"
	header["Content-Disposition"] = "inline"

	body := ""

	for key, value := range header {
		body += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	var encodedMessage bytes.Buffer

	finalMessage := quotedprintable.NewWriter(&encodedMessage)
	finalMessage.Write([]byte(message))
	finalMessage.Close()

	body += "\r\n" + encodedMessage.String()
	if err := smtp.SendMail(
		e.Host+":"+e.Port,
		smtp.PlainAuth("", e.Username, e.Password, e.Host),
		from,
		[]string{to},
		[]byte(body)); err != nil {

		return false, err
	}
	return true, nil

}
