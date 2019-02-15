package app

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// BotType represents social network/messanger which bot uses
type BotType string

//
const (
	Telegram BotType = "TgBot"
	// NOTE slack support is not implemented yet
	Slack BotType = "SlackBot"
)

type Bot interface {
	Run(onExit func())
	Stop()
	Connector() *Connector
	HashToken(int) string
}

func (s *Server) Add(b Bot) {
	s.bots = append(s.bots, b)
	// Docker style management with unique hashes. Since bot name cannot be fetched
	// before request to bot api, first 8 characters from hashed token is used.
	var hash string
	for i := 0; ; i++ {
		hash = b.HashToken(i)
		if _, ok := s.conns[hash]; !ok {
			s.conns[hash] = b.Connector()
			s.logger.Log("msg", fmt.Sprintf("Added connector with hash %s", hash))
			return
		}
	}
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

type Connector struct {
	Input  <-chan Message // messages from Customers
	Output chan<- Message // messages to Customers
	Error  chan error   // Errors for admins
}

// BotOptions is type/token pair which is needed for bot to start working
type BotOptions struct {
	T              BotType
	Token          string
	EmailAuth      bool
	AllowedDomains []string
	Name           string
}

type defaultAnswers struct {
	Welcome         string `yaml:"welcome"`
	WaitingForEmail string `yaml:"waitingforemail"`
	SendingEmail    string `yaml:"sendingemail"`
	EmailSent       string `yaml:"emailsent"`
	WrongFormat     string `yaml:"wrongformat"`
	WrongToken      string `yaml:"wrongtoken"`
	Authorized      string `yaml:"authorized"`
}

// BaseBot is abstraction for useful business logic methods, such as authorization
// sequence etc. Every bot is basicly a REPL, with one difference:
// sometime your command result is given by helpdesker, and sometimes
// it's fully automatic. That single fact explains api design
// and provides ease of testing. Use it as embedded struct for various helpers
// and avoiding duplicating code. See MockBot implementation below
type BaseBot struct {
	s     *Server
	token string // don't log token, use hash below
	hash  string
	ans   defaultAnswers
	quit  chan struct{}
	// below is duplicating of Connector channels, but since real bots
	// inherit base bot and used as interfaces, whose channes needs to be wrapped
	// in separate struct for further usage
	Input  chan Message // messages from Customers
	Output chan Message // messages to Customers
	Error  chan error   // Errors for admins
}

func (b *BaseBot) Connector() *Connector {
	return &Connector{
		Input: b.Input,
		Output: b.Output,
		Error: b.Error,
	}
}

func (b *BaseBot) HashToken(try int) (hash string) {
	hasher := md5.New()
	hasher.Write([]byte(b.token))
	if try > 0 {
		bs := make([]byte, 4)
		binary.LittleEndian.PutUint32(bs, 31415926)
		hasher.Write(bs)
	}
	hash = hex.EncodeToString(hasher.Sum(nil))[0:8]
	b.hash = hash
	return
}

func (b *BaseBot) Stop() {
	close(b.quit)
}

// TgBot represents bot in telegram messanger
type TgBot struct {
	*BaseBot
	api *tgbotapi.BotAPI
	upd tgbotapi.UpdatesChannel
}

// NewTgBot ... FIXME
func NewTgBot(s *Server, token string) (t *TgBot, err error) {
	// TODO it's better to calc hashes here
	base := &BaseBot{
		s:      s,
		token:  token,
		quit:   make(chan struct{}),
		Input:  make(chan Message),
		Output: make(chan Message),
		Error:  make(chan error),
	}
	t = &TgBot{BaseBot: base}
	t.s.logger.Log("status", "Bot initialized")
	return
}

// Run starts polling on telegram bot
func (t *TgBot) Run(onExit func()) {
	t.s.logger.Log("status", fmt.Sprintf("Bot{%s} connecting to telegram api...", t.hash))
	var err error
	if t.api, err = tgbotapi.NewBotAPI(t.token); err != nil {
		return
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	t.upd, err = t.api.GetUpdatesChan(u)
	t.s.logger.Log("msg", fmt.Sprintf("Started polling on telegram bot{%s}", t.hash))
	for {
		select {
		case u := <-t.upd:
			go t.processUpdate(u)
		case <-t.quit:
			t.s.logger.Log("msg", "bot exited")
			onExit()
			return
		}
	}
}

func (t *TgBot) processUpdate(u tgbotapi.Update) {
	if u.Message == nil { // ignore any non-Message Updates
		return
	}
	t.s.logger.Log(
		"msg", "New message for tg bot",
		"from", u.Message.From.UserName,
		"text", u.Message.Text,
	)
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, u.Message.Text)
	msg.ReplyToMessageID = u.Message.MessageID
	t.api.Send(msg)
}

// MockBot takes reader/writer (ex. goes os.Stdin and os.Stdout) and
// simulates bot behavior and logic.
type MockBot struct {
	*BaseBot
}

// NewMockBot initializes bot for development/testing purposes
func NewMockBot(r io.Reader, w io.Writer) (*MockBot, error) {
	return &MockBot{}, nil
}

// Run ...
func (m *MockBot) Run() {

}
