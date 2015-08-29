package telegram

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	POLL_TIMEOUT    = 30
	BOT_API_URL     = "https://api.telegram.org/bot%s/%s"
	MSG_BUFFER_SIZE = 10
)

type Client struct {
	Token        string
	LastUpdateId int
	Active       bool
	messages     chan Message
}

func (t *Client) Connect() {
	log.Println("Connecting...")
	t.messages = make(chan Message, MSG_BUFFER_SIZE)
	t.LastUpdateId = 0
	t.Active = true
	go t.poll()
}

func (t *Client) Disconnect() {
	t.Active = false
}

func (t *Client) Consume() <-chan Message {
	// TODO: Allow multiple consumers
	return (<-chan Message)(t.messages)
}

func (t *Client) call(method string, params map[string]string) (*ApiResponse, error) {
	log.Printf("Calling Telegram %v method", method)

	requestUrl := fmt.Sprintf(BOT_API_URL, t.Token, method)

	values := url.Values{}

	for key, value := range params {
		values.Add(key, value)
	}

	log.Printf("Request url: %v", requestUrl)
	log.Printf("Request params: %v", values.Encode())

	var resp *http.Response
	var err error

	for {
		resp, err = http.PostForm(requestUrl, values)

		if err == nil {
			break
		}

		log.Println(err)
		time.Sleep(1 * time.Second)
		log.Println("Retrying...")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return ParseApiResponse(body), nil
}

func (t *Client) getUpdates() ([]Update, error) {
	resp, err := t.call("getUpdates", map[string]string{
		"timeout": strconv.Itoa(POLL_TIMEOUT),
		"offset":  strconv.Itoa(t.LastUpdateId + 1),
	})

	if err != nil {
		return nil, err
	}
	var updates []Update
	json.Unmarshal(resp.Result, &updates)
	return updates, nil
}

func (t *Client) poll() {
	log.Println("Updates polling started")

	for t.Active {
		resp, err := t.getUpdates()

		if err != nil {
			log.Panicln(err)
		} else {
			log.Printf("%v updates received", len(resp))
			t.processUpdates(resp)
		}

	}
}

func (t *Client) processUpdates(updates []Update) {
	for _, update := range updates {
		if update.UpdateId > t.LastUpdateId {
			t.LastUpdateId = update.UpdateId
		}

		if &update.Message != nil {
			update.Message.client = t
			t.messages <- update.Message
		}

		log.Printf("Update Id: %v", update.UpdateId)
	}
}
