package app

import (
	"fmt"
	"io"

	// "github.com/go-telegram-bot-api/telegram-bot-api"
)

// BotType represents social network/messanger which bot uses
type BotType string

//
const (
	Telegram BotType = "TgBot"
	// NOTE slack support is not implemented yet
	Slack BotType = "SlackBot"
)

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

// InitBots initializes bots for their Options
func InitBots(o []BotOptions) error {
	if len(o) == 0 {
		// TODO init mock bot
		return nil
	} else {
		return fmt.Errorf("There's no support for real bots yet")
	}
}

// BaseBot is abstraction for useful business logic methods, such as authorization
// sequence etc. Every bot is basicly a REPL, with one difference:
// sometime your command result is given by helpdesker, and sometimes
// it's fully automatic. That single fact explains api design
// and provides ease of testing. Use it as embedded struct for various helpers
// and avoiding duplicating code. See MockBot implementation below
type BaseBot struct {
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

// TgBot represents bot in telegram messanger
type TgBot struct {
	s *Server
	*BaseBot
}

// NewTgBot initializes new telegram bot, if check=true bot makes
// single request to telegram API and return error if it wasn't successful
func NewTgBot(token string, check bool) error {
	// u := tgbotapi.NewUpdate(0)
	// u.Timeout = 60
	// updates, err := bot.api.GetUpdatesChan(u)
	return fmt.Errorf("Not implemented yet")
}

// Run starts polling on telegram bot
func (tg *TgBot) Run() {
	// FIXME get quit channel
	tg.s.logger.Log("msg", "Start polling on telegram bot.")

}
