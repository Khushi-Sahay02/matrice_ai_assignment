package main

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

func UploadFolderToAzure(folderPath string) error {
	containerURL := auth()
	// Walk through the folder and upload files to Azure Blob Storage
	err := filepath.Walk(folderPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if !info.IsDir() {
			relativePath, err := filepath.Rel(folderPath, filePath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			// Create a blob URL
			blobURL := containerURL.NewBlockBlobURL(relativePath)

			// Open the file for reading
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer file.Close()

			// Upload the file to Azure Blob Storage
			_, err = azblob.UploadFileToBlockBlob(context.Background(), file, blobURL,
				azblob.UploadToBlockBlobOptions{
					BlockSize:   4 * 1024 * 1024, // 4MB block size
					Parallelism: 16,
				})

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		return nil
	})

	return err
}

// ExtractCompressedFolder extracts the compressed folder to the specified path
func ExtractCompressedFolder(source, destination string) error {
	// Open the source file for reading
	file, err := os.Open(source)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a gzip reader
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	// Create a tar reader
	tarReader := tar.NewReader(gzipReader)

	// Iterate over each file in the tar archive
	for {
		header, err := tarReader.Next()

		// Check for the end of the tar archive
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		// Determine the file path for the extracted file
		target := filepath.Join(destination, header.Name)

		// Create directories as needed
		if header.FileInfo().IsDir() {
			if err := os.MkdirAll(target, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		// Create the file
		file, err := os.Create(target)
		if err != nil {
			return err
		}
		defer file.Close()

		// Copy the file contents from the tar archive to the new file
		if _, err := io.Copy(file, tarReader); err != nil {
			return err
		}
	}

	return nil
}
