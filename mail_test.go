package tinymail

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var user = "95636f24a5e0c5"
var pwd = "a6874bb3a8de46"
var host = "sandbox.smtp.mailtrap.io"
var from = "info@tinymail.com"

func TestSend(t *testing.T) {
	assert := assert.New(t)
	mailer := New(user, pwd, host)
	msg := FromText("this is a test")
	msg.SetFrom(from)
	msg.SetTo("tester.mister@testing.com")
	msg.SetSubject("tinymail")
	msg.SetCC("tester.mister@testing.com")
	msg.SetBCC("tester.mister@testing.com")
	err := mailer.SetMessage(*msg).Send()
	assert.NoError(err, "error while sending email")
}
