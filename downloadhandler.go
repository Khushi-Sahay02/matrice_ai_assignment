package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func downloadhandler(c *gin.Context) {
	blobName := "ksb.zip"

	containerURL := auth()

	// Create a blob URL
	blobURL := containerURL.NewBlobURL(blobName)
	fmt.Println(blobURL)

	// Download compressed data
	compressedData, err := downloadBlob(blobURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download compressed data"})
	}

	// Uncompress data
	uncompressedData, err := unzip(compressedData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unzip data"})
		return
	}

	// Save uncompressed data to local folder
	err = saveToLocalFolder(uncompressedData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save uncompressed data to local folder"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Download and unzip successful"})
}
