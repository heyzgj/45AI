package handler

import (
	"net/http"
	"strconv"

	"github.com/45ai/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type TemplateHandler interface {
	GetAll(c *gin.Context)
	GetByID(c *gin.Context)
}

type templateHandlerImpl struct {
	service service.TemplateService
}

func NewTemplateHandler(service service.TemplateService) TemplateHandler {
	return &templateHandlerImpl{service: service}
}

func (h *templateHandlerImpl) GetAll(c *gin.Context) {
	templates, err := h.service.GetAllTemplates(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, templates)
}

func (h *templateHandlerImpl) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template ID"})
		return
	}

	template, err := h.service.GetTemplateByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
		return
	}
	c.JSON(http.StatusOK, template)
} 