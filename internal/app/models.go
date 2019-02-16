package app

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type model struct {
	ID        uint      `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func newModel() model {
	return model{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Admin represents registered dashboard Helpdesker
type Admin struct {
	model
	Email              string `json:"email" db:"email"`
	Name               string `json:"name" db:"name"`
	HashedPassword     string `json:"hashed_password" db:"hashed_password"`
	IsSuperUser        bool   `json:"is_superuser" db:"is_superuser"`
	IsActive           bool   `json:"is_active" db:"is_active"`
	EmailConfirmed     bool   `json:"email_confirmed" db:"email_confirmed"`
	AuthToken          string `json:"auth_token" db:"auth_token"`
	PasswordResetToken string `json:"password_reset_token" db:"password_reset_token"`
}

func NewAdmin(email string, password string, isSuperUser bool) (a *Admin) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), defaultPasswordHashingCost)
	a = &Admin{
		model:          newModel(),
		Email:          email,
		HashedPassword: string(hash),
		IsSuperUser:    isSuperUser,
		IsActive:       true,
		EmailConfirmed: false,
	}
	return
}

// User represents single chat endpoint of communication
type User struct {
	UserID            int    `json:"userid"`
	ChatID            int64  `json:"chatid"`
	Email             string `json:"email"`
	Name              string `json:"name"`
	Username          string `json:"username"`
	HasUnreadMessages bool   `json:"has_unread_messages"`
	// AuthToken used for email authorization
	AuthToken      []byte `json:"authtoken"`
	IsAuthorized   bool   `json:"isauthorized"`
	IsTokenExpired bool   `json:"is_token_expired"`
	LastMessageAt  int64  `json:"lastMessageAt"`
	// FIXME
	LastMessage Message `json:"last_message"`
	UserPhotoID string  `json:"user_photo_id"`
	IsActive    bool    `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (u User) String() string {
	if len(u.Username) > 0 {
		return fmt.Sprintf("%s %s (@%s)", u.Name, u.Email, u.Username)
	}
	return fmt.Sprintf("%s %s", u.Name, u.Email)
}

// Message is atomic piece of communication between User and Admin
type Message struct {
	model
	FromAdmin   bool `json:"from_bot"`
	IsBroadcast bool `json:"is_broadcast"`
	// I'm not sure what message id's are unique
	Text           string    `json:"text"`
	MessageID      int       `json:"message_id"`
	UserID         int       `json:"user_id"`
	ChatID         int64     `json:"chat_id"`
	Date           int       `json:"send_date"`
	ForwardFrom    int       `json:"forward_from"`
	ForwardDate    int       `json:"forward_date"`
	ReplyToMessage int       `json:"reply_to_message_id"`
	EditDate       int       `json:"edit_date"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	// those are just file ids for FileProxy, currently only Photo and Document
	// are supported
	DocumentID string `json:"file_id"`
	PhotoID    string `json:"photo_id"`
}
