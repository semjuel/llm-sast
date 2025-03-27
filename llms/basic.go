package llms

import (
	"errors"
	"fmt"
	"github.com/semjuel/llm-sast/models"
	"os"
)

const (
	PromptModel    = "prompt"
	LlamaModel     = "llama"
	GemmaModel     = "gemma"
	ChatGPTO1Model = "chatgpto1"
	DeepSeekModel  = "deepseek"
)

type LLMModel interface {
	Name() string
	AnalyzeUrl(filtered models.URLUsageFiltered) (models.URLFilteredResponse, error)
	Send(msg string) error
}

func NewLLMModel(
	name string,
) (LLMModel, error) {
	// @TODO move this to another place
	host := os.Getenv("OPEN_WEB_UI_URL")
	if host == "" {
		return nil, errors.New("host environment variable for LLM api not set")
	}
	url := fmt.Sprintf("%s/api/chat/completions", host)
	token := os.Getenv("OPEN_WEB_UI_TOKEN")

	switch name {
	case PromptModel:
		return NewPromptModel(), nil
	case LlamaModel:
		return NewLLamaModel(url, token), nil
	case GemmaModel:
		return NewGemmaModel(url, token), nil
		//case ChatGPTO1Model:
		//	return NewGPTModel(), nil
		//case DeepSeekModel:
		//	return NewDeepSeekModel(), nil
	}

	return nil, errors.New("unknown LLM model")
}
