package tinymail

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteMessage(t *testing.T) {
	assert := assert.New(t)

	want := `MIME-Version: 1.0
From: test@tinymail.test
To: test.to@tinymail.test
Subject: TestWriteMessage
Cc: test.cc@tinymail.test
Bcc: test.bcc@tinymail.test
Content-Type: text/plain; charset=utf-8

this is a test`

	mailer := New("", "", "")
	msg := FromText("this is a test")
	msg.SetFrom("test@tinymail.test")
	msg.SetTo("test.to@tinymail.test")
	msg.SetSubject("TestWriteMessage")
	msg.SetCC("test.cc@tinymail.test")
	msg.SetBCC("test.bcc@tinymail.test")
	mailer.SetMessage(msg)
	assert.Equal(want, string(mailer.writeMessage()))
}

func TestWriteMessageUrgent(t *testing.T) {
	assert := assert.New(t)

	want := `MIME-Version: 1.0
From: test@tinymail.test
To: test.to@tinymail.test
Subject: TestWriteMessageUrgent
Cc: test.cc@tinymail.test
Bcc: test.bcc@tinymail.test
Priority: urgent
Content-Type: text/plain; charset=utf-8

this is a test`

	mailer := New("", "", "")
	msg := FromText("this is a test")
	msg.SetFrom("test@tinymail.test")
	msg.SetTo("test.to@tinymail.test")
	msg.SetSubject("TestWriteMessageUrgent")
	msg.SetCC("test.cc@tinymail.test")
	msg.SetBCC("test.bcc@tinymail.test")
	msg.SetUrgentPriority()
	mailer.SetMessage(msg)
	assert.Equal(want, string(mailer.writeMessage()))
}

func TestWriteMessageAttach(t *testing.T) {
	assert := assert.New(t)

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

	mailer := New("", "", "")
	mailer.SetBoundary("7b7f6c9583aae2870247062aac5ca1bc1610b22b627ae2c5366bb1394ed0")

	msg := FromText("this is a test")
	msg.SetFrom("test@tinymail.test")
	msg.SetTo("test.to@tinymail.test")
	msg.SetSubject("TestWriteMessageAttach")
	msg.SetCC("test.cc@tinymail.test")
	msg.SetBCC("test.bcc@tinymail.test")
	msg.Attach("TestWriteMessageAttach")

	mailer.SetMessage(msg)

	assert.Equal(want, string(mailer.writeMessage()))
	assert.NoError(os.Remove("TestWriteMessageAttach"))
}
