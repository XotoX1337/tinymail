package tinymail

import (
	"bytes"
	"io"
	"mime"
	"mime/multipart"
	"net/mail"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Bug 3: a message with multiple attachments must produce a valid
// multipart/mixed body that parses back into one text part plus one part
// per attachment.
func TestWriteMessageMultipleAttachments(t *testing.T) {
	test := assert.New(t)

	dir := t.TempDir()
	files := map[string]int{
		"first.bin":  128,
		"second.bin": 256,
		"third.bin":  512,
	}
	mailer, err := New(VALID_MAILER_OPTS)
	test.NoError(err)

	msg := FromString("body text")
	msg.SetFrom("from@tinymail.test")
	msg.SetTo("to@tinymail.test")
	msg.SetSubject("attachments")
	for name, size := range files {
		path := filepath.Join(dir, name)
		test.NoError(os.WriteFile(path, make([]byte, size), 0644))
		test.NoError(msg.Attach(path))
	}
	mailer.SetMessage(msg)

	raw := mailer.writeMessage()

	parsed, err := mail.ReadMessage(bytes.NewReader(raw))
	require.NoError(t, err)

	mediaType, params, err := mime.ParseMediaType(parsed.Header.Get("Content-Type"))
	require.NoError(t, err)
	test.Equal("multipart/mixed", mediaType)
	test.NotEmpty(params["boundary"])

	mr := multipart.NewReader(parsed.Body, params["boundary"])
	var textParts, attachmentParts int
	seen := map[string]bool{}
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		if part.FileName() != "" {
			attachmentParts++
			seen[part.FileName()] = true
		} else {
			textParts++
		}
		io.Copy(io.Discard, part)
	}

	test.Equal(1, textParts, "expected exactly one body part")
	test.Equal(len(files), attachmentParts, "expected one part per attachment")
	for name := range files {
		test.True(seen[name], "attachment %q missing from message", name)
	}
}
