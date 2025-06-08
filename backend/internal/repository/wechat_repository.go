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