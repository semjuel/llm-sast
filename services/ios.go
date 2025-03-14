package services

import "github.com/semjuel/llm-sast/llms"

type iosAnalyzer struct {
	llm llms.LLMModel
}

func NewIOSAnalyzer(llm llms.LLMModel) StaticAnalyzer {
	return iosAnalyzer{
		llm: llm,
	}
}

func (i iosAnalyzer) Analyze(file string) (string, error) {
	//TODO implement me
	return "IOS success", nil
}
