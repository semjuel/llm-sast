package services

import (
	"fmt"
	"github.com/semjuel/llm-sast/llms"
	"github.com/semjuel/llm-sast/models"
	"github.com/semjuel/llm-sast/services/android"
	"github.com/semjuel/llm-sast/utils"
	"log"
)

type androidAnalyzer struct {
	llm llms.LLMModel
}

func NewAndroidAnalyzer(llm llms.LLMModel) StaticAnalyzer {
	return androidAnalyzer{
		llm: llm,
	}
}

func (a androidAnalyzer) Analyze(src string) ([]models.URLFilteredResponse, error) {
	var response []models.URLFilteredResponse

	dest := fmt.Sprintf("uploads/%s", utils.HashString(src))
	err := android.UnzipAPK(src, dest)
	if err != nil {
		return response, err
	}

	err = android.Apk2Java(src, dest)
	if err != nil {
		return response, err
	}
	// apk_2_java
	// dex_2_smali
	// code_an_dic

	// @TODO this handles only url for now
	urlsData := android.ExtractFromSource(dest, []string{})
	unique := uniqueByFilepath(urlsData)
	for _, urlRow := range unique {
		log.Println(urlRow.Filepath)
		log.Println("---START sending requests---")
		// Send request to the LLM and analyze the source
		res, err := a.llm.AnalyzeUrl(urlRow)
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

// @TODO move this function to another place
// uniqueByFilepath returns a slice containing only the unique URLUsageFiltered entries based on Filepath.
func uniqueByFilepath(usages []models.URLUsageFiltered) []models.URLUsageFiltered {
	seen := make(map[string]bool)
	var unique []models.URLUsageFiltered

	for _, usage := range usages {
		if !seen[usage.Filepath] {
			seen[usage.Filepath] = true
			unique = append(unique, usage)
		}
	}

	return unique
}
