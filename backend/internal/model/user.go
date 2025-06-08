package model

import (
	"time"
)

// User represents a user account in the system
type User struct {
	ID           int64     `json:"id" db:"id"`
	WechatOpenID string    `json:"wechat_openid" db:"wechat_openid"`
	Nickname     string    `json:"nickname" db:"nickname"`
	AvatarURL    string    `json:"avatar_url" db:"avatar_url"`
	Credits      int       `json:"credits" db:"credits"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// UserCreateRequest represents the request to create a new user
type UserCreateRequest struct {
	WechatOpenID string `json:"wechat_openid" binding:"required"`
	Nickname     string `json:"nickname"`
	AvatarURL    string `json:"avatar_url"`
}

// UserUpdateRequest represents the request to update user information
type UserUpdateRequest struct {
	Nickname  *string `json:"nickname,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
} 