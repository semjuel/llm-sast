package llms

import "github.com/semjuel/llm-sast/models"

type deepSeekModel struct {
}

func NewDeepSeekModel() LLMModel {
	return deepSeekModel{}
}

func (d deepSeekModel) Send(msg string) error {
	//TODO implement me
	panic("implement me")
}

func (d deepSeekModel) Name() string {
	//TODO implement me
	panic("implement me")
}

func (d deepSeekModel) AnalyzeUrl(filtered models.URLUsageFiltered) (models.URLFilteredResponse, error) {
	//TODO implement me
	panic("implement me")
}
