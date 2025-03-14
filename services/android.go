package services

import (
	"fmt"
	"github.com/semjuel/llm-sast/llms"
	"github.com/semjuel/llm-sast/services/android"
	"github.com/semjuel/llm-sast/utils"
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
	dest := fmt.Sprintf("uploads/%s", utils.HashString(src))
	err := android.UnzipAPK(src, dest)
	if err != nil {
		return "", err
	}

	err = android.Apk2Java(src, dest)
	if err != nil {
		return "", err
	}
	// apk_2_java
	// dex_2_smali
	// code_an_dic

	//TODO implement me
	return "Android success", nil
}
