package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func uploadhandler(c *gin.Context) {
	localPath := "data"
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err.Error())
		return
	}
	localPath = GetLocalPath(file.Filename)
	// Save the uploaded file locally
	if err := c.SaveUploadedFile(file, localPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Extract the compressed folder
	extractPath := GetExtractPath(localPath)
	if err := ExtractCompressedFolder(localPath, extractPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract compressed folder"})
		return
	}

	// Upload the extracted files to Azure Blob Storage
	if err := UploadFolderToAzure(extractPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload to Azure Blob Storage"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Folder uploaded successfully"})
}
