// Package tinymail provides a simple and easy to use interface
// to send smtp emails.
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
	SetBoundary(boundary string)
	Boundary() string
}

type smtpConfig struct {
	auth     smtp.Auth
	user     string
	password string
	host     string
	addr     string
}

type mailer struct {
	message  Message
	boundary string
	config   *smtpConfig
}

func New(user, password, host string) *mailer {
	c := &smtpConfig{
		user:     user,
		password: password,
		host:     host,
		addr:     fmt.Sprintf("%s:%d", host, 587),
	}
	c.auth = smtp.PlainAuth("", c.user, c.password, c.host)
	m := &mailer{
		config: c,
	}
	return m
}
func (m *mailer) Send() error {
	return smtp.SendMail(m.config.addr, m.config.auth, m.config.user, m.message.To(), m.writeMessage())
}

func (m *mailer) SetMessage(msg Message) *mailer {
	m.message = msg
	return m
}

func (m *mailer) SetBoundary(boundary string) *mailer {
	m.boundary = boundary
	return m
}

func (m *mailer) Boundary() string {
	return m.boundary
}

// chunk e mail into parts of 998 characters due to
// RFC5322 2.1.1 Line Length Limits
func (m *mailer) chunkMessage(message string) string {
	var chunks []string
	for len(message) > 998 {
		chunks = append(chunks, message[:998])
		message = message[998:]
	}
	chunks = append(chunks, message)
	return strings.Join(chunks, "\n")
}

func (m *mailer) writeMessage() []byte {
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
	if len(msg.Priority()) > 0 {
		buf.WriteString(fmt.Sprintf("Priority: %s\n", msg.Priority()))
	}

	writer := multipart.NewWriter(buf)
	var boundary string
	if len(m.boundary) > 0 {
		boundary = m.boundary
	} else {
		boundary = writer.Boundary()
	}

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
	return []byte(m.chunkMessage(buf.String()))
}
