package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"time"
)

// MockAppleIAPService implements AppleIAPService for development/testing
type MockAppleIAPService struct {
	creditService CreditService
	logger        *slog.Logger
	bundleID      string
	environment   string
}

// NewMockAppleIAPService creates a new mock Apple IAP service
func NewMockAppleIAPService(creditService CreditService, logger *slog.Logger) AppleIAPService {
	return &MockAppleIAPService{
		creditService: creditService,
		logger:        logger,
		bundleID:      "com.45ai.app",
		environment:   "Sandbox",
	}
}

// ValidateReceipt validates a mock App Store receipt
func (s *MockAppleIAPService) ValidateReceipt(ctx context.Context, req *ReceiptValidationRequest) (*ReceiptValidationResponse, error) {
	s.logger.InfoContext(ctx, "Validating mock App Store receipt",
		"user_id", req.UserID,
		"product_id", req.ProductID,
		"transaction_id", req.TransactionID)

	// Decode receipt data (in real implementation, this would be sent to Apple's servers)
	_, err := base64.StdEncoding.DecodeString(req.ReceiptData)
	if err != nil {
		return &ReceiptValidationResponse{
			Valid:        false,
			ErrorMessage: "Invalid receipt data format",
		}, nil
	}

	// Mock successful validation
	now := time.Now()
	mockReceipt := &AppStoreReceipt{
		Status:      0, // 0 = valid receipt
		Environment: s.environment,
		Receipt: Receipt{
			ReceiptType:            "ProductionSandbox",
			BundleID:               s.bundleID,
			ApplicationVersion:     "1.0.0",
			DownloadID:             123456789,
			VersionExternalID:      987654321,
			ReceiptCreationDate:    now.Format(time.RFC3339),
			ReceiptCreationDateMS:  fmt.Sprintf("%d", now.UnixMilli()),
			OriginalPurchaseDate:   now.Format(time.RFC3339),
			OriginalPurchaseDateMS: fmt.Sprintf("%d", now.UnixMilli()),
			RequestDate:            now.Format(time.RFC3339),
			RequestDateMS:          fmt.Sprintf("%d", now.UnixMilli()),
			InApp: []InAppPurchase{
				{
					Quantity:                "1",
					ProductID:               req.ProductID,
					TransactionID:           req.TransactionID,
					OriginalTransactionID:   req.TransactionID,
					PurchaseDate:            now.Format(time.RFC3339),
					PurchaseDateMS:          fmt.Sprintf("%d", now.UnixMilli()),
					PurchaseDatePST:         now.Format("2006-01-02 15:04:05 America/Los_Angeles"),
					OriginalPurchaseDate:    now.Format(time.RFC3339),
					OriginalPurchaseDateMS:  fmt.Sprintf("%d", now.UnixMilli()),
					OriginalPurchaseDatePST: now.Format("2006-01-02 15:04:05 America/Los_Angeles"),
				},
			},
		},
	}

	response := &ReceiptValidationResponse{
		Valid:       true,
		Receipt:     mockReceipt,
		Environment: s.environment,
	}

	s.logger.InfoContext(ctx, "Mock receipt validation successful",
		"product_id", req.ProductID,
		"transaction_id", req.TransactionID)

	return response, nil
}

// ProcessPurchase processes a validated purchase and awards credits
func (s *MockAppleIAPService) ProcessPurchase(ctx context.Context, req *PurchaseProcessRequest) (*PurchaseProcessResponse, error) {
	s.logger.InfoContext(ctx, "Processing mock IAP purchase",
		"user_id", req.UserID,
		"product_id", req.ProductID,
		"credits", req.Credits)

	// First validate the receipt
	validationReq := &ReceiptValidationRequest{
		ReceiptData:   req.ReceiptData,
		UserID:        req.UserID,
		ProductID:     req.ProductID,
		TransactionID: req.TransactionID,
	}

	validationResp, err := s.ValidateReceipt(ctx, validationReq)
	if err != nil {
		return &PurchaseProcessResponse{
			Success: false,
			Message: fmt.Sprintf("Receipt validation failed: %v", err),
		}, nil
	}

	if !validationResp.Valid {
		return &PurchaseProcessResponse{
			Success: false,
			Message: fmt.Sprintf("Invalid receipt: %s", validationResp.ErrorMessage),
		}, nil
	}

	// Determine credits based on product ID
	credits := s.getCreditsForProduct(req.ProductID)
	if credits == 0 {
		return &PurchaseProcessResponse{
			Success: false,
			Message: "Unknown product ID",
		}, nil
	}

	// Add credits to user account
	description := fmt.Sprintf("Apple IAP purchase - Product: %s", req.ProductID)
	err = s.creditService.AddCredits(ctx, req.UserID, credits, description, &req.TransactionID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to add credits", "error", err)
		return &PurchaseProcessResponse{
			Success: false,
			Message: "Failed to add credits",
		}, nil
	}

	return &PurchaseProcessResponse{
		Success:       true,
		Credits:       credits,
		TransactionID: req.TransactionID,
		Message:       "Purchase processed successfully",
	}, nil
}

// GetProductInfo retrieves product information for IAP products
func (s *MockAppleIAPService) GetProductInfo(ctx context.Context, productIDs []string) (*ProductInfoResponse, error) {
	s.logger.InfoContext(ctx, "Getting mock product info", "product_ids", productIDs)

	products := []IAPProduct{
		{
			ProductID:   "com.45ai.credits.100",
			Credits:     100,
			PriceUSD:    0.99,
			Title:       "100胶卷",
			Description: "基础套餐",
			Popular:     false,
			CreatedAt:   time.Now(),
		},
		{
			ProductID:   "com.45ai.credits.300",
			Credits:     300,
			PriceUSD:    2.99,
			Title:       "300胶卷",
			Description: "热门套餐",
			Popular:     true,
			CreatedAt:   time.Now(),
		},
		{
			ProductID:   "com.45ai.credits.500",
			Credits:     500,
			PriceUSD:    4.99,
			Title:       "500胶卷",
			Description: "超值套餐",
			Popular:     false,
			CreatedAt:   time.Now(),
		},
		{
			ProductID:   "com.45ai.credits.1000",
			Credits:     1000,
			PriceUSD:    9.99,
			Title:       "1000胶卷",
			Description: "豪华套餐",
			Popular:     false,
			CreatedAt:   time.Now(),
		},
	}

	// Filter products if specific IDs requested
	if len(productIDs) > 0 {
		var filteredProducts []IAPProduct
		for _, product := range products {
			for _, requestedID := range productIDs {
				if product.ProductID == requestedID {
					filteredProducts = append(filteredProducts, product)
					break
				}
			}
		}
		products = filteredProducts
	}

	return &ProductInfoResponse{
		Products: products,
	}, nil
}

// getCreditsForProduct returns the number of credits for a given product ID
func (s *MockAppleIAPService) getCreditsForProduct(productID string) int {
	switch productID {
	case "com.45ai.credits.100":
		return 100
	case "com.45ai.credits.300":
		return 300
	case "com.45ai.credits.500":
		return 500
	case "com.45ai.credits.1000":
		return 1000
	default:
		return 0
	}
}

// Real Apple IAP service implementation would go here
// type AppleIAPServiceImpl struct {
//     httpClient   *http.Client
//     logger       *slog.Logger
//     sandboxURL   string
//     productionURL string
//     sharedSecret string
// }
//
// This would implement the actual Apple App Store receipt validation:
// - Send receipt to Apple's verifyReceipt endpoint
// - Handle Apple's response format
// - Implement proper retry logic for network failures
// - Handle different receipt status codes
// - Implement duplicate transaction protection 