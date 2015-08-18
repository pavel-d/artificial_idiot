package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
)

func (t *Client) SendMessage(message *Message) *Message {
	params := map[string]string{
		"chat_id":             strconv.Itoa(message.From.Id),
		"text":                message.Text,
		"reply_to_message_id": strconv.Itoa(message.ReplyToMessageId),
	}

	if message.ReplyKeyboardMarkup != nil {
		keyboard, _ := json.Marshal(message.ReplyKeyboardMarkup)
		params["reply_markup"] = string(keyboard[:])
	}

	resp, _ := t.call("sendMessage", params)

	var replyMessage Message

	json.Unmarshal(resp.Result, &replyMessage)
	return &replyMessage
}

func (t *Client) SendPhoto(url string, caption string, chatId int) (*Message, error) {
	resp, err := http.Get(url)

	if err != nil {
		return &Message{}, err
	}

	defer resp.Body.Close()

	imageBuffer := &bytes.Buffer{}
	writer := multipart.NewWriter(imageBuffer)

	part, err := writer.CreateFormFile("photo", "photo.jpg")

	if err != nil {
		return &Message{}, err
	}

	_, err = io.Copy(part, resp.Body)

	params := map[string]string{
		"chat_id": strconv.Itoa(chatId),
		"caption": caption,
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	err = writer.Close()

	if err != nil {
		return &Message{}, err
	}

	req, _ := http.NewRequest("POST", fmt.Sprintf(BOT_API_URL, t.Token, "sendPhoto"), imageBuffer)
	resp, err = (&http.Client{}).Do(req)

	if err != nil {
		return &Message{}, err
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	log.Println(body[:])

	return &Message{}, nil
}
