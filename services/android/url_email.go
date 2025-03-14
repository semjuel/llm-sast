package android

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
var emailRegex = regexp.MustCompile(`[\w+.-]{1,20}@[\w-]{1,20}\.[\w]{2,10}`)

// ExtractFromSource - function to extract data from Java/Kotlin files
func ExtractFromSource(srcDir string, skipPaths []string) []models.URLUsageFiltered {
	msg := "Extracting Emails and URLs from Source Code!!!"
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
		if strings.HasSuffix(info.Name(), ".java") || strings.HasSuffix(info.Name(), ".kt") {
			fileContent, err := readFileContent(path)
			if err != nil {
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
		}
		return nil
	})

	if err != nil {
		log.Println("Error while walking through files:", err)
	}

	msg = "Email and URL Extraction Completed"
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
	emailMatches := emailRegex.FindAllString(strings.ToLower(content), -1)
	for _, email := range emailMatches {
		if !strings.HasPrefix(email, "//") { // Ignore commented-out emails
			emails = append(emails, email)
		}
	}

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
	for scanner.Scan() {
		content.WriteString(scanner.Text() + "\n")
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return content.String(), nil
}
