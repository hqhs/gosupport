package main

import (
	"github.com/jinzhu/gorm"
)

// TelegramUser represents everyone, who uses `/start` command in any bot chat
type TelegramUser struct {
	gorm.Model
	Id string
}

type Helpdesker struct {
	gorm.Model
}

type Ticket struct {
	gorm.Model
}

type Message struct {
	gorm.Model
}

type Bot struct {
	gorm.Model
}
