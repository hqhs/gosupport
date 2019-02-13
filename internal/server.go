package internal

import (
	// "github.com/jinzhu/gorm"
	// "github.com/jinzhu/gorm/dialects/postgres"
	// "github.com/go-kit/kit/log/level"
	"github.com/go-chi/chi"
	kitlog "github.com/go-kit/kit/log"
)

// Options represents server initialization options
type Options struct {
	Root          string
	Domain        string
	Port          string
	EmailServer   string
	EmailAddress  string
	EmailPassword string
}

// Server contains gosupport server state
type Server struct {
	// path project directory root
	// NOTE: this would cause problems on some FaaS
	root      string
	router    chi.Router
	domain    string
	port      string
	mailer    *Mailer
	logger    kitlog.Logger
	templator *Templator
}

// InitServer initialize new server instance with provided options & logger
func InitServer(
	logger kitlog.Logger,
	templator *Templator,
	o Options,
) *Server {
	if o.Port[0] != ':' {
		o.Port = ":" + o.Port
	}
	s := Server{
		root: o.Root,
		logger: logger,
		templator: templator,
		domain: o.Domain,
		port:   o.Port,
	}
	s.InitRoutes(chi.NewRouter())
	return &s
}
