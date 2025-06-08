package handler

import (
	"net/http"
	"strconv"

	"github.com/45ai/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	GetProfile(c *gin.Context)
	UpdateProfile(c *gin.Context)
	GetTransactions(c *gin.Context)
}

type userHandlerImpl struct {
	userService        service.UserService
	transactionService service.TransactionService
}

func NewUserHandler(userService service.UserService, transactionService service.TransactionService) UserHandler {
	return &userHandlerImpl{
		userService:        userService,
		transactionService: transactionService,
	}
}

func (h *userHandlerImpl) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve user profile"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *userHandlerImpl) UpdateProfile(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
} 