// Package tinymail provides a simple and easy to use interface
// to send smtp emails.
package tinymail

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"strings"
)

const DEFAULT_SMTP_PORT int = 587

type Mailer interface {
	Send() error
	SetBoundary(boundary string)
	Boundary() string
	Config() *smtpConfig
}

type MailerOpts struct {
	User     string
	Password string
	Host     string
	Port     int
}

type smtpConfig struct {
	auth     smtp.Auth
	user     string
	password string
	host     string
	addr     string
	port     int
}

type mailer struct {
	message  Message
	boundary string
	config   *smtpConfig
}

func New(opts MailerOpts) (*mailer, error) {
	if err := validateMailerOpts(opts); err != nil {
		return nil, err
	}
	if opts.Port == 0 {
		opts.Port = DEFAULT_SMTP_PORT
	}
	c := &smtpConfig{
		user:     opts.User,
		password: opts.Password,
		host:     opts.Host,
		port:     opts.Port,
		addr:     fmt.Sprintf("%s:%d", opts.Host, opts.Port),
	}
	c.auth = smtp.PlainAuth("", c.user, c.password, c.host)
	m := &mailer{
		config: c,
	}
	return m, nil
}

func validateMailerOpts(opts MailerOpts) error {
	if opts.User == "" {
		return fmt.Errorf("MailerOpts.User is empty")
	}
	if opts.Password == "" {
		return fmt.Errorf("MailerOpts.Password is empty")
	}
	if opts.Host == "" {
		return fmt.Errorf("MailerOpts.Host ist empty")
	}
	return nil
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

func (m *mailer) Config() *smtpConfig {
	return m.config
}

// splits s line by line into RFC5322 compliant chunks
func (m *mailer) chunkLines(s string) string {
	scanner := bufio.NewScanner(strings.NewReader(s))
	var chunkedLines []string
	for scanner.Scan() {
		chunkedLines = append(chunkedLines, m.chunkString(scanner.Text()))
	}

	return strings.Join(chunkedLines, "\n")
}

// chunk e mail into parts of 998 characters due to
// RFC5322 2.1.1 Line Length Limits
func (m *mailer) chunkString(s string) string {
	var chunks []string
	for len(s) > 998 {
		chunks = append(chunks, s[:998])
		s = s[998:]
	}
	chunks = append(chunks, s)
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
	buf.WriteString(m.chunkLines(msg.Body()))
	if withAttachments {
		for k, v := range msg.Attachments() {
			buf.WriteString(fmt.Sprintf("\n--%s\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s\n", http.DetectContentType(v)))
			buf.WriteString("Content-Transfer-Encoding: base64\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\n\n", k))

			b := make([]byte, base64.StdEncoding.EncodedLen(len(v)))
			base64.StdEncoding.Encode(b, v)
			buf.Write([]byte(m.chunkString(string(b))))
			buf.WriteString(fmt.Sprintf("\n--%s", boundary))
		}

		buf.WriteString("--")
	}
	return buf.Bytes()
}
