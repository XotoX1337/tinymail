package tinymail

import (
	"net/mail"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestMailer wires a mailer to the given mock server. Host must be a
// loopback literal so PlainAuth allows the unencrypted test connection.
func newTestMailer(t *testing.T, srv *mockSMTP) *mailer {
	t.Helper()
	m, err := New(MailerOpts{
		User:     "test",
		Password: "secret",
		Host:     "127.0.0.1",
		Port:     srv.port(),
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return m
}

// Bug 1: every To/CC/BCC recipient must receive an envelope RCPT TO,
// otherwise CC and BCC recipients silently never get the mail.
func TestSendDeliversToAllRecipients(t *testing.T) {
	test := assert.New(t)
	srv := newMockSMTP(t)

	msg := FromString("hello")
	msg.SetFrom("from@tinymail.test")
	msg.SetTo("to1@tinymail.test", "to2@tinymail.test")
	msg.SetCC("cc@tinymail.test")
	msg.SetBCC("bcc@tinymail.test")
	msg.SetSubject("recipients")

	test.NoError(newTestMailer(t, srv).SetMessage(msg).Send())

	rcpts := srv.recipients()
	test.ElementsMatch([]string{
		"to1@tinymail.test",
		"to2@tinymail.test",
		"cc@tinymail.test",
		"bcc@tinymail.test",
	}, rcpts)
}

// Bug 2: CC stays visible in the header, but BCC must never appear in the
// message body delivered to recipients.
func TestSendDoesNotLeakBCC(t *testing.T) {
	test := assert.New(t)
	srv := newMockSMTP(t)

	msg := FromString("hello")
	msg.SetFrom("from@tinymail.test")
	msg.SetTo("to@tinymail.test")
	msg.SetCC("cc@tinymail.test")
	msg.SetBCC("secret@tinymail.test")
	msg.SetSubject("bcc")

	test.NoError(newTestMailer(t, srv).SetMessage(msg).Send())

	data := srv.body()
	test.Contains(data, "Cc: cc@tinymail.test")
	test.NotContains(data, "Bcc:")
	test.NotContains(data, "secret@tinymail.test")
}

// Bug 4: the serialized message must use CRLF line endings per RFC 5322.
// This is asserted on writeMessage() directly because net/smtp's DotWriter
// rewrites bare LF to CRLF on the wire, masking the bug at the mock level.
func TestWriteMessageUsesCRLF(t *testing.T) {
	test := assert.New(t)

	mailer, err := New(VALID_MAILER_OPTS)
	test.NoError(err)

	msg := FromString("line one\nline two")
	msg.SetFrom("from@tinymail.test")
	msg.SetTo("to@tinymail.test")
	msg.SetSubject("crlf")
	mailer.SetMessage(msg)

	data := string(mailer.writeMessage())
	test.Contains(data, "\r\n")
	// No bare LF: stripping all CRLF pairs must leave no stray '\n'.
	test.NotContains(strings.ReplaceAll(data, "\r\n", ""), "\n")
}

// Bug 5: header fields must not allow injecting extra headers via CR/LF.
// A CRLF smuggled into the subject must not produce a real Bcc header or an
// extra envelope recipient.
func TestSendRejectsHeaderInjection(t *testing.T) {
	test := assert.New(t)
	srv := newMockSMTP(t)

	msg := FromString("hello")
	msg.SetFrom("from@tinymail.test")
	msg.SetTo("to@tinymail.test")
	msg.SetSubject("hi\r\nBcc: victim@tinymail.test")

	test.NoError(newTestMailer(t, srv).SetMessage(msg).Send())

	parsed, err := mail.ReadMessage(strings.NewReader(srv.body()))
	require.NoError(t, err)
	// The injected line must not have become a real header.
	test.Empty(parsed.Header.Get("Bcc"))
	test.NotContains(srv.recipients(), "victim@tinymail.test")
}
