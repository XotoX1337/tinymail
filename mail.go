package tinymail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"strings"
)

type Mailer interface {
	Send() error
}

type smtpConfig struct {
	auth     smtp.Auth
	user     string
	password string
	host     string
	addr     string
}

type mail struct {
	message message
	config  *smtpConfig
}

func New(user, password, host string) *mail {
	c := &smtpConfig{
		user:     user,
		password: password,
		host:     host,
		addr:     fmt.Sprintf("%s:%d", host, 587),
	}
	c.auth = smtp.PlainAuth("", c.user, c.password, c.host)
	m := &mail{
		config: c,
	}
	return m
}
func (m *mail) Send() error {
	return smtp.SendMail(m.config.addr, m.config.auth, m.config.user, m.message.To(), m.writeMessage())
}

func (m *mail) SetMessage(msg message) *mail {
	m.message = msg
	return m
}

func (m *mail) writeMessage() []byte {
	msg := m.message
	buf := bytes.NewBuffer(nil)
	buf.WriteString("MIME-Version: 1.0\n")
	withAttachments := len(msg.Attachments()) > 0
	buf.WriteString(fmt.Sprintf("From: %s\n", msg.From()))
	buf.WriteString(fmt.Sprintf("To: %s\n", strings.Join(msg.To(), ",")))
	buf.WriteString(fmt.Sprintf("Subject: %s\n", msg.Subject()))
	if len(msg.CC()) > 0 {
		buf.WriteString(fmt.Sprintf("Cc: %s\n", strings.Join(msg.CC(), ",")))
	}

	if len(msg.BCC()) > 0 {
		buf.WriteString(fmt.Sprintf("Bcc: %s\n", strings.Join(msg.BCC(), ",")))
	}

	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()
	if withAttachments {
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed;\n boundary=%s\n\n", boundary))
		buf.WriteString(fmt.Sprintf("--%s\n", boundary))

	}
	buf.WriteString(fmt.Sprintf("Content-Type: %s\n\n", http.DetectContentType([]byte(msg.Body()))))
	buf.WriteString(msg.Body())
	if withAttachments {
		for k, v := range msg.Attachments() {
			buf.WriteString(fmt.Sprintf("\n--%s\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s\n", http.DetectContentType(v)))
			buf.WriteString("Content-Transfer-Encoding: base64\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\n\n", k))

			b := make([]byte, base64.StdEncoding.EncodedLen(len(v)))
			base64.StdEncoding.Encode(b, v)
			buf.Write(b)
			buf.WriteString(fmt.Sprintf("\n--%s", boundary))
		}

		buf.WriteString("--")
	}

	return buf.Bytes()
}
