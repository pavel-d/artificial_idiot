package google

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
)

const (
	SEARCH_ENDPOINT = "https://ajax.googleapis.com/ajax/services/search/images"
	VERSION         = "1.0"
)

type SearchResponse struct {
	ResponseData    *ResponseData `json:"responseData"`
	ResponseDetails string        `json:"responseDetails"`
	ResponseStatus  int           `json:"responseStatus"`
}

type ResponseData struct {
	Cursor  json.RawMessage `json:"cursor"`
	Results []*Result       `json:"results"`
}

type Result struct {
	Title   string `json:"title"`
	Url     string `json:"url"`
	Content string `json:"content"`
}

func Images(searchTerm string, animated bool) ([]string, error) {
	requestUrl, _ := url.Parse(SEARCH_ENDPOINT)

	query := requestUrl.Query()
	query.Set("v", VERSION)
	query.Set("q", searchTerm)
	query.Set("rsz", "8")
	query.Set("safe", "off")

	if animated {
		query.Set("imgtype", "animated")
	}

	requestUrl.RawQuery = query.Encode()

	resp, err := http.Get(requestUrl.String())

	if err != nil {
		return []string{}, err
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	searchResponse := &SearchResponse{}
	err = json.Unmarshal(body, searchResponse)

	if err != nil {
		return []string{}, err
	}

	imageUrls := make([]string, len(searchResponse.ResponseData.Results))
	for idx, result := range searchResponse.ResponseData.Results {
		imageUrls[idx] = result.Url
	}
	return imageUrls, nil
}

func RandomImage(query string) (string, error) {
	images, err := Images(query, false)

	if err != nil {
		return "", err
	}

	if len(images) == 0 {
		return "", errors.New("No images found")
	}

	return images[rand.Intn(len(images))], nil
}

func RandomGif(query string) (string, error) {
	images, err := Images(query, true)

	if err != nil {
		return "", err
	}

	if len(images) == 0 {
		return "", errors.New("No images found")
	}

	return images[rand.Intn(len(images))], nil
}
