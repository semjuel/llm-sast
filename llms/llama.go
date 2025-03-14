package llms

import "github.com/semjuel/llm-sast/models"

type llamaModel struct {
	url   string
	token string
}

func NewLLamaModel(url, token string) LLMModel {
	return llamaModel{
		url:   url,
		token: token,
	}
}

func (l llamaModel) Name() string {
	return "llama3.3:latest"
}

func (l llamaModel) Send(msg string) error {
	//TODO implement me
	panic("implement me")
}

func (l llamaModel) AnalyzeUrl(filtered models.URLUsageFiltered) (models.URLFilteredResponse, error) {
	//TODO implement me
	panic("implement me")
}
