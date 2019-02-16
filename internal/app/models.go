package app

import (
	"fmt"
	"gopkg.in/jinzhu/gorm.v1"
)

// Admin represents registered dashboard Helpdesker
type Admin struct {
	gorm.Model
	Email          string `gorm:"UNIQUE,INDEX" json:"email"`
	Name           string `json:"name"`
	HashedPassword string `json:"hashed_password"`
	IsSuperUser    bool   `json:"is_superuser"`
	IsActive       bool   `json:"is_active"`
	EmailConfirmed bool   `json:"email_confirmed"`
	// gorm automaticly tracks those 2
	AuthToken          string    `json:"auth_token"`
	PasswordResetToken string    `json:"password_reset_token"`
}

// User represents single chat endpoint of communication
type User struct {
	gorm.Model
	UserID            int    `json:"userid"`
	ChatID            int64  `json:"chatid"`
	Email             string `json:"email"`
	Name              string `json:"name"`
	Username          string `json:"username"`
	HasUnreadMessages bool   `json:"has_unread_messages"`
	// AuthToken used for email authorization
	AuthToken      []byte  `json:"authtoken"`
	IsAuthorized   bool    `json:"isauthorized"`
	IsTokenExpired bool    `json:"is_token_expired"`
	LastMessageAt  int64   `json:"lastMessageAt"`
	// FIXME
	LastMessage    Message `json:"last_message"`
	UserPhotoID    string  `json:"user_photo_id"`
	IsActive       bool    `json:"is_active"`
}

func (u User) String() string {
	if len(u.Username) > 0 {
		return fmt.Sprintf("%s %s (@%s)", u.Name, u.Email, u.Username)
	}
	return fmt.Sprintf("%s %s", u.Name, u.Email)
}

// Message is atomic piece of communication between User and Admin
type Message struct {
	gorm.Model
	FromAdmin   bool `json:"from_bot"`
	IsBroadcast bool `json:"is_broadcast"`
	// I'm not sure what message id's are unique
	MessageID      int    `json:"message_id"`
	UserID         int    `json:"user_id"`
	ChatID         int64  `json:"chat_id"`
	Text           string `json:"text"`
	Date           int    `json:"send_date"`
	ForwardFrom    int    `json:"forward_from"`
	ForwardDate    int    `json:"forward_date"`
	ReplyToMessage int    `json:"reply_to_message_id"`
	EditDate       int    `json:"edit_date"`
	// those are just file ids for FileProxy, currently only Photo and Document
	// are supported
	DocumentID string `json:"file_id"`
	PhotoID    string `json:"photo_id"`
}
