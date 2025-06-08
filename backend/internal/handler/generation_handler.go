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

	job := &service.Job{
		UserID:     userID.(int64),
		TemplateID: templateID,
		ImageData:  imageDataBytes,
	}

	if err := h.queueService.AddJob(c.Request.Context(), job); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add job to queue"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"request_id": "some-request-id"}) // Placeholder
} 