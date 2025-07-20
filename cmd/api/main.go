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
	uploadPDFHandler := &handler.UploadPDF{
		DB: db,
	}
	getPDFHandler := &handler.GetPDF{
		DB: db,
	}
	healthHandler := &handler.Health{
		DB: db,
	}

	// Setup routes
	r.GET("/health", healthHandler.Handle)

	pdfGroup := r.Group("/pdf")
	pdfGroup.POST("/upload", uploadPDFHandler.Handle)
	pdfGroup.GET("/:id", getPDFHandler.Handle)

	// Start server
	log.Fatal(r.Run(":8080"))
}
