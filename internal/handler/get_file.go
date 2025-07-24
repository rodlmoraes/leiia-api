package handler

import (
	"net/http"
	"time"

	"leiia/internal/entity"
	"leiia/internal/handler/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetResponse represents the get file response structure
type GetResponse struct {
	ID         uint    `json:"id"`
	Message    string  `json:"message"`
	Filename   string  `json:"filename"`
	Size       int64   `json:"size"`
	Content    *string `json:"content,omitempty"`
	ParseError *string `json:"parse_error,omitempty"`
	UploadedAt string  `json:"uploaded_at"`
}

// GetFile handles file retrieval requests
type GetFile struct {
	DB *gorm.DB
}

// Handle processes the file retrieval request
func (h *GetFile) Handle(c *gin.Context) {
	// Extract ID from URL parameter
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.NewError("File ID is required"))
		return
	}

	// Find file in database
	var fileRecord entity.File
	result := h.DB.First(&fileRecord, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, response.NewError("File not found"))
		} else {
			c.JSON(http.StatusInternalServerError, response.NewError("Database error: "+result.Error.Error()))
		}
		return
	}

	// Prepare response
	response := GetResponse{
		ID:         fileRecord.ID,
		Filename:   fileRecord.Filename,
		Size:       fileRecord.FileSize,
		UploadedAt: fileRecord.UploadedAt.Format(time.RFC3339),
	}

	if fileRecord.ParseError != nil {
		response.Message = "File found but parsing failed"
		response.ParseError = fileRecord.ParseError
	} else {
		response.Message = "File found and parsed successfully"
		response.Content = fileRecord.ParsedText
	}

	// Send JSON response using Gin
	c.JSON(http.StatusOK, response)
}
