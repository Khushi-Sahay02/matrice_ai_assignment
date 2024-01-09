package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/cheggaaa/pb.v1"
)

// BBox represents a bounding box in COCO format
type BBox struct {
	X float64
	Y float64
	W float64
	H float64
}

// Category represents a COCO category
type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Image represents a COCO image
type Image struct {
	ID       int    `json:"id"`
	FileName string `json:"file_name"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

// Annotation represents a COCO annotation
type Annotation struct {
	ImageID    int       `json:"image_id"`
	CategoryID int       `json:"category_id"`
	BBox       []float64 `json:"bbox"`
	Size       []float64 `json:"size"`
}

// COCOData represents the COCO JSON structure
type COCOData struct {
	Images      []Image      `json:"images"`
	Annotations []Annotation `json:"annotations"`
	Categories  []Category   `json:"categories"`
}

func convertBBoxCOCO2Yolo(imgWidth, imgHeight int, bbox []float64) BBox {
	dw := 1.0 / float64(imgWidth)
	dh := 1.0 / float64(imgHeight)

	x, y, w, h := bbox[0], bbox[1], bbox[2], bbox[3]

	xCenter := x + w/2.0
	yCenter := y + h/2.0

	xRel := xCenter * dw
	yRel := yCenter * dh
	wRel := w * dw
	hRel := h * dh

	return BBox{X: xRel, Y: yRel, W: wRel, H: hRel}
}

func convertCOCOJsonToYoloTxt(collection *mongo.Collection, outputPath, jsonFile string) {

	jsonData := readCOCOJson(jsonFile)

	labelFile := filepath.Join(outputPath, "_darknet.labels")
	writeLabels(labelFile, jsonData.Categories)

	bar := pb.StartNew(len(jsonData.Images)).Prefix("Annotation txt for each image")

	for _, image := range jsonData.Images {
		annoInImage := getAnnotationsInImage(jsonData.Annotations, image.ID)
		annoTxt := filepath.Join(outputPath, image.FileName[:len(image.FileName)-len(filepath.Ext(image.FileName))]+".txt")

		writeAnnotations(annoTxt, image.Width, image.Height, annoInImage)
		err := SaveToMongoDB(collection, image)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		bar.Increment()
	}

	bar.Finish()
	fmt.Println("Converting COCO Json to YOLO txt finished!")
}

func readCOCOJson(jsonFile string) COCOData {
	file, err := os.Open(jsonFile)
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		os.Exit(1)
	}
	defer file.Close()

	var jsonData COCOData
	err = json.NewDecoder(file).Decode(&jsonData)
	if err != nil {
		fmt.Println("Error decoding JSON data:", err)
		os.Exit(1)
	}

	return jsonData
}

func writeLabels(labelFile string, categories []Category) {
	sort.Slice(categories, func(i, j int) bool {
		return categories[i].ID < categories[j].ID
	})

	file, err := os.Create(labelFile)
	if err != nil {
		fmt.Println("Error creating labels file:", err)
		os.Exit(1)
	}
	defer file.Close()

	for _, category := range categories {
		file.WriteString(category.Name + "\n")
	}
}

func getAnnotationsInImage(annotations []Annotation, imageID int) []Annotation {
	var annoInImage []Annotation
	for _, anno := range annotations {
		if anno.ImageID == imageID {
			annoInImage = append(annoInImage, anno)
		}
	}
	return annoInImage
}

func writeAnnotations(annoFile string, imgWidth, imgHeight int, annotations []Annotation) {
	file, err := os.Create(annoFile)
	if err != nil {
		fmt.Println("Error creating annotation file:", err)
		os.Exit(1)
	}
	defer file.Close()

	for _, anno := range annotations {
		category := strconv.Itoa(anno.CategoryID)
		bbox := convertBBoxCOCO2Yolo(imgWidth, imgHeight, anno.BBox)
		file.WriteString(fmt.Sprintf("%s %.6f %.6f %.6f %.6f\n", category, bbox.X, bbox.Y, bbox.W, bbox.H))
	}
}

func converter(c *gin.Context) {
	outputPath := "labels"
	jsonFile := "instances_val2017.json"
	collection := connectMongoDB()
	convertCOCOJsonToYoloTxt(collection, outputPath, jsonFile)
}
