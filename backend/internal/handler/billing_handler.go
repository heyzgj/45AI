package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/45ai/backend/internal/service"
	"github.com/45ai/backend/internal/model"
)

// BillingHandler handles billing-related HTTP requests
type BillingHandler struct {
	wechatPayService    service.WechatPayService
	appleIAPService     service.AppleIAPService
	creditService       service.CreditService
	transactionService  service.TransactionService
}

// NewBillingHandler creates a new billing handler
func NewBillingHandler(
	wechatPayService service.WechatPayService,
	appleIAPService service.AppleIAPService,
	creditService service.CreditService,
	transactionService service.TransactionService,
) *BillingHandler {
	return &BillingHandler{
		wechatPayService:   wechatPayService,
		appleIAPService:    appleIAPService,
		creditService:      creditService,
		transactionService: transactionService,
	}
}

// PurchaseRequest represents a unified purchase validation request
type PurchaseRequest struct {
	ProductID   string `json:"product_id" binding:"required"`
	Platform    string `json:"platform" binding:"required"` // "wechat" or "apple"
	Receipt     string `json:"receipt" binding:"required"`
	Description string `json:"description"`
}

// PurchaseResponse represents a unified purchase validation response
type PurchaseResponse struct {
	Success       bool   `json:"success"`
	Credits       int    `json:"credits"`
	TransactionID string `json:"transaction_id"`
	Message       string `json:"message"`
	Platform      string `json:"platform"`
}

// ValidatePurchase validates a purchase from either WeChat Pay or Apple IAP
// @Summary Validate purchase
// @Description Validates a purchase receipt from WeChat Pay or Apple IAP and awards credits
// @Tags billing
// @Accept json
// @Produce json
// @Param request body PurchaseRequest true "Purchase validation request"
// @Success 200 {object} model.SuccessResponse{data=PurchaseResponse}
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/v1/billing/purchase [post]
func (h *BillingHandler) ValidatePurchase(c *gin.Context) {
	var req PurchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("INVALID_REQUEST", err.Error()))
		return
	}

	// Get user from context (set by auth middleware)
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.NewErrorResponse("UNAUTHORIZED", "User not found in context"))
		return
	}
	
	userID, ok := userIDInterface.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("SERVER_ERROR", "Invalid user ID format"))
		return
	}

	// Normalize platform string
	platform := strings.ToLower(req.Platform)
	
	var response *PurchaseResponse
	var err error

	switch platform {
	case "wechat":
		response, err = h.validateWechatPurchase(c, userID, &req)
	case "apple":
		response, err = h.validateApplePurchase(c, userID, &req)
	default:
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("INVALID_PLATFORM", "Platform must be 'wechat' or 'apple'"))
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("VALIDATION_FAILED", err.Error()))
		return
	}

	if response.Success {
		c.JSON(http.StatusOK, model.NewSuccessResponse(response))
	} else {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("VALIDATION_FAILED", response.Message))
	}
}

// validateWechatPurchase validates a WeChat Pay purchase
func (h *BillingHandler) validateWechatPurchase(c *gin.Context, userID int64, req *PurchaseRequest) (*PurchaseResponse, error) {
	// For WeChat Pay, the receipt is typically the order ID
	validationReq := &service.PaymentValidationRequest{
		OrderID:       req.Receipt,
		TransactionID: "", // Will be filled from query
		ResultCode:    "SUCCESS",
	}

	validationResp, err := h.wechatPayService.ValidatePayment(c.Request.Context(), validationReq)
	if err != nil {
		return &PurchaseResponse{
			Success:  false,
			Message:  err.Error(),
			Platform: "wechat",
		}, nil
	}

	return &PurchaseResponse{
		Success:       validationResp.Success,
		Credits:       validationResp.Credits,
		TransactionID: req.Receipt,
		Message:       validationResp.Message,
		Platform:      "wechat",
	}, nil
}

