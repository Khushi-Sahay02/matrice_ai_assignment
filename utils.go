package main

import (
	"context"
	"time"

	"path/filepath"

	"go.mongodb.org/mongo-driver/mongo"
)

func GetLocalPath(filename string) string {
	return filepath.Join("uploads", filename)
}

func GetExtractPath(localPath string) string {
	return filepath.Join("extracted", localPath)
}

func SaveToMongoDB(collection *mongo.Collection, image Image) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, image)
	return err
}
