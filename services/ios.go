package services

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/semjuel/llm-sast/llms"
	"os"
)

type iosAnalyzer struct {
	llm llms.LLMModel
}

func NewIOSAnalyzer(llm llms.LLMModel) StaticAnalyzer {
	return iosAnalyzer{
		llm: llm,
	}
}

func (i iosAnalyzer) Analyze(src string) (string, error) {
	b, err := os.ReadFile(src)
	if err != nil {
		return "", err
	}

	reader := bytes.NewReader(b)
	zipReader, err := zip.NewReader(reader, int64(len(b)))
	if err != nil {
		return "", err
	}

	dest := fmt.Sprintf("uploads/%s", hashString(src))
	err = unzip(zipReader, dest)
	if err != nil {
		return "", err
	}

	//TODO implement me
	return "IOS success", nil
}