// validateApplePurchase validates an Apple IAP purchase
func (h *BillingHandler) validateApplePurchase(c *gin.Context, userID int64, req *PurchaseRequest) (*PurchaseResponse, error) {
	// For Apple IAP, extract transaction ID from receipt if needed
	// In a real implementation, you'd parse the receipt to get transaction ID
	transactionID := "apple_" + req.Receipt[:min(8, len(req.Receipt))]

	processReq := &service.PurchaseProcessRequest{
		UserID:        userID,
		ProductID:     req.ProductID,
		TransactionID: transactionID,
		ReceiptData:   req.Receipt,
		Description:   req.Description,
	}

	processResp, err := h.appleIAPService.ProcessPurchase(c.Request.Context(), processReq)
	if err != nil {
		return &PurchaseResponse{
			Success:  false,
			Message:  err.Error(),
			Platform: "apple",
		}, nil
	}

	return &PurchaseResponse{
		Success:       processResp.Success,
		Credits:       processResp.Credits,
		TransactionID: processResp.TransactionID,
		Message:       processResp.Message,
		Platform:      "apple",
	}, nil
}

// GetUserCredits returns the user's current credit balance
// @Summary Get user credits
// @Description Returns the current credit balance for the authenticated user
// @Tags billing
// @Produce json
// @Success 200 {object} model.SuccessResponse{data=map[string]int}
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/v1/billing/credits [get]
func (h *BillingHandler) GetUserCredits(c *gin.Context) {
	// Get user from context (set by auth middleware)
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.NewErrorResponse("UNAUTHORIZED", "User not found in context"))
		return
	}
	
	userID, ok := userIDInterface.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("SERVER_ERROR", "Invalid user ID format"))
		return
	}

	credits, err := h.creditService.GetBalance(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("CREDITS_FETCH_FAILED", err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(map[string]int{
		"credits": credits,
	}))
}

// GetUserTransactions returns the user's transaction history
// @Summary Get user transactions
// @Description Returns paginated transaction history for the authenticated user
// @Tags billing
// @Produce json
// @Param limit query int false "Number of transactions to return (default: 20)"
// @Param offset query int false "Number of transactions to skip (default: 0)"
// @Success 200 {object} model.SuccessResponse{data=map[string]interface{}}
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/v1/me/transactions [get]
func (h *BillingHandler) GetUserTransactions(c *gin.Context) {
	// Get user from context (set by auth middleware)
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.NewErrorResponse("UNAUTHORIZED", "User not found in context"))
		return
	}
	
	userID, ok := userIDInterface.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("SERVER_ERROR", "Invalid user ID format"))
		return
	}

	// Parse pagination parameters
	limit := 20 // default
	offset := 0 // default
	
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit := parseInt(limitStr, limit); parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}
	
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsedOffset := parseInt(offsetStr, offset); parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// Get transactions
	transactions, err := h.transactionService.GetTransactionsByUserID(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("TRANSACTIONS_FETCH_FAILED", err.Error()))
		return
	}

	// Get total count for pagination
	totalCount, err := h.transactionService.CountTransactionsByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("TRANSACTION_COUNT_FAILED", err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(map[string]interface{}{
		"transactions": transactions,
		"pagination": map[string]interface{}{
			"limit":       limit,
			"offset":      offset,
			"total_count": totalCount,
			"has_more":    offset+limit < totalCount,
		},
	}))
}

// RegisterRoutes registers billing routes
func (h *BillingHandler) RegisterRoutes(r *gin.RouterGroup) {
	billing := r.Group("/billing")
	{
		billing.POST("/purchase", h.ValidatePurchase)
		billing.GET("/credits", h.GetUserCredits)
	}
}

// RegisterUserRoutes registers user-specific billing routes
func (h *BillingHandler) RegisterUserRoutes(r *gin.RouterGroup) {
	me := r.Group("/me")
	{
		me.GET("/transactions", h.GetUserTransactions)
	}
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Helper function to parse int with default
func parseInt(s string, defaultVal int) int {
	if val, err := strconv.Atoi(s); err == nil {
		return val
	}
	return defaultVal
} 