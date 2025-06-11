package service

import (
	"context"
	"crypto/md5"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// MockWechatPayService implements WechatPayService for development/testing
type MockWechatPayService struct {
	creditService CreditService
	logger        *slog.Logger
	appID         string
	mchID         string
	apiKey        string
}

// NewMockWechatPayService creates a new mock WeChat Pay service
func NewMockWechatPayService(creditService CreditService, logger *slog.Logger) WechatPayService {
	return &MockWechatPayService{
		creditService: creditService,
		logger:        logger,
		appID:         "mock_app_id",
		mchID:         "mock_mch_id", 
		apiKey:        "mock_api_key",
	}
}

// CreatePayOrder creates a mock pre-pay order
func (s *MockWechatPayService) CreatePayOrder(ctx context.Context, req *PayOrderRequest) (*PayOrderResponse, error) {
	s.logger.InfoContext(ctx, "Creating mock pay order", 
		"user_id", req.UserID,
		"amount", req.Amount,
		"description", req.Description)

	// Generate mock order ID
	orderID := fmt.Sprintf("ORDER_%d_%s", time.Now().Unix(), uuid.New().String()[:8])
	
	// Generate mock prepay ID
	prepayID := fmt.Sprintf("wx%d%s", time.Now().Unix(), uuid.New().String()[:16])
	
	// Generate mock payment parameters
	nonceStr := uuid.New().String()[:32]
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	packageStr := fmt.Sprintf("prepay_id=%s", prepayID)
	
	// Generate mock sign (in real implementation, this would use WeChat's signing algorithm)
	sign := s.generateMockSign(packageStr, nonceStr, timeStamp)

	response := &PayOrderResponse{
		PrepayID:  prepayID,
		OrderID:   orderID,
		Package:   packageStr,
		NonceStr:  nonceStr,
		TimeStamp: timeStamp,
		Sign:      sign,
	}

	s.logger.InfoContext(ctx, "Mock pay order created", 
		"order_id", orderID,
		"prepay_id", prepayID)

	return response, nil
}

// ValidatePayment validates a mock payment result
func (s *MockWechatPayService) ValidatePayment(ctx context.Context, req *PaymentValidationRequest) (*PaymentValidationResponse, error) {
	s.logger.InfoContext(ctx, "Validating mock payment", 
		"order_id", req.OrderID,
		"result_code", req.ResultCode)

	// Mock validation - in development, always succeed
	if req.ResultCode == "SUCCESS" || req.ResultCode == "" {
		// For mock, assume 100 credits for 10 yuan (1000 cents)
		credits := 100
		
		// Extract user ID from order (in real implementation, this would be stored)
		// For mock, we'll hardcode a user ID
		userID := int64(1)
		
		// Add credits to user account
		description := fmt.Sprintf("WeChat Pay purchase - Order: %s", req.OrderID)
		err := s.creditService.AddCredits(ctx, userID, credits, description, &req.TransactionID)
		if err != nil {
			s.logger.ErrorContext(ctx, "Failed to add credits", "error", err)
			return &PaymentValidationResponse{
				Success: false,
				OrderID: req.OrderID,
				Message: "Failed to add credits",
			}, nil
		}

		return &PaymentValidationResponse{
			Success: true,
			OrderID: req.OrderID,
			Credits: credits,
			Message: "Payment successful",
		}, nil
	}

	return &PaymentValidationResponse{
		Success: false,
		OrderID: req.OrderID,
		Message: fmt.Sprintf("Payment failed: %s", req.ErrCodeDes),
	}, nil
}

// QueryPayment queries mock payment status
func (s *MockWechatPayService) QueryPayment(ctx context.Context, orderID string) (*PaymentQueryResponse, error) {
	s.logger.InfoContext(ctx, "Querying mock payment", "order_id", orderID)

	// Mock response - assume payment is successful
	return &PaymentQueryResponse{
		OrderID:       orderID,
		TransactionID: fmt.Sprintf("TX_%s", uuid.New().String()[:16]),
		TradeState:    "SUCCESS",
		Amount:        1000, // 10 yuan in cents
		PaidAt:        time.Now(),
	}, nil
}

// generateMockSign generates a mock signature for development
func (s *MockWechatPayService) generateMockSign(packageStr, nonceStr, timeStamp string) string {
	// This is a mock sign - real implementation would use WeChat's algorithm
	data := fmt.Sprintf("appId=%s&nonceStr=%s&package=%s&signType=MD5&timeStamp=%s&key=%s",
		s.appID, nonceStr, packageStr, timeStamp, s.apiKey)
	
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// Real WeChat Pay service implementation would go here
// type WechatPayServiceImpl struct {
//     appID      string
//     mchID      string
//     apiKey     string
//     httpClient *http.Client
//     logger     *slog.Logger
// }
//
// This would implement the actual WeChat Pay API calls:
// - Unified Order API for creating prepay orders
// - Payment notification handling
// - Order query API
// - Proper signature verification 