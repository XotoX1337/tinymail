package tinymail

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
)

type Message interface {
	FromText(text string) *message
	FromTemplateString(tpl string) *message
	FromTemplateFile(file string) *message
	SetFrom(from string)
	SetTo(to []string)
	SetCC(cc []string)
	SetBCC(bcc []string)
	SetSubject(s string)
	Attach(files []string)
}
type message struct {
	from        string
	to          []string
	cc          []string
	bcc         []string
	subject     string
	body        string
	attachments map[string][]byte
}

func (m *message) SetFrom(from string) {
	m.from = from
}
func (m *message) From() string {
	return m.from
}

func (m *message) SetTo(to ...string) {
	m.to = to
}

func (m *message) To() []string {
	return m.to
}

func (m *message) SetCC(cc ...string) {
	m.cc = cc
}

func (m *message) CC() []string {
	return m.cc
}

func (m *message) SetBCC(bcc ...string) {
	m.bcc = bcc
}

func (m *message) BCC() []string {
	return m.bcc
}

func (m *message) SetSubject(subject string) {
	m.subject = subject
}

func (m *message) Subject() string {
	return m.subject
}

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

func (m *message) Attachments() map[string][]byte {
	return m.attachments
}

func (m *message) Body() string {
	return m.body
}

func FromText(text string) *message {
	m := new()
	m.body = text
	return m
}

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
