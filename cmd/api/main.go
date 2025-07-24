package main

import (
	"leiia/internal/database"
	"leiia/internal/handler"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	db, err := database.Initialize()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Create Gin router
	r := gin.Default()

	// Create handler instances
	uploadFileHandler := &handler.UploadFile{
		DB: db,
	}
	getFileHandler := &handler.GetFile{
		DB: db,
	}
	healthHandler := &handler.Health{
		DB: db,
	}

	// Setup routes
	r.GET("/health", healthHandler.Handle)

	fileGroup := r.Group("/file")
	fileGroup.POST("/upload", uploadFileHandler.Handle)
	fileGroup.GET("/:id", getFileHandler.Handle)

	// Start server
	log.Fatal(r.Run(":8080"))
}
