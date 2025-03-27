package llms

import (
	"fmt"
	"github.com/semjuel/llm-sast/assets"
	"github.com/semjuel/llm-sast/models"
	"os"
	"strings"
)

type promptModel struct {
}

func NewPromptModel() LLMModel {
	return promptModel{}
}

func (p promptModel) Name() string {
	return "prompt"
}

func (p promptModel) Send(msg string) error {
	//TODO implement me
	panic("implement me")
}

func (p promptModel) AnalyzeUrl(filtered models.URLUsageFiltered) (models.URLFilteredResponse, error) {
	result := models.URLFilteredResponse{}

	dest := fmt.Sprintf("prompts/%s.txt", strings.Replace(filtered.Filepath, "/", "_", -1))

	combined := string(assets.AndroidPrompt) + filtered.Content
	err := os.WriteFile(dest, []byte(combined), 0644)
	if err != nil {
		return result, err
	}

	return result, nil
}
