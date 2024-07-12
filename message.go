package tinymail

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
)

type Message interface {
	SetFrom(from string)
	From() string
	SetTo(to ...string)
	To() []string
	SetCC(cc ...string)
	CC() []string
	SetBCC(bcc ...string)
	BCC() []string
	SetSubject(s string)
	Subject() string
	Attach(files ...string) error
	Attachments() map[string][]byte
	Body() string
	SetUrgentPriority()
	SetNonUrgentPriority()
	SetNormalPriority()
	Priority() string
}
type message struct {
	from        string
	to          []string
	cc          []string
	bcc         []string
	subject     string
	body        string
	priority    string
	attachments map[string][]byte
}

// SetFrom sets the sender email.
func (m *message) SetFrom(from string) {
	m.from = from
}

// From returns the sender email.
func (m *message) From() string {
	return m.from
}

// SetTo sets the receiver addresses.
func (m *message) SetTo(to ...string) {
	m.to = to
}

// To returns receivers.
func (m *message) To() []string {
	return m.to
}

// SetCC sets the CC recipients.
func (m *message) SetCC(cc ...string) {
	m.cc = cc
}

// CC returns the CC recipients.
func (m *message) CC() []string {
	return m.cc
}

// SetBCC sets the BCC recipients.
func (m *message) SetBCC(bcc ...string) {
	m.bcc = bcc
}

// BCC returns the BCC recipients.
func (m *message) BCC() []string {
	return m.bcc
}

// SetSubject sets the subject.
func (m *message) SetSubject(subject string) {
	m.subject = subject
}

// Subject returns the Subject.
func (m *message) Subject() string {
	return m.subject
}

// Attach attaches files
//
// Returns an error if one of files could not be read.
func (m *message) Attach(files ...string) error {
	for _, file := range files {
		b, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		_, fileName := filepath.Split(file)
		m.attachments[fileName] = b
	}
	return nil
}

// Attachments returns the attachments.
func (m *message) Attachments() map[string][]byte {
	return m.attachments
}

// Body returns the body.
func (m *message) Body() string {
	return m.body
}

// SetNormalPriority sets the email priority to 'normal'.
func (m *message) SetNormalPriority() {
	m.priority = "normal"
}

// SetUrgentPriority sets the email priority to 'urgent'.
func (m *message) SetUrgentPriority() {
	m.priority = "urgent"
}

// SetNonUrgentPriority sets the email priority to 'non-urgent'.
func (m *message) SetNonUrgentPriority() {
	m.priority = "non-urgent"
}

// Priority returns the priority.
func (m *message) Priority() string {
	return m.priority
}

// FromString creates a new message with content from given string.
func FromString(str string) *message {
	m := new()
	m.body = str
	return m
}

// FromTemplateString creates a new message with content from parsed template string.
//
// Returns an error if the template string could not be parsed.
func FromTemplateString(data any, tpl string) (*message, error) {
	buff := bytes.Buffer{}
	template, _ := template.New("tinymail").Parse(tpl)
	err := template.Execute(&buff, data)
	if err != nil {
		return nil, err
	}
	m := new()
	m.body = buff.String()
	return m, nil
}

// FromTemplateFile creates a new message with content from parsed template file.
//
// Returns an error if the template file could not be parsed.
func FromTemplateFile(data any, filenames ...string) (*message, error) {
	buff := bytes.Buffer{}
	template, _ := template.ParseFiles(filenames...)
	err := template.Execute(&buff, data)
	if err != nil {
		return nil, err
	}
	m := new()
	m.body = buff.String()
	return m, nil
}

// new creates a new empty message.
func new() *message {
	return &message{
		from:        "",
		to:          []string{},
		cc:          []string{},
		bcc:         []string{},
		subject:     "",
		body:        "",
		attachments: map[string][]byte{},
	}
}
