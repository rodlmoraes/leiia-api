package handler

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"leiia/internal/entity"
	"leiia/internal/handler/response"

	"github.com/gin-gonic/gin"
	"github.com/ledongthuc/pdf"
	"gorm.io/gorm"
)

const (
	maxFileSize = 10 << 20 // 10 MB
)

// UploadResponse represents the upload response structure
type UploadResponse struct {
	ID         uint    `json:"id"`
	Message    string  `json:"message"`
	Filename   string  `json:"filename"`
	Size       int64   `json:"size"`
	Content    *string `json:"content,omitempty"`
	ParseError *string `json:"parse_error,omitempty"`
	UploadedAt string  `json:"uploaded_at"`
}

// UploadFile handles file upload requests
type UploadFile struct {
	DB *gorm.DB
}

// Handle processes the file upload request
func (h *UploadFile) Handle(c *gin.Context) {
	// Parse multipart form
	err := c.Request.ParseMultipartForm(maxFileSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewError("Failed to parse multipart form: "+err.Error()))
		return
	}

	// Get the file from form data
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewError("Failed to get file from form: "+err.Error()))
		return
	}

	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewError("Failed to open file: "+err.Error()))
		return
	}
	defer file.Close()

	// Validate file extension
	if !strings.HasSuffix(strings.ToLower(fileHeader.Filename), ".pdf") {
		c.JSON(http.StatusBadRequest, response.NewError("File must be a PDF"))
		return
	}

	// Validate content type
	if fileHeader.Header.Get("Content-Type") != "application/pdf" {
		c.JSON(http.StatusBadRequest, response.NewError("Content-Type must be application/pdf"))
		return
	}

	// Read file content into memory
	fileData, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewError("Failed to read file content: "+err.Error()))
		return
	}

	// Create File record
	filename := filepath.Base(fileHeader.Filename)
	fileRecord := entity.File{
		Filename:     filename,
		OriginalName: fileHeader.Filename,
		FileSize:     fileHeader.Size,
		ContentType:  fileHeader.Header.Get("Content-Type"),
		FileData:     fileData,
		UploadedAt:   time.Now(),
	}

	// Parse file content
	parsedText, parseErr := parseFileContent(fileData)
	if parseErr != nil {
		log.Printf("Warning: failed to parse file content: %v", parseErr)
		errorMsg := parseErr.Error()
		fileRecord.ParseError = &errorMsg
	} else {
		fileRecord.ParsedText = &parsedText
	}

	// Save to database
	result := h.DB.Create(&fileRecord)
	if result.Error != nil {
		log.Printf("Failed to save file to database: %v", result.Error)
		c.JSON(http.StatusInternalServerError, response.NewError("Failed to save file: "+result.Error.Error()))
		return
	}

	// Prepare response
	response := UploadResponse{
		ID:         fileRecord.ID,
		Filename:   fileRecord.Filename,
		Size:       fileRecord.FileSize,
		UploadedAt: fileRecord.UploadedAt.Format(time.RFC3339),
	}

	if fileRecord.ParseError != nil {
		response.Message = "File uploaded successfully but failed to parse content"
		response.ParseError = fileRecord.ParseError
	} else {
		response.Message = "File uploaded and parsed successfully"
		response.Content = fileRecord.ParsedText
	}

	// Send JSON response using Gin
	c.JSON(http.StatusOK, response)
}

// parseFileContent extracts text content from file data
func parseFileContent(fileData []byte) (string, error) {
	// Create a reader from the file data
	reader := bytes.NewReader(fileData)

	// Open file using the data reader
	fileReader, err := pdf.NewReader(reader, int64(len(fileData)))
	if err != nil {
		return "", fmt.Errorf("failed to create file reader: %w", err)
	}

	// Use GetPlainText() for human-readable text extraction
	var buf bytes.Buffer
	textReader, err := fileReader.GetPlainText()
	if err != nil {
		return "", fmt.Errorf("failed to extract text: %w", err)
	}

	_, err = buf.ReadFrom(textReader)
	if err != nil {
		return "", fmt.Errorf("failed to read text content: %w", err)
	}

	content := buf.String()
	if strings.TrimSpace(content) == "" {
		return "", fmt.Errorf("no text content found in file")
	}

	return content, nil
}
