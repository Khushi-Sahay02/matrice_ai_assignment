package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func auth() azblob.ContainerURL {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	accountName := os.Getenv("ACCOUNTNAME")
	accountKey := os.Getenv("ACCOUNTKEY")
	containerName := os.Getenv("CONTAINERNAME")

	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		fmt.Println("error : Failed to create storage account credential")
		os.Exit(1)
	}

	// Create a pipeline using the credential
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	// Create a blob URL
	// Create a service URL
	serviceURL, err := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", "imagesmatriceai"))
	if err != nil {
		fmt.Println("error : Failed to parse service URL")
		os.Exit(1)
	}

	// Create a service URL using the parsed URL
	service := azblob.NewServiceURL(*serviceURL, p)

	// Create a container URL
	containerURL := service.NewContainerURL(containerName)

	return containerURL
}

func connectMongoDB() *mongo.Collection {
	// Replace with your MongoDB Atlas connection string
	uri := "mongodb+srv://sahaykhushi350:Khushi123@cluster0.9njbln9.mongodb.net/?retryWrites=true&w=majority"

	// Create a client
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB Atlas!")

	database := client.Database("imageinfo")
	collection := database.Collection("matriceai")
	return collection
}
