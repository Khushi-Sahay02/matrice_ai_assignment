package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/download", downloadhandler)
	r.POST("/upload", uploadhandler)
	r.GET("/converter", converter)
	r.Run(":8080")
}
