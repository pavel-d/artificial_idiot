package telegram

import (
	"log"
	"strings"
)

type Bot struct {
	Client          *Client
	CommandHandlers map[string]CommandHandler
	OneOffHandlers  map[int]OneOffHandler
}

type OneOffHandler func(message Message, bot *Bot)
type CommandHandler func(message Message, bot *Bot, params []string)

func (t *Bot) Start() {
	for msg := range t.Client.Consume() {

		if t.OneOffHandlers != nil {
			handler := t.OneOffHandlers[msg.From.Id]

			if handler != nil {
				go handler(msg, t)
				delete(t.OneOffHandlers, msg.From.Id)
			}
		}

		if t.CommandHandlers != nil && strings.HasPrefix(msg.Text, "/") {
			command := strings.Split(msg.Text[1:], " ")[0]
			handler := t.CommandHandlers[command]
			if handler != nil {
				go handler(msg, t, parseCommandArgs(command, &msg))
			} else {
				msg.Reply("Unrecognized command. What do you mean?")
			}
		}
	}
}

func parseCommandArgs(command string, message *Message) []string {
	if len(message.Text) < (len(command) + 3) {
		return []string{}
	}
	return strings.Split(message.Text, " ")[1:]
}

func (t *Bot) Once(userId int, handler OneOffHandler) {
	t.OneOffHandlers[userId] = handler
}

func (t *Bot) OnCommand(command string, handler CommandHandler) {
	log.Printf("Registered command handler for: /%s", command)
	t.CommandHandlers[command] = handler
}

func MakeBot(client *Client) *Bot {
	return &Bot{client, make(map[string]CommandHandler), make(map[int]OneOffHandler)}
}
