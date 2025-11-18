package tinymail

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var VALID_MAILER_OPTS = MailerOpts{
	User:     "test",
	Password: "secret",
	Host:     "test.com",
}

var MISSING_USER_MAILER_OPTS = MailerOpts{
	Password: "secret",
	Host:     "test.com",
}

var MISSING_PASSWORD_MAILER_OPTS = MailerOpts{
	User: "test",
	Host: "test.com",
}

var MISSING_HOST_MAILER_OPTS = MailerOpts{
	User:     "test",
	Password: "secret",
}

var CUSTOM_PORT_MAILER_OPTS = MailerOpts{
	User:     "test",
	Password: "secret",
	Host:     "test.com",
	Port:     123,
}

func TestWriteMessage(t *testing.T) {
	test := assert.New(t)

	want := `MIME-Version: 1.0
From: test@tinymail.test
To: test.to@tinymail.test
Subject: TestWriteMessage
Cc: test.cc@tinymail.test
Bcc: test.bcc@tinymail.test
Content-Type: text/plain; charset=utf-8

this is a test`

	mailer, err := New(VALID_MAILER_OPTS)
	test.NoError(err)

	msg := FromString("this is a test")
	msg.SetFrom("test@tinymail.test")
	msg.SetTo("test.to@tinymail.test")
	msg.SetSubject("TestWriteMessage")
	msg.SetCC("test.cc@tinymail.test")
	msg.SetBCC("test.bcc@tinymail.test")

	mailer.SetMessage(msg)

	test.Equal(want, string(mailer.writeMessage()))
}

func TestWriteMessageUrgent(t *testing.T) {
	test := assert.New(t)

	want := `MIME-Version: 1.0
From: test@tinymail.test
To: test.to@tinymail.test
Subject: TestWriteMessageUrgent
Cc: test.cc@tinymail.test
Bcc: test.bcc@tinymail.test
Priority: urgent
Content-Type: text/plain; charset=utf-8

this is a test`

	mailer, err := New(VALID_MAILER_OPTS)
	test.NoError(err)

	msg := FromString("this is a test")
	msg.SetFrom("test@tinymail.test")
	msg.SetTo("test.to@tinymail.test")
	msg.SetSubject("TestWriteMessageUrgent")
	msg.SetCC("test.cc@tinymail.test")
	msg.SetBCC("test.bcc@tinymail.test")
	msg.SetUrgentPriority()

	mailer.SetMessage(msg)
	test.Equal(want, string(mailer.writeMessage()))
}

func TestWriteMessageAttach(t *testing.T) {
	test := assert.New(t)

	want := `MIME-Version: 1.0
From: test@tinymail.test
To: test.to@tinymail.test
Subject: TestWriteMessageAttach
Cc: test.cc@tinymail.test
Bcc: test.bcc@tinymail.test
Content-Type: multipart/mixed;
 boundary=7b7f6c9583aae2870247062aac5ca1bc1610b22b627ae2c5366bb1394ed0

--7b7f6c9583aae2870247062aac5ca1bc1610b22b627ae2c5366bb1394ed0
Content-Type: text/plain; charset=utf-8

this is a test
--7b7f6c9583aae2870247062aac5ca1bc1610b22b627ae2c5366bb1394ed0
Content-Type: application/octet-stream
Content-Transfer-Encoding: base64
Content-Disposition: attachment; filename=TestWriteMessageAttach

AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=
--7b7f6c9583aae2870247062aac5ca1bc1610b22b627ae2c5366bb1394ed0--`

	os.WriteFile("TestWriteMessageAttach", make([]byte, 512), 0644)

	mailer, err := New(VALID_MAILER_OPTS)
	test.NoError(err)

	mailer.SetBoundary("7b7f6c9583aae2870247062aac5ca1bc1610b22b627ae2c5366bb1394ed0")

	msg := FromString("this is a test")
	msg.SetFrom("test@tinymail.test")
	msg.SetTo("test.to@tinymail.test")
	msg.SetSubject("TestWriteMessageAttach")
	msg.SetCC("test.cc@tinymail.test")
	msg.SetBCC("test.bcc@tinymail.test")
	msg.Attach("TestWriteMessageAttach")

	mailer.SetMessage(msg)

	test.Equal(want, string(mailer.writeMessage()))
	test.NoError(os.Remove("TestWriteMessageAttach"))
}

func TestDefaultPortOpt(t *testing.T) {
	test := assert.New(t)

	mailer, err := New(VALID_MAILER_OPTS)
	test.NoError(err)

	config := mailer.Config()
	test.Equal(DEFAULT_SMTP_PORT, config.port)
}

func TestMissingUsernameInMailerOpts(t *testing.T) {
	test := assert.New(t)

	mailer, err := New(MISSING_USER_MAILER_OPTS)

	test.Error(err)
	test.Nil(mailer)
}

func TestMissingPasswordInMailerOpts(t *testing.T) {
	test := assert.New(t)

	mailer, err := New(MISSING_PASSWORD_MAILER_OPTS)

	test.Error(err)
	test.Nil(mailer)
}

func TestMissingHostInMailerOpts(t *testing.T) {
	test := assert.New(t)

	mailer, err := New(MISSING_HOST_MAILER_OPTS)

	test.Error(err)
	test.Nil(mailer)
}

func TestCustomPortInMailerOpts(t *testing.T) {
	test := assert.New(t)

	mailer, err := New(CUSTOM_PORT_MAILER_OPTS)
	test.NoError(err)

	config := mailer.Config()
	test.Equal(123, config.port)
}
