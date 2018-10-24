package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

// Runbot starts bot in goroutine
func Runbot(token string) {
	bot, err := tgbotapi.NewBotAPI(token)
	Must(err)

	bot.Debug = true
	log.Printf("Started bot: %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	// Set timeout for requests to Telegram servers before bot stops polling (in seconds)
	u.Timeout = 600

	updates, err := bot.GetUpdatesChan(u)
	updates.Clear()

	go func() {
		for update := range updates {
			if update.Message == nil {
				continue
			}
			// log.Println(update)
		}
	}()

	return
}
