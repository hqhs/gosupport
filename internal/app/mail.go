package app

import (
	// "gopkg.in/gomail.v2"
	"github.com/go-kit/kit/log"
	"github.com/hqhs/gosupport/pkg/templator"
)

type authOptions struct {
	host     string
	port     string
	user     string
	password string
}

// Mailer interface provides api for sending mails to users
type Mailer interface {
	SendAuthMail(AuthMail) error
}

// mockMailer just logs letters to logger
type mockMailer struct {
	templator    *templator.Templator
	logger       log.Logger
	emailTimeout int // in seconds
}

// NewMockMailer initializes mock email client. It logs new messages to provided
// logger.
// NOTE: I'm not sure tampletor need to be in args.
func NewMockMailer(t *templator.Templator, l log.Logger) Mailer {
	return &mockMailer{t, l, 30}
}

func (m *mockMailer) SendAuthMail(a AuthMail) error {
	m.logger.Log("msg", "Email sent",
		"Receiver", a.Receiver,
		"Subject", a.Subject,
		"Body", a.Body,
		"URL", a.URL,
	)
	return nil
}

// AuthMail represents data needed for user/admin authorization email.
type AuthMail struct {
	Receiver string
	Subject  string
	Body     string
	URL      string
	URLName  string
}
