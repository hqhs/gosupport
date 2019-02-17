package app

import (
	"database/sql"
	"net/http"
	"context"
	"fmt"
	"sync"

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
	QuitCh chan struct{}
	Ctx    context.Context
	Cancel func()

	Secret string
	DB     *sql.DB
	// Unexported fields
	// project directory root
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
	ctx, cancel := context.WithCancel(context.Background())
	s := Server{
		QuitCh:    make(chan struct{}),
		Ctx:       ctx,
		Cancel:    cancel,
		root:      o.Root,
		logger:    l,
		templator: t,
		mailer:    m,
		domain:    o.Domain,
		port:      o.Port,
		DB:        db,
		Secret:    o.Secret,
		// this three separate entities should be one struct with Bot Interface, Hub and connector (?), though
		// I'm not sure last one is necessary
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
			s.bots[hash] = b
			c := b.Connector()
			s.conns[hash] = c
			s.hubs[hash] = NewHub(c.Input)
			s.logger.Log("msg", fmt.Sprintf("Added connector with hash %s", hash))
			return
		}
	}
}

// ServeRouter starts http server
func (s *Server) ListenAndServe() {
	s.RunBots()
	s.logger.Log("status", "Starting serving routes")
	httpServer := &http.Server{Addr: s.port, Handler: s.router}
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			s.logger.Log("err", err, "msg", "Gracefully shutting down the server")
			return
		}
	}()
	<-s.Ctx.Done()
	s.logger.Log("status", "Waiting then all bots are done...")
	httpServer.Shutdown(context.Background())
	s.DB.Close()
}

func (s *Server) RunBots() {
	for k, b := range s.bots {
		onExit := func() { s.botGroup.Add(-1) }
		go b.Run(onExit)
		go s.hubs[k].run()
		s.botGroup.Add(1)
	}
}

func (s *Server) Shutdown() {
	s.StopBots()
	s.botGroup.Wait()
	s.Cancel()
}

func (s *Server) StopBots() {
	s.logger.Log("status", "Performing graceful shutdown of bots...")
	for _, b := range s.bots {
		// FIXME close hubs
		b.Stop()
	}
}
