package tinymail

import (
	"net/smtp"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Bug 6: the LOGIN auth mechanism must refuse to send credentials over an
// unencrypted connection, and must answer the server's challenges robustly
// regardless of casing or trailing punctuation.
func TestLoginAuthRefusesUnencrypted(t *testing.T) {
	test := assert.New(t)
	a := loginAuth("user", "pass")

	_, _, err := a.Start(&smtp.ServerInfo{Name: "mail.example.com", TLS: false})
	test.Error(err, "LOGIN auth must not start on an unencrypted connection")
}

func TestLoginAuthStartsOverTLS(t *testing.T) {
	test := assert.New(t)
	a := loginAuth("user", "pass")

	proto, _, err := a.Start(&smtp.ServerInfo{Name: "mail.example.com", TLS: true})
	test.NoError(err)
	test.Equal("LOGIN", proto)
}

func TestLoginAuthAnswersChallenges(t *testing.T) {
	test := assert.New(t)
	a := loginAuth("user", "pass")
	if _, _, err := a.Start(&smtp.ServerInfo{Name: "mail.example.com", TLS: true}); err != nil {
		t.Fatalf("start: %v", err)
	}

	cases := []struct {
		prompt string
		want   string
	}{
		{"Username:", "user"},
		{"username", "user"},
		{"Username", "user"},
		{"Password:", "pass"},
		{"PASSWORD:", "pass"},
	}
	for _, c := range cases {
		got, err := a.Next([]byte(c.prompt), true)
		test.NoErrorf(err, "prompt %q", c.prompt)
		test.Equalf(c.want, string(got), "prompt %q", c.prompt)
	}
}
