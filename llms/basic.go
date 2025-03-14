package llms

import (
	"errors"
)

const (
	LlamaModel     = "llama"
	ChatGPTO1Model = "chatgpto1"
	DeepSeekModel  = "deepseek"
)

type LLMModel interface {
	Send(msg string) error
}

func NewLLMModel(
	name string,
) (LLMModel, error) {
	switch name {
	case LlamaModel:
		return NewLLamaModel(), nil
	case ChatGPTO1Model:
		return NewGPTModel(), nil
	case DeepSeekModel:
		return NewDeepSeekModel(), nil
	}

	return nil, errors.New("unknown LLM model")
}
