package app

import (
	"gopkg.in/jinzhu/gorm.v1"
	// "github.com/go-kit/kit/log/level"
	"github.com/go-chi/chi"
	kitlog "github.com/go-kit/kit/log"
	"github.com/hqhs/gosupport/pkg/templator"
)

// Options represents server initialization options
type Options struct {
	Root          string
	Domain        string
	Port          string
	EmailServer   string
	EmailAddress  string
	EmailPassword string
	DatabaseURL   string
	ServeStatic   bool
	StaticFiles   string
	DbOptions     DbOptions
	Secret        string
}

// Server contains gosupport server state
type Server struct {
	// path project directory root
	// NOTE: this would cause problems on some FaaS
	QuitCh chan struct{}
	Secret string
	DB     *gorm.DB
	// Unexported fields
	root   string
	router chi.Router
	domain string
	port   string
	mailer Mailer
	// NOTE: use sqlmock for testing
	logger    kitlog.Logger
	templator *templator.Templator
}

// InitServer initialize new server instance with provided options & logger
func InitServer(
	l kitlog.Logger,
	t *templator.Templator,
	m Mailer,
	db *gorm.DB,
	o Options,
) *Server {
	if o.Port[0] != ':' {
		o.Port = ":" + o.Port
	}
	s := Server{
		QuitCh:    make(chan struct{}),
		root:      o.Root,
		logger:    l,
		templator: t,
		mailer:    m,
		domain:    o.Domain,
		port:      o.Port,
		DB:        db,
		Secret:    o.Secret,
	}
	s.InitRoutes(chi.NewRouter(), o.StaticFiles)
	return &s
}
