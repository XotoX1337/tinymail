package tinymail

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const tplFile string = "test_template.html"

func TestFromText(t *testing.T) {
	assert := assert.New(t)
	msg := FromText("this is a text")
	assert.Equal("this is a text", msg.Body(), "msg body should be 'this is a test'")
}

func TestFromTemplateString(t *testing.T) {
	assert := assert.New(t)
	tplString, err := os.ReadFile(tplFile)
	assert.NoError(err, "error reading tpl file")
	msg, err := FromTemplateString(nil, string(tplString))
	assert.NoError(err, "error creating msg from template string")
	assert.Equal(string(tplString), msg.Body())
}

func TestFromTemplateFile(t *testing.T) {
	assert := assert.New(t)
	tplString, err := os.ReadFile(tplFile)
	assert.NoError(err, "error reading tpl file")
	msg, err := FromTemplateFile(nil, tplFile)
	assert.NoError(err, "error creating msg from template file")
	assert.Equal(string(tplString), msg.Body())
}

func TestAttach(t *testing.T) {
	assert := assert.New(t)
	msg := FromText("TestAttach")
	msg.Attach(tplFile)
	for file, bytes := range msg.Attachments() {
		b, err := os.ReadFile(file)
		assert.NoError(err, "error reading file")
		assert.Equal(b, bytes, "file contents are not equal")
		assert.Equal(file, tplFile, "file name is not equal")
	}
}

func TestSetFrom(t *testing.T) {
	assert := assert.New(t)
	msg := FromText("TestSetFrom")
	msg.SetFrom("test@testing.com")
	assert.Equal("test@testing.com", msg.From())
	msg.SetFrom("test@changed.com")
	assert.Equal("test@changed.com", msg.From())
}

func TestSetTo(t *testing.T) {
	assert := assert.New(t)
	msg := FromText("TestSetTo")
	msg.SetTo("tester@testing.com")
	assert.Equal([]string{"tester@testing.com"}, msg.To())
	msg.SetTo("testerino@testing.com")
	assert.Equal([]string{"testerino@testing.com"}, msg.To())
	msg.SetTo("testerino@testing.com", "tester@testing.com")
	assert.Equal([]string{"testerino@testing.com", "tester@testing.com"}, msg.To())
}

func TestSetCC(t *testing.T) {
	assert := assert.New(t)
	msg := FromText("TestSetCC")
	msg.SetCC("tester@testing.com")
	assert.Equal([]string{"tester@testing.com"}, msg.CC())
	msg.SetCC("testerino@testing.com")
	assert.Equal([]string{"testerino@testing.com"}, msg.CC())
	msg.SetCC("testerino@testing.com", "tester@testing.com")
	assert.Equal([]string{"testerino@testing.com", "tester@testing.com"}, msg.CC())
}

func TestSetBCC(t *testing.T) {
	assert := assert.New(t)
	msg := FromText("TestSetBCC")
	msg.SetBCC("tester@testing.com")
	assert.Equal([]string{"tester@testing.com"}, msg.BCC())
	msg.SetBCC("testerino@testing.com")
	assert.Equal([]string{"testerino@testing.com"}, msg.BCC())
	msg.SetBCC("testerino@testing.com", "tester@testing.com")
	assert.Equal([]string{"testerino@testing.com", "tester@testing.com"}, msg.BCC())
}

func TestSetSubject(t *testing.T) {
	assert := assert.New(t)
	msg := FromText("TestSetSubject")
	msg.SetSubject("test")
	assert.Equal("test", msg.Subject())
	msg.SetSubject("changed")
	assert.Equal("changed", msg.Subject())
}
