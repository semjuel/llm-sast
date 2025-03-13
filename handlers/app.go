package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"strings"
)

func Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get uploaded file: " + err.Error(),
		})
		return
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExtensions := map[string]bool{
		".apk": true,
		".ipa": true,
		".zip": true,
	}
	if !allowedExtensions[ext] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Only .apk, .ipa, and .zip files are allowed",
		})
		return
	}

	// Create a destination path and save the uploaded file
	dst := filepath.Join("uploads", file.Filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save file: " + err.Error(),
		})
		return
	}

	// @TODO

	// @TODO change the response
	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded successfully",
		"filename": file.Filename,
	})
}
