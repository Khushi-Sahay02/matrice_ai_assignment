package main

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

func downloadBlob(blobURL azblob.BlobURL) ([]byte, error) {
	getBlobResponse, err := blobURL.Download(context.Background(), 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false, azblob.ClientProvidedKeyOptions{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer getBlobResponse.Body(azblob.RetryReaderOptions{}).Close()
	// Read the content into a byte slice
	compressedData, err := io.ReadAll(getBlobResponse.Body(azblob.RetryReaderOptions{}))
	if err != nil {
		return nil, err
	}

	return compressedData, nil
}

func unzip(compressedData []byte) ([]byte, error) {
	r, err := zip.NewReader(strings.NewReader(string(compressedData)), int64(len(compressedData)))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var result []byte

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer rc.Close()

		content, err := io.ReadAll(rc)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		result = append(result, content...)
	}

	return result, nil
}

func saveToLocalFolder(data []byte) error {
	localFolderPath := "data"
	if _, err := os.Stat(localFolderPath); os.IsNotExist(err) {
		err := os.Mkdir(localFolderPath, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	localFilePath := filepath.Join(localFolderPath, "uncompressed_data.jpg")
	err := os.WriteFile(localFilePath, data, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return nil
}
