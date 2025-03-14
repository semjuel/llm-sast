package llms

type deepSeekModel struct {
}

func NewDeepSeekModel() LLMModel {
	return deepSeekModel{}
}

func (d deepSeekModel) Send(msg string) error {
	//TODO implement me
	panic("implement me")
}
