package app

import (
	// "github.com/jinzhu/gorm"
	// "github.com/jinzhu/gorm/dialects/postgres"
	// "github.com/go-kit/kit/log/level"
	"github.com/go-chi/chi"
	kitlog "github.com/go-kit/kit/log"
)

// Options represents server initialization options
type Options struct {
	Domain        string
	Port          string
	EmailServer   string
	EmailAddress  string
	EmailPassword string
}

// Server contains gosupport server state
type Server struct {
	r *chi.Mux
}

// InitServer initialize new server instance with provided options & logger
func InitServer(logger kitlog.Logger, o Options) (*Server, error) {
	if logger == nil {
		logger = kitlog.NewNopLogger()
	}
	s := Server{}
	s.InitRoutes(chi.NewRouter())
	return &s, nil
}
