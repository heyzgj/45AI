package service

import (
	"context"
	"time"
)

// AppleIAPService defines the interface for Apple In-App Purchase operations
type AppleIAPService interface {
	// ValidateReceipt validates an App Store receipt
	ValidateReceipt(ctx context.Context, req *ReceiptValidationRequest) (*ReceiptValidationResponse, error)
	
	// ProcessPurchase processes a validated purchase and awards credits
	ProcessPurchase(ctx context.Context, req *PurchaseProcessRequest) (*PurchaseProcessResponse, error)
	
	// GetProductInfo retrieves product information for IAP products
	GetProductInfo(ctx context.Context, productIDs []string) (*ProductInfoResponse, error)
}

// ReceiptValidationRequest represents a request to validate an App Store receipt
type ReceiptValidationRequest struct {
	ReceiptData   string `json:"receipt_data" binding:"required"`
	UserID        int64  `json:"user_id" binding:"required"`
	ProductID     string `json:"product_id" binding:"required"`
	TransactionID string `json:"transaction_id"`
}

// ReceiptValidationResponse represents the response from receipt validation
type ReceiptValidationResponse struct {
	Valid         bool                     `json:"valid"`
	Receipt       *AppStoreReceipt        `json:"receipt,omitempty"`
	ErrorMessage  string                  `json:"error_message,omitempty"`
	Environment   string                  `json:"environment"` // "Sandbox" or "Production"
}

// PurchaseProcessRequest represents a request to process a validated purchase
type PurchaseProcessRequest struct {
	UserID        int64  `json:"user_id"`
	ProductID     string `json:"product_id"`
	TransactionID string `json:"transaction_id"`
	ReceiptData   string `json:"receipt_data"`
	Credits       int    `json:"credits"`
	Description   string `json:"description"`
}

// PurchaseProcessResponse represents the response from processing a purchase
type PurchaseProcessResponse struct {
	Success       bool   `json:"success"`
	Credits       int    `json:"credits"`
	TransactionID string `json:"transaction_id"`
	Message       string `json:"message"`
}

// ProductInfoResponse represents product information response
type ProductInfoResponse struct {
	Products []IAPProduct `json:"products"`
}

// AppStoreReceipt represents an App Store receipt structure
type AppStoreReceipt struct {
	Status             int                    `json:"status"`
	Environment        string                 `json:"environment"`
	Receipt            Receipt               `json:"receipt"`
	LatestReceiptInfo  []InAppPurchase       `json:"latest_receipt_info,omitempty"`
	PendingRenewalInfo []PendingRenewalInfo  `json:"pending_renewal_info,omitempty"`
}

// Receipt represents the main receipt data
type Receipt struct {
	ReceiptType         string          `json:"receipt_type"`
	BundleID            string          `json:"bundle_id"`
	ApplicationVersion  string          `json:"application_version"`
	DownloadID          int64           `json:"download_id"`
	VersionExternalID   int64           `json:"version_external_identifier"`
	ReceiptCreationDate string          `json:"receipt_creation_date"`
	ReceiptCreationDateMS string        `json:"receipt_creation_date_ms"`
	OriginalPurchaseDate string         `json:"original_purchase_date"`
	OriginalPurchaseDateMS string       `json:"original_purchase_date_ms"`
	RequestDate         string          `json:"request_date"`
	RequestDateMS       string          `json:"request_date_ms"`
	InApp               []InAppPurchase `json:"in_app"`
}

// InAppPurchase represents an in-app purchase transaction
type InAppPurchase struct {
	Quantity                string `json:"quantity"`
	ProductID               string `json:"product_id"`
	TransactionID           string `json:"transaction_id"`
	OriginalTransactionID   string `json:"original_transaction_id"`
	PurchaseDate            string `json:"purchase_date"`
	PurchaseDateMS          string `json:"purchase_date_ms"`
	PurchaseDatePST         string `json:"purchase_date_pst"`
	OriginalPurchaseDate    string `json:"original_purchase_date"`
	OriginalPurchaseDateMS  string `json:"original_purchase_date_ms"`
	OriginalPurchaseDatePST string `json:"original_purchase_date_pst"`
	ExpiresDate             string `json:"expires_date,omitempty"`
	ExpiresDateMS           string `json:"expires_date_ms,omitempty"`
	ExpiresDatePST          string `json:"expires_date_pst,omitempty"`
	WebOrderLineItemID      string `json:"web_order_line_item_id,omitempty"`
	IsTrialPeriod           string `json:"is_trial_period,omitempty"`
	IsInIntroOfferPeriod    string `json:"is_in_intro_offer_period,omitempty"`
}

// PendingRenewalInfo represents pending renewal information
type PendingRenewalInfo struct {
	ExpirationIntent   string `json:"expiration_intent,omitempty"`
	AutoRenewProductID string `json:"auto_renew_product_id,omitempty"`
	RetryFlag          string `json:"retry_flag,omitempty"`
	AutoRenewStatus    string `json:"auto_renew_status,omitempty"`
}

// IAPProduct represents an in-app purchase product
type IAPProduct struct {
	ProductID     string    `json:"product_id"`
	Credits       int       `json:"credits"`
	PriceUSD      float64   `json:"price_usd"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Popular       bool      `json:"popular"`
	CreatedAt     time.Time `json:"created_at"`
} 