package telegram_api

type TelegramBot struct {
	Client          *TelegramClient
	MessageHandlers map[string]MessageHandler
}

type MessageHandler func(message Message, client TelegramClient)

func (t *TelegramBot) Start() {
	// go func() {
	for msg := range t.Client.Consume() {

		if t.MessageHandlers != nil {
			for _, handler := range t.MessageHandlers {
				handler(msg, *t.Client)
			}
		}
	}
	// }()
}

func (t *TelegramBot) On(command string, handler MessageHandler) {
	t.MessageHandlers[command] = handler
}

func MakeBot(client *TelegramClient) *TelegramBot {
	return &TelegramBot{client, make(map[string]MessageHandler)}
}
