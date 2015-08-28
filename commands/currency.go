package commands

import (
	"fmt"
	"github.com/pavel-d/artificial_idiot/telegram"
	"golang.org/x/net/websocket"
	"log"
	"strings"
)

func CurrencyHandler(message telegram.Message, bot *telegram.Bot, params []string) {
	keyboard := telegram.KeyboardForOptions("Книжка", "НБУ", "RUB")
	message.ReplyWithKeyboardMarkup("Какой курс?", keyboard)

	bot.Once(message.From.Id, func(message telegram.Message, bot *telegram.Bot) {
		if message.Text == "RUB" {
			message.Reply(fetchRub())
		}
	})
}

func fetchRub() string {
	origin := "http://zenrus.ru/"
	url := "ws://zenrus.ru:8888/"

	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	defer ws.Close()

	var msg = make([]byte, 512)

	if n, err = ws.Read(msg); err != nil {
		log.Fatal(err)
	}

	rates := strings.Split(string(msg[:n]), ";")

	return fmt.Sprintf("USD/RUB - %s; EUR/RUB - %s", rates[0], rates[1])
}
