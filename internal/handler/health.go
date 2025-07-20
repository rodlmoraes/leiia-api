package handler

import (
	"net/http"

	"leiia/internal/handler/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthResponse represents the health check response structure
type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

// Health handles health check requests
type Health struct {
	DB *gorm.DB
}

// Handle processes the health check request
func (h *Health) Handle(c *gin.Context) {
	// Check database connectivity
	sqlDB, err := h.DB.DB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewError("Database connection error"))
		return
	}

	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewError("Database ping failed"))
		return
	}

	response := HealthResponse{
		Status:   "healthy",
		Database: "connected",
	}
	c.JSON(http.StatusOK, response)
}
