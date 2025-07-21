package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"news-release/internal/config"
	"news-release/internal/model"
	"news-release/internal/repository"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type WxLoginResponse struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid,omitempty"`
	ErrCode    int    `json:"errcode,omitempty"`
	ErrMsg     string `json:"errmsg,omitempty"`
}

// UserService 用户服务接口
type UserService interface {
	Login(ctx context.Context, code string) (string, error)
}

// UserServiceImpl 用户服务实现
type UserServiceImpl struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo repository.UserRepository, cfg *config.Config) UserService {
	return &UserServiceImpl{userRepo: userRepo, cfg: cfg}
}

// Login 微信登录逻辑
func (s *UserServiceImpl) Login(ctx context.Context, code string) (string, error) {
	// 调用微信接口
	wxResp, err := s.getFromWechat(code)
	if err != nil {
		return "", err
	}

	// 查找或创建用户
	userID, err := s.findOrCreateUser(ctx, wxResp.OpenID, wxResp.SessionKey, wxResp.UnionID)
	if err != nil {
		return "", fmt.Errorf("处理用户信息失败: %v", err)
	}

	// 生成登录状态 Token
	token, err := s.generateToken(wxResp.OpenID, userID)
	if err != nil {
		return "", fmt.Errorf("生成Token失败: %v", err)
	}

	return token, nil
}

// getFromWechat 调用微信接口
func (s *UserServiceImpl) getFromWechat(code string) (WxLoginResponse, error) {
	var wxResp WxLoginResponse

	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", s.cfg.Wechat.AppID, s.cfg.Wechat.AppSecret, code)

	resp, err := http.Get(url)
	if err != nil {
		return wxResp, err
	}
	defer resp.Body.Close()

	// 读取微信响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return wxResp, fmt.Errorf("读取微信响应失败: %v", err)
	}

	// 解析微信响应
	err = json.Unmarshal(body, &wxResp)
	if err != nil {
		return wxResp, fmt.Errorf("解析微信响应失败: %v", err)
	}

	if wxResp.ErrCode != 0 {
		return wxResp, fmt.Errorf("微信登录错误: %d - %s", wxResp.ErrCode, wxResp.ErrMsg)
	}

	return wxResp, nil
}

// 查找或创建用户
func (s *UserServiceImpl) findOrCreateUser(ctx context.Context, openID, sessionKey, unionID string) (int, error) {
	// 查找用户
	user, err := s.userRepo.GetUserByOpenID(ctx, openID)
	if err != nil {
		return user.UserID, err
	}

	now := time.Now()

	// 如果用户不存在，创建新用户
	if user == nil {
		user = &model.User{
			OpenID:        openID,
			SessionKey:    sessionKey,
			UnionID:       unionID,
			LastLoginTime: now,
		}

		if err := s.userRepo.Create(ctx, user); err != nil {
			return user.UserID, err
		}
	} else {
		// 如果用户存在，更新session_key和登录时间
		user.SessionKey = sessionKey
		user.LastLoginTime = now

		if err := s.userRepo.Update(ctx, user); err != nil {
			return user.UserID, err
		}
	}

	return user.UserID, nil
}

// 生成JWT Token
func (s *UserServiceImpl) generateToken(openID string, userID int) (string, error) {
	// 创建令牌声明
	claims := jwt.MapClaims{
		"openid": openID,
		"userid": userID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // 令牌有效期24小时
		"iat":    time.Now().Unix(),
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	return token.SignedString([]byte(s.cfg.JWT.JwtSecret))
}
