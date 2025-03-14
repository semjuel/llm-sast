package services

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/semjuel/llm-sast/llms"
	"github.com/semjuel/llm-sast/models"
	"github.com/semjuel/llm-sast/utils"
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

func (i iosAnalyzer) Analyze(src string) ([]models.URLFilteredResponse, error) {
	var response []models.URLFilteredResponse

	b, err := os.ReadFile(src)
	if err != nil {
		return response, err
	}

	reader := bytes.NewReader(b)
	zipReader, err := zip.NewReader(reader, int64(len(b)))
	if err != nil {
		return response, err
	}

	dest := fmt.Sprintf("uploads/%s", utils.HashString(src))
	err = utils.Unzip(zipReader, dest)
	if err != nil {
		return response, err
	}

	//TODO implement me
	return response, nil
}
