package services

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/semjuel/llm-sast/llms"
	"github.com/semjuel/llm-sast/models"
	"github.com/semjuel/llm-sast/services/ios"
	"github.com/semjuel/llm-sast/utils"
	"log"
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

	urlsData := ios.ExtractFromSource(dest, []string{})
	unique := uniqueByFilepath(urlsData)
	for _, urlRow := range unique {
		log.Println(urlRow.Filepath)
		log.Println("---START sending requests---")
		// Send request to the LLM and analyze the source
		res, err := i.llm.AnalyzeUrl(urlRow)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Println("---END sending requests---")
		response = append(response, res)
	}

	log.Println("Responses")
	log.Println(response)

	return response, nil
}
