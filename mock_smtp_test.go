package tinymail

import (
	"bufio"
	"net"
	"strings"
	"sync"
	"testing"
)

// mockSMTP is a minimal in-process SMTP server used to capture the SMTP
// conversation (envelope recipients and the raw DATA payload) so the send
// paths can be tested without a real mail server.
type mockSMTP struct {
	ln net.Listener

	mu       sync.Mutex
	mailFrom string
	rcpts    []string
	data     string
}

// newMockSMTP starts a mock server on a random loopback port and stops it
// when the test finishes. The loopback host (127.0.0.1) lets net/smtp's
// PlainAuth send credentials over the unencrypted test connection.
func newMockSMTP(t *testing.T) *mockSMTP {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("mock smtp listen: %v", err)
	}
	m := &mockSMTP{ln: ln}
	go m.serve()
	t.Cleanup(func() { ln.Close() })
	return m
}

func (m *mockSMTP) port() int {
	return m.ln.Addr().(*net.TCPAddr).Port
}

func (m *mockSMTP) serve() {
	conn, err := m.ln.Accept()
	if err != nil {
		return
	}
	defer conn.Close()

	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)
	reply := func(s string) {
		w.WriteString(s + "\r\n")
		w.Flush()
	}

	reply("220 mock ready")
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return
		}
		line = strings.TrimRight(line, "\r\n")
		cmd := strings.ToUpper(line)

		switch {
		case strings.HasPrefix(cmd, "EHLO"), strings.HasPrefix(cmd, "HELO"):
			reply("250-mock greets you")
			reply("250 AUTH PLAIN LOGIN")
		case strings.HasPrefix(cmd, "AUTH"):
			reply("235 2.7.0 authentication succeeded")
		case strings.HasPrefix(cmd, "MAIL FROM"):
			m.mu.Lock()
			m.mailFrom = line
			m.mu.Unlock()
			reply("250 2.1.0 ok")
		case strings.HasPrefix(cmd, "RCPT TO"):
			m.mu.Lock()
			m.rcpts = append(m.rcpts, addr(line))
			m.mu.Unlock()
			reply("250 2.1.5 ok")
		case strings.HasPrefix(cmd, "DATA"):
			reply("354 end data with <CR><LF>.<CR><LF>")
			var sb strings.Builder
			for {
				dl, err := r.ReadString('\n')
				if err != nil {
					break
				}
				if dl == ".\r\n" || dl == ".\n" {
					break
				}
				sb.WriteString(dl)
			}
			m.mu.Lock()
			m.data = sb.String()
			m.mu.Unlock()
			reply("250 2.0.0 ok queued")
		case strings.HasPrefix(cmd, "QUIT"):
			reply("221 2.0.0 bye")
			return
		default:
			reply("250 ok")
		}
	}
}

// recipients returns the addresses that received an envelope RCPT TO.
func (m *mockSMTP) recipients() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]string(nil), m.rcpts...)
}

// body returns the raw DATA payload as received by the server.
func (m *mockSMTP) body() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.data
}

// addr extracts the address out of a "RCPT TO:<addr>" command.
func addr(line string) string {
	if i := strings.Index(line, "<"); i >= 0 {
		if j := strings.Index(line, ">"); j > i {
			return line[i+1 : j]
		}
	}
	return line
}
