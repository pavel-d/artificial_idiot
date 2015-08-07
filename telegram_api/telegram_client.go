package telegram_api

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

const (
	POLL_TIMEOUT    = 30
	BOT_API_URL     = "https://api.telegram.org/bot"
	MSG_BUFFER_SIZE = 10
)

type TelegramClient struct {
	Token        string
	LastUpdateId int
	Active       bool
	messages     chan Message
}

func (t *TelegramClient) Connect() {
	log.Println("Connecting...")
	t.messages = make(chan Message, MSG_BUFFER_SIZE)
	t.LastUpdateId = 0
	t.Active = true
	go t.poll()
}

func (t *TelegramClient) Disconnect() {
	t.Active = false
}

func (t *TelegramClient) Consume() <-chan Message {
	// TODO: Allow multiple consumers
	return (<-chan Message)(t.messages)
}

func (t *TelegramClient) SendMessage(chatId int, text string) {
	params := map[string]string{
		"chat_id": strconv.Itoa(chatId),
		"text":    text,
	}

	t.call("sendMessage", params)
}

func (t *TelegramClient) call(method string, params map[string]string) (*ApiResponse, error) {
	log.Printf("Calling Telegram %v method", method)

	requestUrl, err := url.Parse(t.getApiUrl() + method)
	q := requestUrl.Query()

	for key, value := range params {
		q.Set(key, value)
	}

	requestUrl.RawQuery = q.Encode()

	log.Printf("Request url: %v", requestUrl)
	resp, err := http.Get(requestUrl.String())

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return ParseApiResponse(body), nil
}

func (t *TelegramClient) getApiUrl() string {
	return BOT_API_URL + t.Token + "/"
}

func (t *TelegramClient) poll() {
	log.Println("Updates polling started")

	for t.Active {
		resp, err := t.call("getUpdates", map[string]string{"timeout": "5", "offset": strconv.Itoa(t.LastUpdateId + 1)})

		if err != nil {
			log.Panicln(err)
		} else {
			log.Printf("%v updates received", len(resp.Result))
			if resp.Ok {
				t.processUpdates(resp)
			} else {
				log.Println("Response status is not Ok")
			}
		}

	}
}

func (t *TelegramClient) processUpdates(apiResponse *ApiResponse) {
	for _, update := range apiResponse.Result {
		if update.UpdateId > t.LastUpdateId {
			t.LastUpdateId = update.UpdateId
		}

		if &update.Message != nil {
			t.messages <- update.Message
		}

		log.Printf("Update Id: %v", update.UpdateId)
	}
}
