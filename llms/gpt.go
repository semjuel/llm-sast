package llms

import "github.com/semjuel/llm-sast/models"

type gptModel struct {
}

func NewGPTModel() LLMModel {
	return gptModel{}
}

func (g gptModel) Send(msg string) error {
	//TODO implement me
	panic("implement me")
}

func (g gptModel) Name() string {
	//TODO implement me
	panic("implement me")
}

func (g gptModel) AnalyzeUrl(filtered models.URLUsageFiltered) (models.URLFilteredResponse, error) {
	//TODO implement me
	panic("implement me")
}
