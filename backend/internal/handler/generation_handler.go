package handler

import (
	"net/http"
	"strconv"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/45ai/backend/internal/service"
)

type GenerationHandler interface {
	GenerateImage(c *gin.Context)
	GetGenerationStatus(c *gin.Context)
	GetGenerationResult(c *gin.Context)
}

type generationHandlerImpl struct {
	service      service.GenerationService
	queueService service.QueueService
}

func NewGenerationHandler(service service.GenerationService, queueService service.QueueService) GenerationHandler {
	return &generationHandlerImpl{service: service, queueService: queueService}
}

func (h *generationHandlerImpl) GenerateImage(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	templateID, err := strconv.Atoi(c.PostForm("template_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template_id"})
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image file is required"})
		return
	}

	imageData, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open image"})
		return
	}
	defer imageData.Close()

	// Read image data into a byte slice
	imageDataBytes, err := io.ReadAll(imageData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read image data"})
		return
	}

	// Add job to queue and get job ID
	jobID, err := h.queueService.AddJob(c.Request.Context(), userID.(int64), templateID, imageDataBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add job to queue"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"job_id": jobID,
		"status": "pending",
		"message": "Image generation job queued successfully",
	})
}

func (h *generationHandlerImpl) GetGenerationStatus(c *gin.Context) {
	jobID := c.Param("job_id")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "job_id is required"})
		return
	}

	status, err := h.service.GetGenerationStatus(c.Request.Context(), jobID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}

	c.JSON(http.StatusOK, status)
}

func (h *generationHandlerImpl) GetGenerationResult(c *gin.Context) {
	jobID := c.Param("job_id")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "job_id is required"})
		return
	}

	result, err := h.service.GetGenerationResult(c.Request.Context(), jobID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
} 