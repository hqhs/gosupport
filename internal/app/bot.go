package app

import (
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
}


func (s *Server) Add(b Bot) {
	s.bots = append(s.bots, b)
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
	ans defaultAnswers
}

// TgBot represents bot in telegram messanger
type TgBot struct {
	*BaseBot
	s   *Server
	api *tgbotapi.BotAPI
	upd tgbotapi.UpdatesChannel
	quit chan struct{}
}

// NewTgBot initializes new telegram bot, bot makes single request to telegram
// API as part of initialization process and return error if it wasn't successful.
// DB, logger, domain, port, and mailer (this uses templator internally) is used
// from server
func NewTgBot(s *Server, token string) (t *TgBot, err error) {
	t = &TgBot{s: s, quit: make(chan struct{})}
	t.s.logger.Log("status", "Connecting to telegram api...")
	if t.api, err = tgbotapi.NewBotAPI(token); err != nil {
		return
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	t.upd, err = t.api.GetUpdatesChan(u)
	t.s.logger.Log("status", "Bot initialized")
	return
}

// Run starts polling on telegram bot
func (t *TgBot) Run(onExit func()) {
	// FIXME get quit channel
	t.s.logger.Log("msg", "Start polling on telegram bot.")
	for {
		select {
		case u := <- t.upd:
			go t.processUpdate(u)
		case <-t.quit:
			t.s.logger.Log("msg", "bot exited")
			onExit()
			return
		}
	}
}

func (t *TgBot) Stop() {
	close(t.quit)
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
