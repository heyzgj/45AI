package handler

import (
	"net/http"
	"strconv"
	"fmt"

	"github.com/45ai/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/45ai/backend/internal/model"
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

func (h *userHandlerImpl) GetTransactions(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// Parse pagination parameters
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 10 // Default limit
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0 // Default offset
	}

	// Get transactions
	transactions, err := h.transactionService.GetTransactionsByUserID(c.Request.Context(), userID.(int64), limit, offset)
	if err != nil {
		// Log the actual error for debugging
		c.Header("X-Debug-Error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve transactions", "debug": err.Error()})
		return
	}

	// Debug logging
	c.Header("X-Debug-Transactions-Nil", fmt.Sprintf("%v", transactions == nil))
	c.Header("X-Debug-Transactions-Len", fmt.Sprintf("%d", len(transactions)))

	// Get total count for pagination
	totalCount, err := h.transactionService.CountTransactionsByUserID(c.Request.Context(), userID.(int64))
	if err != nil {
		// If count fails, continue without it
		totalCount = 0
	}

	// Force empty array instead of null
	var transactionsArray []model.Transaction
	if transactions != nil && len(transactions) > 0 {
		transactionsArray = transactions
	} else {
		transactionsArray = make([]model.Transaction, 0)
	}

	// Return response with pagination info
	c.JSON(http.StatusOK, gin.H{
		"transactions": transactionsArray,
		"pagination": gin.H{
			"limit":       limit,
			"offset":      offset,
			"total_count": totalCount,
		},
	})
} 