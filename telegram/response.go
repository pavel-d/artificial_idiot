package telegram

import (
	"encoding/json"
	"log"
)

type ApiResponse struct {
	Ok     bool            `json:"ok"`
	Result json.RawMessage `json:"result"`
}

type Update struct {
	Message  Message `json:"message"`
	UpdateId int     `json:"update_id"`
}

type Message struct {
	From                User                   `json:"from"`
	MessageId           int                    `json:"message_id"`
	Chat                map[string]interface{} `json:"chat"`
	Date                int                    `json:"date"`
	Text                string                 `json:"text"`
	ReplyToMessageId    int                    `json:"-"`
	ReplyKeyboardMarkup *ReplyKeyboardMarkup   `json:"-"`
	client              *Client                `json:"-"`
}

func (m *Message) Reply(text string) *Message {
	return m.ReplyWithKeyboardMarkup(text, nil)
}

func (m *Message) ReplyWithKeyboardMarkup(text string, keyboard *ReplyKeyboardMarkup) *Message {
	message := &Message{
		From:                m.From,
		Text:                text,
		ReplyKeyboardMarkup: keyboard,
	}
	return m.client.SendMessage(message)
}

type ReplyKeyboardMarkup struct {
	Keyboard        [][]string `json:"keyboard"`
	ResizeKeyboard  bool       `json:"resize_keyboard"`
	OneTimeKeyboard bool       `json:"one_time_keyboard"`
	Selective       bool       `json:"selective"`
}

func KeyboardForOptions(options ...string) *ReplyKeyboardMarkup {
	replyOptions := make([]string, len(options))

	for idx, option := range options {
		replyOptions[idx] = option
	}

	return &ReplyKeyboardMarkup{
		Keyboard: [][]string{
			replyOptions,
		},
		OneTimeKeyboard: true,
	}
}

type User struct {
	Id        int    `json:"id"`
	LastName  string `json:"last_name"`
	FirstName string `json:"first_name"`
}

func ParseApiResponse(jsonBlob []byte) *ApiResponse {
	apiResponse := &ApiResponse{}

	err := json.Unmarshal(jsonBlob, &apiResponse)

	if err != nil {
		log.Printf("Failed to parse api response: %v", err)
		log.Printf("Response body: %v", string(jsonBlob[:]))
	}

	if !apiResponse.Ok {
		log.Printf("Api returned error: %v", apiResponse.Result)
	}

	return apiResponse
}
