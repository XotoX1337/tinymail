package tinymail

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const tplString string = `
<!doctype html>
<html>
  <head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
    <title>Simple Transactional Email</title>
  </head>
  <body style="background-color: #f6f6f6; font-family: sans-serif; -webkit-font-smoothing: antialiased; font-size: 14px; line-height: 1.4; margin: 0; padding: 0; -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%;">
    THIS IS A TEST
  </body>
</html>
`

func TestFromString(t *testing.T) {
	assert := assert.New(t)
	msg := FromString("this is a test")
	want := &message{
		from:        "",
		to:          []string{},
		cc:          []string{},
		bcc:         []string{},
		subject:     "",
		body:        "this is a test",
		attachments: map[string][]byte{},
	}
	assert.Equal(want, msg)
}

func TestFromTemplateString(t *testing.T) {
	assert := assert.New(t)
	msg, err := FromTemplateString(nil, string(tplString))
	assert.NoError(err, "error creating msg from template string")
	want := &message{
		from:        "",
		to:          []string{},
		cc:          []string{},
		bcc:         []string{},
		subject:     "",
		body:        tplString,
		attachments: map[string][]byte{},
	}
	assert.Equal(want, msg)
}

func TestFromTemplateFile(t *testing.T) {
	assert := assert.New(t)
	testFile := "test_template.html"
	err := os.WriteFile(testFile, []byte(tplString), 0644)
	assert.NoErrorf(err, "error writing %s", testFile)
	msg, err := FromTemplateFile(nil, testFile)
	assert.NoErrorf(err, "error creating msg from %s", testFile)
	want := &message{
		from:        "",
		to:          []string{},
		cc:          []string{},
		bcc:         []string{},
		subject:     "",
		body:        tplString,
		attachments: map[string][]byte{},
	}
	assert.Equal(want, msg)
	assert.NoErrorf(os.Remove(testFile), "error deleting %s", testFile)
}

func TestAttach(t *testing.T) {
	assert := assert.New(t)
	fileName := "test_attach"
	fileContent := make([]byte, 1024)
	msg := FromString("TestAttach")
	err := os.WriteFile(fileName, fileContent, 0644)
	assert.NoErrorf(err, "error creating %s", fileName)
	msg.Attach(fileName)
	want := &message{
		from:        "",
		to:          []string{},
		cc:          []string{},
		bcc:         []string{},
		subject:     "",
		body:        "TestAttach",
		attachments: map[string][]byte{fileName: fileContent},
	}
	assert.Equal(want, msg)
	assert.NoError(os.Remove(fileName))
}

func TestAttachMultiple(t *testing.T) {
	assert := assert.New(t)

	nameFile1 := "testFile1"
	contentFile1 := make([]byte, 512)
	nameFile2 := "testFile2"
	contentFile2 := make([]byte, 1024)
	nameFile3 := "testFile3"
	contentFile3 := make([]byte, 2048)

	assert.NoError(os.WriteFile(nameFile1, contentFile1, 0644))
	assert.NoError(os.WriteFile(nameFile2, contentFile2, 0644))
	assert.NoError(os.WriteFile(nameFile3, contentFile3, 0644))

	want := &message{
		from:    "",
		to:      []string{},
		cc:      []string{},
		bcc:     []string{},
		subject: "",
		body:    "TestAttachMultiple",
		attachments: map[string][]byte{
			nameFile1: contentFile1,
			nameFile2: contentFile2,
			nameFile3: contentFile3,
		},
	}

	msg := FromString("TestAttachMultiple")
	msg.Attach(nameFile1, nameFile2, nameFile3)

	assert.Equal(want, msg)

	assert.NoError(os.Remove(nameFile1))
	assert.NoError(os.Remove(nameFile2))
	assert.NoError(os.Remove(nameFile3))
}

