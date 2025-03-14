package services

import (
	"github.com/semjuel/llm-sast/llms"
	"os"
	"path/filepath"
)

type androidAnalyzer struct {
	llm llms.LLMModel
}

func NewAndroidAnalyzer(llm llms.LLMModel) StaticAnalyzer {
	return androidAnalyzer{
		llm: llm,
	}
}

func (a androidAnalyzer) Analyze(src string) (string, error) {
	b, err := os.ReadFile(src)
	if err != nil {
		return "", err
	}

	ext := filepath.Ext(src)
	if ext == ".zip" {
		// Merge apk file inside the zip archive
		b, err = mergeApks(src)
		if err != nil {
			return "", err
		}
	}

	err = unzipAPK(b, "uploads")
	if err != nil {
		return "", err
	}

	//TODO implement me
	return "Android success", nil
}
