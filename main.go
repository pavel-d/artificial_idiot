package main

import (
	"github.com/pavel-d/artificial_idiot/telegram_api"
	"log"
	"os"
)

func main() {
	telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	if len(telegramToken) == 0 {
		log.Fatal("Please set TELEGRAM_BOT_TOKEN env variable")
	}

	client := &telegram_api.TelegramClient{Token: telegramToken}
	client.Connect()

	bot := telegram_api.MakeBot(client)

	bot.On("help", func(message telegram_api.Message, client telegram_api.TelegramClient) {
		client.SendMessage(message.From.Id, message.Text)
	})

	bot.Start()
}
