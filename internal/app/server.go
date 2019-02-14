package app

import (
	// "github.com/jinzhu/gorm"
	// "github.com/jinzhu/gorm/dialects/postgres"
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
}

// Server contains gosupport server state
type Server struct {
	// path project directory root
	// NOTE: this would cause problems on some FaaS
	quitCh    chan struct{}
	root      string
	router    chi.Router
	domain    string
	port      string
	mailer    Mailer
	db        Database
	logger    kitlog.Logger
	templator *templator.Templator
}

// InitServer initialize new server instance with provided options & logger
func InitServer(
	l kitlog.Logger,
	t *templator.Templator,
	m Mailer,
	db Database,
	o Options,
) *Server {
	if o.Port[0] != ':' {
		o.Port = ":" + o.Port
	}
	s := Server{
		root:      o.Root,
		logger:    l,
		templator: t,
		mailer:    m,
		domain:    o.Domain,
		port:      o.Port,
	}
	s.InitRoutes(chi.NewRouter(), o.StaticFiles)
	return &s
}