func TestSetFrom(t *testing.T) {
	assert := assert.New(t)
	want := &message{
		from:        "test@testing.com",
		to:          []string{},
		cc:          []string{},
		bcc:         []string{},
		subject:     "",
		body:        "TestSetFrom",
		attachments: map[string][]byte{},
	}
	msg := FromString("TestSetFrom")
	msg.SetFrom("test@testing.com")
	assert.Equal(want, msg)
	assert.Equal(want.From(), msg.From())
}

func TestSetTo(t *testing.T) {
	assert := assert.New(t)
	want := &message{
		from:        "",
		to:          []string{"tester@testing.com"},
		cc:          []string{},
		bcc:         []string{},
		subject:     "",
		body:        "TestSetTo",
		attachments: map[string][]byte{},
	}
	msg := FromString("TestSetTo")
	msg.SetTo("tester@testing.com")
	assert.Equal(want, msg)
	assert.Equal([]string{"tester@testing.com"}, msg.To())
	msg.SetTo("testerino@testing.com")
	assert.Equal([]string{"testerino@testing.com"}, msg.To())
	msg.SetTo("testerino@testing.com", "tester@testing.com")
	assert.Equal([]string{"testerino@testing.com", "tester@testing.com"}, msg.To())
}

func TestSetToMultiple(t *testing.T) {
	assert := assert.New(t)
	want := &message{
		from:        "",
		to:          []string{"tester1@testing.com", "tester2@testing.com"},
		cc:          []string{},
		bcc:         []string{},
		subject:     "",
		body:        "TestSetToMultiple",
		attachments: map[string][]byte{},
	}
	msg := FromString("TestSetToMultiple")
	msg.SetTo("tester1@testing.com", "tester2@testing.com")
	assert.Equal(want, msg)

}

func TestSetCC(t *testing.T) {
	assert := assert.New(t)
	want := &message{
		from:        "",
		to:          []string{},
		cc:          []string{"tester@testing.com"},
		bcc:         []string{},
		subject:     "",
		body:        "TestSetCC",
		attachments: map[string][]byte{},
	}
	msg := FromString("TestSetCC")
	msg.SetCC("tester@testing.com")
	assert.Equal(want, msg)
}

func TestSetCCMultiple(t *testing.T) {
	assert := assert.New(t)
	want := &message{
		from:        "",
		to:          []string{},
		cc:          []string{"tester1@testing.com", "tester2@testing.com"},
		bcc:         []string{},
		subject:     "",
		body:        "TestSetCCMultiple",
		attachments: map[string][]byte{},
	}
	msg := FromString("TestSetCCMultiple")
	msg.SetCC("tester1@testing.com", "tester2@testing.com")
	assert.Equal(want, msg)
}

func TestSetBCC(t *testing.T) {
	assert := assert.New(t)
	want := &message{
		from:        "",
		to:          []string{},
		cc:          []string{},
		bcc:         []string{"tester@testing.com"},
		subject:     "",
		body:        "TestSetBCC",
		attachments: map[string][]byte{},
	}
	msg := FromString("TestSetBCC")
	msg.SetBCC("tester@testing.com")
	assert.Equal(want, msg)
}

func TestSetBCCMultiple(t *testing.T) {
	assert := assert.New(t)
	want := &message{
		from:        "",
		to:          []string{},
		cc:          []string{},
		bcc:         []string{"tester1@testing.com", "tester2@testing.com"},
		subject:     "",
		body:        "TestSetBCCMultiple",
		attachments: map[string][]byte{},
	}
	msg := FromString("TestSetBCCMultiple")
	msg.SetBCC("tester1@testing.com", "tester2@testing.com")
	assert.Equal(want, msg)
}

func TestSetSubject(t *testing.T) {
	assert := assert.New(t)
	want := &message{
		from:        "",
		to:          []string{},
		cc:          []string{},
		bcc:         []string{},
		subject:     "Test",
		body:        "TestSetSubject",
		attachments: map[string][]byte{},
	}
	msg := FromString("TestSetSubject")
	msg.SetSubject("Test")
	assert.Equal(want, msg)
}
