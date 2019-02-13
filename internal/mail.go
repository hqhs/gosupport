package internal

import (
	"html/template"

	// "gopkg.in/gomail.v2"
	"github.com/go-kit/kit/log"
)

type authOptions struct {
	host     string
	port     string
	user     string
	password string
}

// Mailer interface provides api for sending mails to users
type Mailer interface {
	SendAuthMail() error
}

// mockMailer just logs letters to logger
type mockMailer struct {
	t            *template.Template
	l            log.Logger
	emailTimeout int // in seconds
}

func newMockMailer(t *template.Template, l log.Logger) (*mockMailer, error) {
	m := &mockMailer{}
	return m, nil
}

func (m *mockMailer) SendAuthMail(a authMail) error {
	return nil
}

type authMail struct {
	Receiver string
	Subject  string
	Body     string
	URL      string
	URLName  string
}
