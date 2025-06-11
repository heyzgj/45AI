package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/45ai/backend/internal/service"
	"github.com/45ai/backend/internal/model"
)

// PaymentHandler handles payment-related HTTP requests
type PaymentHandler struct {
	wechatPayService service.WechatPayService
	appleIAPService  service.AppleIAPService
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(wechatPayService service.WechatPayService, appleIAPService service.AppleIAPService) *PaymentHandler {
	return &PaymentHandler{
		wechatPayService: wechatPayService,
		appleIAPService:  appleIAPService,
	}
}

// CreatePayOrder creates a new payment order
// @Summary Create payment order
// @Description Creates a WeChat Pay pre-order for credit purchase
// @Tags payment
// @Accept json
// @Produce json
// @Param request body service.PayOrderRequest true "Pay order request"
// @Success 200 {object} model.SuccessResponse{data=service.PayOrderResponse}
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/v1/payment/create-order [post]
func (h *PaymentHandler) CreatePayOrder(c *gin.Context) {
	var req service.PayOrderRequest
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
	
	req.UserID = userID

	// Get user's openid from context (set by auth middleware)
	openID := c.GetString("openid")
	if openID == "" {
		c.JSON(http.StatusUnauthorized, model.NewErrorResponse("UNAUTHORIZED", "OpenID not found in context"))
		return
	}
	
	req.OpenID = openID

	response, err := h.wechatPayService.CreatePayOrder(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("PAY_ORDER_FAILED", err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(response))
}

// ValidatePayment validates a payment result
// @Summary Validate payment
// @Description Validates a WeChat Pay payment result and adds credits
// @Tags payment
// @Accept json
// @Produce json
// @Param request body service.PaymentValidationRequest true "Payment validation request"
// @Success 200 {object} model.SuccessResponse{data=service.PaymentValidationResponse}
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/v1/payment/validate [post]
func (h *PaymentHandler) ValidatePayment(c *gin.Context) {
	var req service.PaymentValidationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("INVALID_REQUEST", err.Error()))
		return
	}

	response, err := h.wechatPayService.ValidatePayment(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("VALIDATION_FAILED", err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(response))
}

// QueryPayment queries payment status
// @Summary Query payment status
// @Description Queries the status of a payment order
// @Tags payment
// @Produce json
// @Param order_id path string true "Order ID"
// @Success 200 {object} model.SuccessResponse{data=service.PaymentQueryResponse}
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/v1/payment/query/{order_id} [get]
func (h *PaymentHandler) QueryPayment(c *gin.Context) {
	orderID := c.Param("order_id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("INVALID_REQUEST", "Order ID is required"))
		return
	}

	response, err := h.wechatPayService.QueryPayment(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("QUERY_FAILED", err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(response))
}

// ValidateAppleReceipt validates an App Store receipt
// @Summary Validate App Store receipt
// @Description Validates an Apple App Store receipt and processes the purchase
// @Tags payment
// @Accept json
// @Produce json
// @Param request body service.ReceiptValidationRequest true "Receipt validation request"
// @Success 200 {object} model.SuccessResponse{data=service.ReceiptValidationResponse}
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/v1/payment/apple/validate-receipt [post]
func (h *PaymentHandler) ValidateAppleReceipt(c *gin.Context) {
	var req service.ReceiptValidationRequest
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
	
	req.UserID = userID

	response, err := h.appleIAPService.ValidateReceipt(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("RECEIPT_VALIDATION_FAILED", err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(response))
}

// ProcessApplePurchase processes an Apple IAP purchase
// @Summary Process Apple IAP purchase
// @Description Processes a validated Apple IAP purchase and awards credits
// @Tags payment
// @Accept json
// @Produce json
// @Param request body service.PurchaseProcessRequest true "Purchase process request"
// @Success 200 {object} model.SuccessResponse{data=service.PurchaseProcessResponse}
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/v1/payment/apple/process-purchase [post]
func (h *PaymentHandler) ProcessApplePurchase(c *gin.Context) {
	var req service.PurchaseProcessRequest
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
	
	req.UserID = userID

	response, err := h.appleIAPService.ProcessPurchase(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("PURCHASE_PROCESSING_FAILED", err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(response))
}

// GetAppleProducts gets Apple IAP product information
// @Summary Get Apple IAP products
// @Description Retrieves information about available Apple IAP products
// @Tags payment
// @Produce json
// @Param product_ids query string false "Comma-separated product IDs"
// @Success 200 {object} model.SuccessResponse{data=service.ProductInfoResponse}
// @Failure 500 {object} model.ErrorResponse
// @Router /api/v1/payment/apple/products [get]
func (h *PaymentHandler) GetAppleProducts(c *gin.Context) {
	productIDsParam := c.Query("product_ids")
	var productIDs []string
	if productIDsParam != "" {
		// Split comma-separated product IDs
		productIDs = []string{productIDsParam} // Simplified for now
	}

	response, err := h.appleIAPService.GetProductInfo(c.Request.Context(), productIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("PRODUCT_INFO_FAILED", err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(response))
}

// SimplePurchase handles a simplified purchase request for development
// @Summary Simple purchase
// @Description Simplified purchase endpoint for development/testing
// @Tags payment
// @Accept json
// @Produce json
// @Param request body SimplePurchaseRequest true "Simple purchase request"
// @Success 200 {object} model.SuccessResponse{data=SimplePurchaseResponse}
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/v1/billing/purchase [post]
func (h *PaymentHandler) SimplePurchase(c *gin.Context) {
	var req SimplePurchaseRequest
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

	// For development, simulate successful purchase and add credits
	var credits int
	if req.Platform == "wechat" {
		// Map amount to credits (simplified)
		credits = int(req.Amount * 10) // 1 yuan = 10 credits
	} else if req.Platform == "apple" {
		// Map product ID to credits
		switch req.ProductID {
		case "com.45ai.credits.50":
			credits = 50
		case "com.45ai.credits.120":
			credits = 120
		case "com.45ai.credits.300":
			credits = 300
		case "com.45ai.credits.600":
			credits = 600
		default:
			credits = 100 // Default
		}
	}

	// Add credits to user account (mock transaction)
	// In a real implementation, this would be done after payment verification
	transactionID := fmt.Sprintf("MOCK_%d_%d", userID, time.Now().Unix())
	
	// For development, we'll simulate adding credits directly
	// In production, this would be done through the payment service after verification
	
	response := SimplePurchaseResponse{
		Success:       true,
		Credits:       credits,
		TransactionID: transactionID,
		Message:       "Purchase successful (development mode)",
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(response))
}

// SimplePurchaseRequest represents a simplified purchase request
type SimplePurchaseRequest struct {
	Amount      float64 `json:"amount,omitempty"`      // For WeChat Pay
	ProductID   string  `json:"product_id,omitempty"`  // For Apple IAP
	Title       string  `json:"title,omitempty"`
	Description string  `json:"description"`
	Platform    string  `json:"platform"` // "wechat" or "apple"
}

// SimplePurchaseResponse represents a simplified purchase response
type SimplePurchaseResponse struct {
	Success       bool   `json:"success"`
	Credits       int    `json:"credits"`
	TransactionID string `json:"transaction_id"`
	Message       string `json:"message"`
}

// RegisterRoutes registers payment routes
func (h *PaymentHandler) RegisterRoutes(r *gin.RouterGroup) {
	payment := r.Group("/payment")
	{
		// WeChat Pay routes
		payment.POST("/create-order", h.CreatePayOrder)
		payment.POST("/validate", h.ValidatePayment)
		payment.GET("/query/:order_id", h.QueryPayment)
		
		// Apple IAP routes
		apple := payment.Group("/apple")
		{
			apple.POST("/validate-receipt", h.ValidateAppleReceipt)
			apple.POST("/process-purchase", h.ProcessApplePurchase)
			apple.GET("/products", h.GetAppleProducts)
		}
	}
} 