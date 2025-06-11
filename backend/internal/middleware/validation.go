package middleware

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/45ai/backend/internal/model"
)

// ValidationMiddleware creates a new validator instance with custom rules
type ValidationMiddleware struct {
	validator *validator.Validate
}

// NewValidationMiddleware creates a new validation middleware instance
func NewValidationMiddleware() *ValidationMiddleware {
	v := validator.New()
	
	// Register custom field name function to use JSON tags
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	
	// Register custom validators
	registerCustomValidators(v)
	
	return &ValidationMiddleware{
		validator: v,
	}
}

// ValidateJSON middleware validates JSON request body against struct
func (vm *ValidationMiddleware) ValidateJSON() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// This middleware only validates if there's a JSON body
		if c.Request.Header.Get("Content-Type") != "application/json" {
			c.Next()
			return
		}
		
		// Continue to next handler - validation happens in individual handlers
		// using ShouldBindJSON which automatically triggers validation
		c.Next()
	})
}

// ValidateStruct validates a struct and returns formatted errors
func (vm *ValidationMiddleware) ValidateStruct(obj interface{}) *ValidationErrors {
	err := vm.validator.Struct(obj)
	if err == nil {
		return nil
	}
	
	var validationErrors []ValidationError
	
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrs {
			validationErrors = append(validationErrors, ValidationError{
				Field:   fieldError.Field(),
				Tag:     fieldError.Tag(),
				Message: getErrorMessage(fieldError),
				Value:   fieldError.Value(),
			})
		}
	}
	
	return &ValidationErrors{
		Errors: validationErrors,
	}
}

// ValidationError represents a single validation error
type ValidationError struct {
	Field   string      `json:"field"`
	Tag     string      `json:"tag"`
	Message string      `json:"message"`
	Value   interface{} `json:"value,omitempty"`
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

// Error implements the error interface
func (ve *ValidationErrors) Error() string {
	var messages []string
	for _, err := range ve.Errors {
		messages = append(messages, err.Message)
	}
	return strings.Join(messages, "; ")
}

// HandleValidationErrors middleware handles validation errors from Gin's ShouldBindJSON
func (vm *ValidationMiddleware) HandleValidationErrors() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Next()
		
		// Check if there were any errors during processing
		if len(c.Errors) > 0 {
			lastError := c.Errors.Last()
			
			// Check if it's a validation error
			if validationErrs, ok := lastError.Err.(validator.ValidationErrors); ok {
				var errors []ValidationError
				
				for _, fieldError := range validationErrs {
					errors = append(errors, ValidationError{
						Field:   fieldError.Field(),
						Tag:     fieldError.Tag(),
						Message: getErrorMessage(fieldError),
						Value:   fieldError.Value(),
					})
				}
				
				c.JSON(http.StatusBadRequest, model.NewErrorResponse("VALIDATION_FAILED", "Request validation errors"))
				c.Abort()
				return
			}
		}
	})
}

// registerCustomValidators registers custom validation rules
func registerCustomValidators(v *validator.Validate) {
	// Phone number validation (simple Chinese mobile number)
	v.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
		phone := fl.Field().String()
		if len(phone) != 11 {
			return false
		}
		// Chinese mobile numbers start with 1 and second digit is 3,4,5,6,7,8,9
		if phone[0] != '1' {
			return false
		}
		secondDigit := phone[1]
		return secondDigit >= '3' && secondDigit <= '9'
	})
	
	// WeChat OpenID validation
	v.RegisterValidation("openid", func(fl validator.FieldLevel) bool {
		openid := fl.Field().String()
		return len(openid) >= 20 && len(openid) <= 32
	})
	
	// Credit amount validation (positive integer)
	v.RegisterValidation("credits", func(fl validator.FieldLevel) bool {
		credits := fl.Field().Int()
		return credits > 0 && credits <= 10000 // Max 10,000 credits per transaction
	})
	
	// Template ID validation
	v.RegisterValidation("template_id", func(fl validator.FieldLevel) bool {
		templateID := fl.Field().Int()
		return templateID > 0
	})
}

// getErrorMessage returns a user-friendly error message for validation errors
func getErrorMessage(fe validator.FieldError) string {
	field := fe.Field()
	tag := fe.Tag()
	param := fe.Param()
	
	switch tag {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email address"
	case "min":
		if fe.Kind() == reflect.String {
			return field + " must be at least " + param + " characters long"
		}
		return field + " must be at least " + param
	case "max":
		if fe.Kind() == reflect.String {
			return field + " must be at most " + param + " characters long"
		}
		return field + " must be at most " + param
	case "len":
		return field + " must be exactly " + param + " characters long"
	case "phone":
		return field + " must be a valid Chinese mobile phone number"
	case "openid":
		return field + " must be a valid WeChat OpenID"
	case "credits":
		return field + " must be a positive number between 1 and 10,000"
	case "template_id":
		return field + " must be a valid template ID"
	case "oneof":
		return field + " must be one of: " + param
	case "url":
		return field + " must be a valid URL"
	case "uuid":
		return field + " must be a valid UUID"
	case "numeric":
		return field + " must be a number"
	case "alpha":
		return field + " must contain only letters"
	case "alphanum":
		return field + " must contain only letters and numbers"
	default:
		return field + " is invalid"
	}
}

// Common validation functions that can be used in handlers

// ValidateAndBind validates and binds JSON request to struct
func (vm *ValidationMiddleware) ValidateAndBind(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			var errors []ValidationError
			
			for _, fieldError := range validationErrs {
				errors = append(errors, ValidationError{
					Field:   fieldError.Field(),
					Tag:     fieldError.Tag(),
					Message: getErrorMessage(fieldError),
					Value:   fieldError.Value(),
				})
			}
			
			c.JSON(http.StatusBadRequest, model.NewErrorResponse("VALIDATION_FAILED", "Request validation errors"))
			return err
		}
		
		// Non-validation error (JSON parsing error, etc.)
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("INVALID_REQUEST", "Failed to parse request body"))
		return err
	}
	
	return nil
}

// Example usage in handlers:
/*
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := h.validator.ValidateAndBind(c, &req); err != nil {
		return // Error response already sent
	}
	
	// Process valid request...
}

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required,phone"`
	OpenID   string `json:"openid" validate:"required,openid"`
	Nickname string `json:"nickname" validate:"required,min=2,max=20"`
}
*/ 