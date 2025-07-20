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

// UploadPDF handles PDF upload requests
type UploadPDF struct {
	DB *gorm.DB
}

// Handle processes the PDF upload request
func (h *UploadPDF) Handle(c *gin.Context) {
	// Parse multipart form
	err := c.Request.ParseMultipartForm(maxFileSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewError("Failed to parse multipart form: "+err.Error()))
		return
	}

	// Get the file from form data
	fileHeader, err := c.FormFile("pdf")
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

	// Create PDFFile record
	filename := filepath.Base(fileHeader.Filename)
	pdfFile := entity.PDFFile{
		Filename:     filename,
		OriginalName: fileHeader.Filename,
		FileSize:     fileHeader.Size,
		ContentType:  fileHeader.Header.Get("Content-Type"),
		FileData:     fileData,
		UploadedAt:   time.Now(),
	}

	// Parse PDF content
	parsedText, parseErr := parsePDFContent(fileData)
	if parseErr != nil {
		log.Printf("Warning: failed to parse PDF content: %v", parseErr)
		errorMsg := parseErr.Error()
		pdfFile.ParseError = &errorMsg
	} else {
		pdfFile.ParsedText = &parsedText
	}

	// Save to database
	result := h.DB.Create(&pdfFile)
	if result.Error != nil {
		log.Printf("Failed to save PDF to database: %v", result.Error)
		c.JSON(http.StatusInternalServerError, response.NewError("Failed to save file: "+result.Error.Error()))
		return
	}

	// Prepare response
	response := UploadResponse{
		ID:         pdfFile.ID,
		Filename:   pdfFile.Filename,
		Size:       pdfFile.FileSize,
		UploadedAt: pdfFile.UploadedAt.Format(time.RFC3339),
	}

	if pdfFile.ParseError != nil {
		response.Message = "PDF uploaded successfully but failed to parse content"
		response.ParseError = pdfFile.ParseError
	} else {
		response.Message = "PDF uploaded and parsed successfully"
		response.Content = pdfFile.ParsedText
	}

	// Send JSON response using Gin
	c.JSON(http.StatusOK, response)
}

// parsePDFContent extracts text content from PDF file data
func parsePDFContent(fileData []byte) (string, error) {
	// Create a reader from the file data
	reader := bytes.NewReader(fileData)

	// Open PDF using the data reader
	pdfReader, err := pdf.NewReader(reader, int64(len(fileData)))
	if err != nil {
		return "", fmt.Errorf("failed to create PDF reader: %w", err)
	}

	// Use GetPlainText() for human-readable text extraction
	var buf bytes.Buffer
	textReader, err := pdfReader.GetPlainText()
	if err != nil {
		return "", fmt.Errorf("failed to extract text: %w", err)
	}

	_, err = buf.ReadFrom(textReader)
	if err != nil {
		return "", fmt.Errorf("failed to read text content: %w", err)
	}

	content := buf.String()
	if strings.TrimSpace(content) == "" {
		return "", fmt.Errorf("no text content found in PDF")
	}

	return content, nil
}
