package repository

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/45ai/backend/internal/config"
	"github.com/45ai/backend/internal/model"
)

type WechatRepository interface {
	Code2Session(code string) (*model.WechatLoginResponse, error)
}

type wechatRepositoryImpl struct {
	cfg config.WeChatConfig
}

func NewWechatRepository(cfg config.WeChatConfig) WechatRepository {
	return &wechatRepositoryImpl{cfg: cfg}
}

func (r *wechatRepositoryImpl) Code2Session(code string) (*model.WechatLoginResponse, error) {
	// Development testing: Allow test codes that start with "test_"
	// This simulates WeChat's API response for development testing
	if len(code) > 5 && code[:5] == "test_" {
		// Generate a unique OpenID based on the test code
		return &model.WechatLoginResponse{
			OpenID:  fmt.Sprintf("test_openid_%s", code[5:]), // Use code suffix as unique identifier
			UnionID: fmt.Sprintf("test_union_%s", code[5:]),
		}, nil
	}

	// Production: Use real WeChat API
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		r.cfg.AppID, r.cfg.AppSecret, code)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var wechatResp model.WechatLoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&wechatResp); err != nil {
		return nil, err
	}

	if wechatResp.OpenID == "" {
		return nil, fmt.Errorf("wechat api error: %v", wechatResp)
	}

	return &wechatResp, nil
} 