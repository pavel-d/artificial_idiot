package commands

import (
	"github.com/pavel-d/artificial_idiot/telegram"
)

func HelpHandler(message telegram.Message, bot *telegram.Bot, params []string) {
	helpMsg := `Hey, you can use the following commands:

              /image <search term> - Send an image
              /help - Display this message`

	message.Reply(helpMsg)
}
