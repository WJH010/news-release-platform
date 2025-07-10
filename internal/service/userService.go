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
)

// UserService 用户服务接口
type UserService interface {
	Login(ctx context.Context, code string) (*model.User, error)
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
func (s *UserServiceImpl) Login(ctx context.Context, code string) (*model.User, error) {
	// 调用微信接口换取 openid 和 session_key
	openid, err := s.getOpenIDFromWechat(code)
	if err != nil {
		return nil, err
	}

	// 检查数据库中是否存在该 openid 对应的用户
	user, err := s.userRepo.GetUserByOpenID(ctx, openid)
	if err != nil {
		return nil, err
	}

	if user == nil {
		// 新用户，创建用户记录
		newUser := &model.User{
			OpenID: openid,
		}
		err := s.userRepo.CreateUser(ctx, newUser)
		if err != nil {
			return nil, err
		}
		user = newUser
	}

	return user, nil
}

// getOpenIDFromWechat 调用微信接口换取 openid
func (s *UserServiceImpl) getOpenIDFromWechat(code string) (string, error) {
	//url需要去自己申请
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", s.cfg.Wechat.AppID, s.cfg.Wechat.Secret, code)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 使用 io.ReadAll 替换 ioutil.ReadAll
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		OpenID     string `json:"openid"`
		SessionKey string `json:"session_key"`
		ErrCode    int    `json:"errcode"`
		ErrMsg     string `json:"errmsg"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	if result.ErrCode != 0 {
		return "", fmt.Errorf("微信接口返回错误: %d - %s", result.ErrCode, result.ErrMsg)
	}

	return result.OpenID, nil
}
