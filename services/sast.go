package services

import (
	"errors"
	"github.com/semjuel/llm-sast/llms"
)

type StaticAnalyzer interface {
	Analyze(dst string) (string, error)
}

func NewStaticAnalyzer(ext string, llm llms.LLMModel) (StaticAnalyzer, error) {
	switch ext {
	case ".apk", ".zip":
		return NewAndroidAnalyzer(llm), nil
	case ".ipa":
		return NewIOSAnalyzer(llm), nil
	}

	return nil, errors.New("unknown LLM model")
}
