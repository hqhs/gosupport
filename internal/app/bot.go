package app

import "fmt"

// BotType represents social network/messanger which bot uses
type BotType string

//
const (
	Telegram BotType = "TgBot"
	Slack    BotType = "SlackBot"
)

// BotOptions is type/token pair which is needed for bot to start working
type BotOptions struct {
	T BotType
	Token string
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

// Connecter provides abstraction for varios available bot apis
type Connector interface {
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

// type Bot interface {

// }
