package app

import (
	"fmt"
	"sync"
	"database/sql"

	// "github.com/go-kit/kit/log/level"
	"github.com/go-chi/chi"
	kitlog "github.com/go-kit/kit/log"
	"github.com/hqhs/gosupport/pkg/templator"
	_ "github.com/lib/pq" // postgres driver
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
	TgBotTokens   string
}

// Server contains gosupport server state
type Server struct {
	// path project directory root
	QuitCh chan struct{}
	Secret string
	DB     *sql.DB
	// Unexported fields
	root   string // NOTE: this would cause problems on some FaaS
	router chi.Router
	domain string
	port   string
	mailer Mailer
	// NOTE: use sqlmock for testing
	logger    kitlog.Logger
	templator *templator.Templator
	// Since bot is a REPL (read bot.go) we need to store all bot
	// interfaces and make sure each return chan with messages.
	bots  map[string]Bot
	conns map[string]*Connector
	hubs  map[string]*Hub
	// hubs []Hub
	botGroup *sync.WaitGroup
}

// InitServer initialize new server instance with provided options & logger
func InitServer(
	l kitlog.Logger,
	t *templator.Templator,
	m Mailer,
	db *sql.DB,
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
		bots:      make(map[string]Bot, 0),
		conns:     make(map[string]*Connector),
		hubs:      make(map[string]*Hub),
		botGroup:  &sync.WaitGroup{},
	}
	s.InitRoutes(chi.NewRouter(), o.StaticFiles)
	return &s
}

func (s *Server) Add(b Bot) {
	// Docker style management with unique hashes. Since bot name cannot be fetched
	// before request to bot api, first 8 characters from hashed token is used.
	var hash string
	for i := 0; ; i++ {
		hash = b.HashToken(i)
		if _, ok := s.conns[hash]; !ok {
			s.conns[hash] = b.Connector()
			s.logger.Log("msg", fmt.Sprintf("Added connector with hash %s", hash))
			break
		}
	}
	s.bots[hash] = b
	// TODO fix broadcasting
	broadcast := make(chan []byte)
	s.hubs[hash] = NewHub(broadcast)
}

func (s *Server) RunBots() {
	for _, b := range s.bots {
		onExit := func() { s.botGroup.Add(-1) }
		go b.Run(onExit)
		s.botGroup.Add(1)
	}
}

func (s *Server) StopBots() {
	s.logger.Log("status", "Performing graceful shutdown of bots...")
	for _, b := range s.bots {
		b.Stop()
	}
}
