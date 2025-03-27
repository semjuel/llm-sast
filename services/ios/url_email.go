package ios

import (
	"bufio"
	"github.com/semjuel/llm-sast/models"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Regular Expressions
var urlRegex = regexp.MustCompile(`((?:https?://|s?ftps?://|file://|javascript:|data:|www\d{0,3}[.])[\w().=/;,#:@?&~*+!$%{}-]+)`)

func ExtractFromSource(srcDir string, skipPaths []string) []models.URLUsageFiltered {
	msg := "Extracting URLs from Source Code!!!"
	log.Println(msg)

	var results []models.URLUsageFiltered

	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip directories specified in `skipPaths`
		for _, skip := range skipPaths {
			if strings.Contains(path, skip) {
				return nil
			}
		}

		// Process only .java and .kt files
		if strings.HasSuffix(info.Name(), ".m") ||
			strings.HasSuffix(info.Name(), ".swift") {
			fileContent, err := readFileContent(path)
			if err != nil {
				log.Println(path, err.Error())
				return nil // Skip files that cannot be read
			}

			relativePath, _ := filepath.Rel(srcDir, path)
			// @TODO add email usage
			urls, _ := extractEmailsAndURLs(fileContent, relativePath)
			for _, url := range urls {
				result := models.URLUsageFiltered{
					Url:      url,
					Filepath: relativePath,
					Content:  fileContent,
				}

				results = append(results, result)
			}
		} else if strings.HasSuffix(info.Name(), ".so") {
			// @TODO implement this, A *.so file is a binary shared object
			// use hex dump or debug/elf or disassembler or decompiler
		}
		return nil
	})

	if err != nil {
		log.Println("Error while walking through files:", err)
	}

	msg = "URL Extraction Completed"
	log.Println(msg)

	return results
}

// Function to extract URLs and Emails from a file
func extractEmailsAndURLs(content, relativePath string) ([]string, []string) {
	var urls, emails []string

	// Extract URLs
	matches := urlRegex.FindAllString(strings.ToLower(content), -1)
	for _, url := range matches {
		urls = append(urls, url)
	}

	// Extract Emails
	// @TODO implement this

	return urls, emails
}

// Read file content safely
func readFileContent(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var content strings.Builder
	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 1024*1024)
	scanner.Buffer(buf, 20*1024*1024)
	for scanner.Scan() {
		content.WriteString(scanner.Text() + "\n")
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return content.String(), nil
}
