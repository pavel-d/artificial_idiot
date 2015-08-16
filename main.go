package main

import (
	"github.com/pavel-d/artificial_idiot/commands"
	"github.com/pavel-d/artificial_idiot/telegram"
	"log"
	"os"
)

func main() {
	telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	if len(telegramToken) == 0 {
		log.Fatal("Please set TELEGRAM_BOT_TOKEN env variable")
	}

	client := &telegram.Client{Token: telegramToken}
	client.Connect()

	bot := telegram.MakeBot(client)

	bot.OnCommand("help", commands.HelpHandler)
	bot.OnCommand("image", commands.GoogleImageHandler)

	bot.Start()
}
