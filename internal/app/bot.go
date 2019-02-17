package app

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"io"
	"time"

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
	GetFileDirectURL(string) (string, error)
}

type Connector struct {
	Input  chan []byte // messages from Customers
	Output chan<- Message // messages to Customers
	Error  chan error     // Errors for admins
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
	Input  chan []byte // messages from Customers
	Output chan Message // messages to Customers
	Error  chan error   // Errors for admins
}

func (b *BaseBot) Connector() *Connector {
	return &Connector{
		Input:  b.Input,
		Output: b.Output,
		Error:  b.Error,
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
func NewTgBot(ctx context.Context, s *Server, token string) (t *TgBot, err error) {
	// TODO it's better to calc hashes here
	base := &BaseBot{
		s:      s,
		token:  token,
		quit:   make(chan struct{}),
		Input:  make(chan []byte),
		Output: make(chan Message),
		Error:  make(chan error),
	}
	t = &TgBot{BaseBot: base}

	ready := make(chan struct{})
	go func() {
		t.s.logger.Log("status", "connecting to telegram api...")
		var err error
		if t.api, err = tgbotapi.NewBotAPI(t.token); err != nil {
			return
		}
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 5
		t.upd, err = t.api.GetUpdatesChan(u)
		ready <- struct{}{}
	}()
	select {
	case <-ready:
		t.s.logger.Log("status", "Bot initialized")
		return
	case <-ctx.Done():
		t.s.logger.Log("msg", "exiting bot")
	}
	return
}

// Run starts polling on telegram bot
func (t *TgBot) Run(onExit func()) {
	// Telegram servers is blocked in some countries and connecting may take
	// a lot of time. But graceful reload should stop setup, therefore separate
	// setup method.
	rate := time.Second / 30
	throttle := time.Tick(rate)
	for {
		select {
		case m := <-t.Output:
			<-throttle
			t.s.logger.Log("sending", m.Text)
			go t.sendMessage(m)
		case u := <-t.upd:
			go t.processUpdate(u)
		case <-t.quit:
			t.s.logger.Log("msg", "bot exited")
			onExit()
			return
		}
	}
}

func (t *TgBot) sendMessage(m Message) {
	// FIXME this little hash works, but for real world
	// we need to fetch user from message.UserID and then provide user.ChatID
	toSend := tgbotapi.NewMessage(int64(m.UserID), m.Text)
	_, err := t.api.Send(toSend)
	if err != nil {
		t.s.logger.Log("err", err, "then", "during sending message from admin to user")
	}
	// FIXME save messages from dashboard
	// d.SaveMessage(int(message.ChatID), m)
	// broadcast for other admins
	// resend := t.parseMessage(&response)
	// t.Input <- *resend
}

func (t *TgBot) GetFileDirectURL(fileID string) (string, error) {
	return t.api.GetFileDirectURL(fileID)
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
	ctx := context.Background()
	user := &User{}
	msg := t.parseMessage(u.Message)
	var reply string

	conn, _ := t.s.DB.Conn(ctx)
	defer conn.Close()
	userID := u.Message.From.ID
	// try to get user from database
	query := "SELECT user_id FROM users WHERE user_id=$1"
	if err := conn.QueryRowContext(ctx, query, userID).Scan(&user.UserID); err != nil {
		switch err {
		case sql.ErrNoRows:
			t.parseUserInfo(user, u)
			if err = t.saveUser(ctx, conn, user); err != nil {
				reply = "Internal server error. Try again later"
				t.s.logger.Log("err", err.Error())
				goto response
			}
		default:
			t.s.logger.Log("err", err.Error())
			reply = "Internal server error. Try again later"
			goto response
		}
	}
	// user updated after message is saved
	if err := t.saveMessage(ctx, conn, msg); err != nil {
		t.s.logger.Log("err", err.Error())
		reply = "Internal server error. Try again later"
		goto response
	}
	reply = "Resending your message to admins..."
response:
	output := tgbotapi.NewMessage(u.Message.Chat.ID, reply)
	output.ReplyToMessageID = u.Message.MessageID
	t.api.Send(output)
	payload, _ := json.Marshal(*msg)
	t.Input <- payload
}

func (t *TgBot) parseUserInfo(u *User, update tgbotapi.Update) {
	u.UserID = update.Message.From.ID
	u.ChatID = update.Message.Chat.ID
	u.Email = ""
	u.Name = update.Message.From.FirstName + " " + update.Message.From.LastName
	u.Username = update.Message.From.UserName
	u.IsAuthorized = false
	u.AuthToken = []byte{}
	u.IsTokenExpired = false
	u.UserPhotoID = ""
	u.UpdatedAt = time.Now()
	u.CreatedAt = time.Now()
}

func (t *TgBot) saveUser(ctx context.Context, conn *sql.Conn, u *User) (err error) {
	query := `INSERT INTO users(user_id, created_at, updated_at,
			chat_id, email, name, username, has_unread_messages)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	_, err = conn.ExecContext(ctx, query, u.UserID, u.CreatedAt, u.UpdatedAt,
		u.ChatID, u.Email, u.Name, u.Username, u.HasUnreadMessages)
	return
}

func (t *TgBot) saveMessage(ctx context.Context, conn *sql.Conn, m *Message) (err error) {
	query := `INSERT INTO messages(user_id, message_id, is_broadcast,
			from_admin, created_at, updated_at, text, reply_to_message,
			document_id, photo_id)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err = conn.ExecContext(ctx, query, m.UserID, m.MessageID, false,
		false, time.Now(), time.Now(), m.Text, m.ReplyToMessage,
		m.DocumentID, m.PhotoID)
	if err != nil {
		return err
	}
	query = `UPDATE users SET last_message_id=($1) WHERE user_id=($2)`
	_, err = conn.ExecContext(ctx, query, m.MessageID, m.UserID)
	return
}

func (t *TgBot) parseMessage(message *tgbotapi.Message) *Message {
	isBot := message.From.IsBot
	m := Message{
		FromAdmin:  isBot,
		MessageID:  message.MessageID,
		UserID:     message.From.ID,
		ChatID:     message.Chat.ID,
		Text:       message.Text,
		Date:       message.Date,
		DocumentID: "",
		PhotoID:    "",
	}
	if message.ForwardFrom != nil {
		m.ForwardFrom = message.ForwardFrom.ID
		m.ForwardDate = message.ForwardDate
	}
	if message.ReplyToMessage != nil {
		m.ReplyToMessage = message.ReplyToMessage.MessageID
	}
	if message.Document != nil {
		m.DocumentID = message.Document.FileID
	}
	if message.Photo != nil {
		if p := *message.Photo; len(p) > 0 {
			m.PhotoID = p[len(p)-1].FileID
		}
	}
	return &m
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
