package main

import (
	"github.com/pavel-d/artificial_idiot/commands"
	"github.com/pavel-d/artificial_idiot/telegram"
	"github.com/pavel-d/artificial_idiot/tools"
	"log"
	"os"
	"strconv"
)

func main() {
	telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	healthCheckPort, _ := strconv.Atoi(os.Getenv("PORT"))

	if healthCheckPort == 0 {
		healthCheckPort = 3000
	}

	if len(telegramToken) == 0 {
		log.Fatal("Please set TELEGRAM_BOT_TOKEN env variable")
	}

	client := &telegram.Client{Token: telegramToken}
	client.Connect()

	bot := telegram.MakeBot(client)

	bot.OnCommand("help", commands.HelpHandler)

	bot.OnCommand("img", commands.ImageFinder)
	bot.OnCommand("gif", commands.GifFinder)
	bot.OnCommand("g", commands.GoogleSearch)
	// bot.OnCommand("video", commands.VideoFinder)

	bot.OnCommand("currency", commands.CurrencyHandler)

	tools.StartHealthCheck(healthCheckPort)
	bot.Start()
}
