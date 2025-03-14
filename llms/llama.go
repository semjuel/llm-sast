package llms

type llamaModel struct {
}

func NewLLamaModel() LLMModel {
	return llamaModel{}
}

func (l llamaModel) Send(msg string) error {
	//TODO implement me
	panic("implement me")
}
