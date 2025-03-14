package llms

type gptModel struct {
}

func NewGPTModel() LLMModel {
	return gptModel{}
}

func (g gptModel) Send(msg string) error {
	//TODO implement me
	panic("implement me")
}
