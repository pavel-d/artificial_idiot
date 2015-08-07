package telegram_api

import (
	"encoding/json"
	"log"
)

type ApiResponse struct {
	Ok     bool         `json:"ok"`
	Result []ResultItem `json:"result"`
}

type ResultItem struct {
	Message  Message `json:"message"`
	UpdateId int     `json:"update_id"`
}

type Message struct {
	From      User                   `json:"from"`
	MessageId int                    `json:"message_id"`
	Chat      map[string]interface{} `json:"chat"`
	Date      int                    `json:"date"`
	Text      string                 `json:"text"`
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
	return apiResponse
}

func ParseMessage(jsonBlob []byte) *Message {
	message := &Message{}

	err := json.Unmarshal(jsonBlob, &message)

	if err != nil {
		log.Printf("Failed to parse api response: %v", err)
		log.Printf("Response body: %v", string(jsonBlob[:]))
	}
	return message
}
