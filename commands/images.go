package commands

import (
	"github.com/pavel-d/artificial_idiot/telegram"
	"github.com/pavel-d/artificial_idiot/tools/google"
	"strings"
)

func GoogleImageHandler(message telegram.Message, bot *telegram.Bot, params []string) {
	searchPhrase := strings.Join(params, " ")

	image, err := google.RandomImage(searchPhrase)

	if err != nil {
		message.Reply(err.Error())
		return
	}
	message.Reply(image)
	keyboard := telegram.KeyboardForOptions("Yes", "No")
	message.ReplyWithKeyboardMarkup("Want more?", keyboard)

	bot.Once(message.From.Id, func(message telegram.Message, bot *telegram.Bot) {
		if message.Text == "Yes" {
			GoogleImageHandler(message, bot, params)
		}
	})
}
