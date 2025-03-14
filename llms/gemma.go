package llms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/semjuel/llm-sast/models"
	"io"
	"log"
	"net/http"
)

type gemmaModel struct {
	url   string
	token string
}

func NewGemmaModel(url, token string) LLMModel {
	return gemmaModel{
		url:   url,
		token: token,
	}
}

func (g gemmaModel) Name() string {
	return "gemma3:27b"
}

func (g gemmaModel) Send(msg string) error {
	//TODO implement me
	panic("implement me")
}

func (g gemmaModel) AnalyzeUrl(filtered models.URLUsageFiltered) (models.URLFilteredResponse, error) {
	result := models.URLFilteredResponse{}

	question := "Analyze the code below again and give an answer to whether requests for links from this file will be sent." +
		" Provide the answer in json format in this form:" +
		"\n[\n{\n\"url\": \"http://...\",\n\"request\": true,\n\"description\": \"text\"\n}\n]\nwhere \"URL\" - " +
		"is a found URL\n\"request\" - could be \"false\" - if request will not be sent, \"true\" - " +
		"if request will be sent, \"nil\" - if you are not sure\n\"description\" - " +
		"you argumentation why you made such decision, the answer shouldn't be longer " +
		"then 2 sentences.\nYour response must include all urls found in the code.\nHere is a code to analyze:\n"
	question += filtered.Content

	reqBody := models.ChatCompletionRequest{
		Model: g.Name(),
		Messages: []models.Message{
			{
				Role:    "user",
				Content: question,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return result, err
	}

	req, err := http.NewRequest("POST", g.url, bytes.NewBuffer(jsonData))
	if err != nil {
		return result, err
	}

	// Set the appropriate headers.
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", g.token))
	req.Header.Set("Content-Type", "application/json")

	// Use http.Client to send the request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	// Read and print the response.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	log.Println(string(body))

	return result, nil
}
