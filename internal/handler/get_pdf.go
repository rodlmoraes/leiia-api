package handler

import (
	"net/http"
	"time"

	"leiia/internal/entity"
	"leiia/internal/handler/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetResponse represents the get PDF response structure
type GetResponse struct {
	ID         uint    `json:"id"`
	Message    string  `json:"message"`
	Filename   string  `json:"filename"`
	Size       int64   `json:"size"`
	Content    *string `json:"content,omitempty"`
	ParseError *string `json:"parse_error,omitempty"`
	UploadedAt string  `json:"uploaded_at"`
}

// GetPDF handles PDF retrieval requests
type GetPDF struct {
	DB *gorm.DB
}

// Handle processes the PDF retrieval request
func (h *GetPDF) Handle(c *gin.Context) {
	// Extract ID from URL parameter
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.NewError("PDF ID is required"))
		return
	}

	// Find PDF file in database
	var pdfFile entity.PDFFile
	result := h.DB.First(&pdfFile, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, response.NewError("PDF not found"))
		} else {
			c.JSON(http.StatusInternalServerError, response.NewError("Database error: "+result.Error.Error()))
		}
		return
	}

	// Prepare response
	response := GetResponse{
		ID:         pdfFile.ID,
		Filename:   pdfFile.Filename,
		Size:       pdfFile.FileSize,
		UploadedAt: pdfFile.UploadedAt.Format(time.RFC3339),
	}

	if pdfFile.ParseError != nil {
		response.Message = "PDF found but parsing failed"
		response.ParseError = pdfFile.ParseError
	} else {
		response.Message = "PDF found and parsed successfully"
		response.Content = pdfFile.ParsedText
	}

	// Send JSON response using Gin
	c.JSON(http.StatusOK, response)
}
