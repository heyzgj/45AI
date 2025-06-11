package service

import (
	"context"
	"time"
)

// WechatPayService defines the interface for WeChat Pay operations
type WechatPayService interface {
	// CreatePayOrder creates a pre-pay order with WeChat Pay
	CreatePayOrder(ctx context.Context, req *PayOrderRequest) (*PayOrderResponse, error)
	
	// ValidatePayment validates a payment result from the frontend
	ValidatePayment(ctx context.Context, req *PaymentValidationRequest) (*PaymentValidationResponse, error)
	
	// QueryPayment queries the payment status from WeChat Pay
	QueryPayment(ctx context.Context, orderID string) (*PaymentQueryResponse, error)
}

// PayOrderRequest represents a request to create a pay order
type PayOrderRequest struct {
	UserID      int64  `json:"user_id" binding:"required"`
	Amount      int    `json:"amount" binding:"required,min=1"`        // Amount in cents
	Description string `json:"description" binding:"required"`
	OpenID      string `json:"openid" binding:"required"`
}

// PayOrderResponse represents the response from creating a pay order
type PayOrderResponse struct {
	PrepayID  string `json:"prepay_id"`
	OrderID   string `json:"order_id"`
	Package   string `json:"package"`
	NonceStr  string `json:"nonce_str"`
	TimeStamp string `json:"time_stamp"`
	Sign      string `json:"sign"`
}

// PaymentValidationRequest represents a request to validate payment
type PaymentValidationRequest struct {
	OrderID       string `json:"order_id" binding:"required"`
	TransactionID string `json:"transaction_id"`
	ResultCode    string `json:"result_code"`
	ErrCode       string `json:"err_code,omitempty"`
	ErrCodeDes    string `json:"err_code_des,omitempty"`
}

// PaymentValidationResponse represents the response from payment validation
type PaymentValidationResponse struct {
	Success   bool   `json:"success"`
	OrderID   string `json:"order_id"`
	Credits   int    `json:"credits"`
	Message   string `json:"message"`
}

// PaymentQueryResponse represents the response from payment query
type PaymentQueryResponse struct {
	OrderID       string    `json:"order_id"`
	TransactionID string    `json:"transaction_id"`
	TradeState    string    `json:"trade_state"`
	Amount        int       `json:"amount"`
	PaidAt        time.Time `json:"paid_at"`
} 