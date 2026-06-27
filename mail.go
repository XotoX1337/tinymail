// Package tinymail provides a simple and easy to use interface
// to send smtp emails.
package tinymail

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"strings"
)

const DEFAULT_SMTP_PORT int = 587

// crlf is the RFC 5322 line separator used throughout the serialized message.
const crlf = "\r\n"

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
	TLS      bool
}

type smtpConfig struct {
	auth     smtp.Auth
	user     string
	password string
	host     string
	addr     string
	port     int
	tls      bool
}

type mailer struct {
	message  Message
	boundary string
	config   *smtpConfig
}

type smtpLoginAuth struct {
	username, password string
}

// New returns a new Mailer instance
//
// Returns an error, if opts could not be validated
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
		tls:      opts.TLS,
	}
	c.auth = smtp.PlainAuth("", c.user, c.password, c.host)
	m := &mailer{
		config: c,
	}
	return m, nil
}

// validateMailerOpts returns an error if
//
//   - [MailerOpts.User] is empty
//   - [MailerOpts.Password] is empty
//   - [MailerOpts.Host] is empty
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

// Send sends the message with [net/smtp.SendMail]
func (m *mailer) Send() error {
	if m.config.tls {
		return m.sendTLS()
	} else {
		return m.sendPlain()
	}

}

func (m *mailer) sendTLS() error {
	tlsConfig := &tls.Config{
		ServerName: m.config.host,
	}
	c, err := smtp.Dial(m.config.addr)
	if err != nil {
		return err
	}

	if err = c.StartTLS(tlsConfig); err != nil {
		return err
	}
	if ok, auths := c.Extension("AUTH"); ok {
		if strings.Contains(auths, "LOGIN") &&
			!strings.Contains(auths, "PLAIN") {
			m.config.auth = loginAuth(m.config.user, m.config.password)
		}
	} else {
		return errors.New("no authentication method found")
	}

	if err = c.Auth(m.config.auth); err != nil {
		return err
	}

	if err = c.Mail(m.config.user); err != nil {
		return err
	}

	for _, rcpt := range m.recipients() {
		if err = c.Rcpt(rcpt); err != nil {
			return err
		}
	}
	writer, err := c.Data()
	if err != nil {
		return err
	}
	_, err = writer.Write(m.writeMessage())
	if err != nil {
		return err
	}
	err = writer.Close()
	if err != nil {
		return err
	}

	return c.Quit()
}

func (m *mailer) sendPlain() error {
	return smtp.SendMail(m.config.addr, m.config.auth, m.config.user, m.recipients(), m.writeMessage())
}

// recipients returns every envelope recipient: To, CC and BCC. All of them
// need a RCPT TO, otherwise CC/BCC recipients never receive the message.
func (m *mailer) recipients() []string {
	msg := m.message
	rcpts := make([]string, 0, len(msg.To())+len(msg.CC())+len(msg.BCC()))
	rcpts = append(rcpts, msg.To()...)
	rcpts = append(rcpts, msg.CC()...)
	rcpts = append(rcpts, msg.BCC()...)
	return rcpts
}

// SetMessage sets the message
func (m *mailer) SetMessage(msg Message) *mailer {
	m.message = msg
	return m
}

// SetBoundary sets the boundary string
func (m *mailer) SetBoundary(boundary string) *mailer {
	m.boundary = boundary
	return m
}

// Boundary returns the boundary string
func (m *mailer) Boundary() string {
	return m.boundary
}

// Config returns the SMTP Config
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

	return strings.Join(chunkedLines, crlf)
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
	return strings.Join(chunks, crlf)
}

// sanitizeHeader strips CR and LF from a header value to prevent header
// injection (an attacker-controlled value breaking out into new headers).
func sanitizeHeader(s string) string {
	return strings.NewReplacer("\r", "", "\n", "").Replace(s)
}

// writeMessage serializes the message into RFC 5322 / MIME wire format.
func (m *mailer) writeMessage() []byte {
	msg := m.message
	buf := bytes.NewBuffer(nil)
	withAttachments := len(msg.Attachments()) > 0

	buf.WriteString("MIME-Version: 1.0" + crlf)
	buf.WriteString(fmt.Sprintf("From: %s%s", sanitizeHeader(msg.From()), crlf))
	buf.WriteString(fmt.Sprintf("To: %s%s", sanitizeHeader(strings.Join(msg.To(), ",")), crlf))
	buf.WriteString(fmt.Sprintf("Subject: %s%s", sanitizeHeader(msg.Subject()), crlf))
	if len(msg.CC()) > 0 {
		buf.WriteString(fmt.Sprintf("Cc: %s%s", sanitizeHeader(strings.Join(msg.CC(), ",")), crlf))
	}
	// Bcc is intentionally omitted from the headers: blind recipients must not
	// be visible to anyone. They are delivered via the SMTP envelope only.
	if len(msg.Priority()) > 0 {
		buf.WriteString(fmt.Sprintf("Priority: %s%s", sanitizeHeader(msg.Priority()), crlf))
	}

	var boundary string
	if withAttachments {
		boundary = m.boundary
		if boundary == "" {
			boundary = multipart.NewWriter(buf).Boundary()
		}
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed;%s boundary=%s%s%s", crlf, boundary, crlf, crlf))
		buf.WriteString(fmt.Sprintf("--%s%s", boundary, crlf))
	}

	buf.WriteString(fmt.Sprintf("Content-Type: %s%s%s", http.DetectContentType([]byte(msg.Body())), crlf, crlf))
	buf.WriteString(m.chunkLines(msg.Body()))

	if withAttachments {
		for name, content := range msg.Attachments() {
			buf.WriteString(fmt.Sprintf("%s--%s%s", crlf, boundary, crlf))
			buf.WriteString(fmt.Sprintf("Content-Type: %s%s", http.DetectContentType(content), crlf))
			buf.WriteString("Content-Transfer-Encoding: base64" + crlf)
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s%s%s", sanitizeHeader(name), crlf, crlf))

			b := make([]byte, base64.StdEncoding.EncodedLen(len(content)))
			base64.StdEncoding.Encode(b, content)
			buf.Write([]byte(m.chunkString(string(b))))
		}
		buf.WriteString(fmt.Sprintf("%s--%s--", crlf, boundary))
	}
	return buf.Bytes()
}

func loginAuth(username, password string) smtp.Auth {
	return &smtpLoginAuth{username, password}
}

func (a *smtpLoginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	if !server.TLS {
		return "", nil, errors.New("tinymail: refusing to send LOGIN credentials over an unencrypted connection")
	}
	return "LOGIN", nil, nil
}

func (a *smtpLoginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if !more {
		return nil, nil
	}
	// Servers phrase the challenge differently ("Username:", "username", ...);
	// normalize before matching.
	switch strings.ToLower(strings.TrimRight(strings.TrimSpace(string(fromServer)), ":")) {
	case "username":
		return []byte(a.username), nil
	case "password":
		return []byte(a.password), nil
	default:
		return nil, fmt.Errorf("tinymail: unexpected LOGIN challenge %q", fromServer)
	}
}
