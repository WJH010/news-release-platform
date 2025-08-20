package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"news-release/internal/config"
	"news-release/internal/user/dto"
	"news-release/internal/user/model"
	"news-release/internal/user/repository"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// WxLoginResponse 微信登录请求参数
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
	UpdateUserInfo(ctx context.Context, userID int, req dto.UserUpdateRequest) error
	GetUserByID(ctx context.Context, userID int) (*model.User, error)
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
func (svc *UserServiceImpl) Login(ctx context.Context, code string) (string, error) {
	// 调用微信接口
	wxResp, err := svc.getFromWechat(code)
	if err != nil {
		return "", err
	}

	// 查找或创建用户
	userID, err := svc.findOrCreateUser(ctx, wxResp.OpenID, wxResp.SessionKey, wxResp.UnionID)
	if err != nil {
		return "", fmt.Errorf("处理用户信息失败: %v", err)
	}

	// 生成登录状态 Token
	token, err := svc.generateToken(wxResp.OpenID, userID)
	if err != nil {
		return "", fmt.Errorf("生成Token失败: %v", err)
	}

	return token, nil
}

// getFromWechat 调用微信接口
func (svc *UserServiceImpl) getFromWechat(code string) (WxLoginResponse, error) {
	var wxResp WxLoginResponse

	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", svc.cfg.Wechat.AppID, svc.cfg.Wechat.AppSecret, code)

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
func (svc *UserServiceImpl) findOrCreateUser(ctx context.Context, openID, sessionKey, unionID string) (int, error) {
	// 查找用户
	user, err := svc.userRepo.GetUserByOpenID(ctx, openID)
	if err != nil {
		return 0, err
	}

	now := time.Now()

	// 如果用户不存在，创建新用户
	if user == nil {
		// 生成默认昵称和头像
		defaultNickname := fmt.Sprintf("微信用户%s", openID[len(openID)-5:]) // 拼接OpenID的后5位作为默认昵称
		defaultAvatar := "http://47.113.194.28:9000/news-platform/images/202508/1754126743005963551.webp"
		user = &model.User{
			OpenID:        openID,
			SessionKey:    sessionKey,
			UnionID:       unionID,
			Nickname:      defaultNickname,
			AvatarURL:     defaultAvatar,
			LastLoginTime: now,
		}

		if err := svc.userRepo.Create(ctx, user); err != nil {
			return user.UserID, err
		}
	} else {
		// 如果用户存在，更新session_key和登录时间
		if err := svc.userRepo.UpdateSessionAndLoginTime(ctx, user.UserID, sessionKey); err != nil {
			return 0, err
		}
	}

	return user.UserID, nil
}

// 生成JWT Token
func (svc *UserServiceImpl) generateToken(openID string, userID int) (string, error) {
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
	return token.SignedString([]byte(svc.cfg.JWT.JwtSecret))
}

// UpdateUserInfo 更新用户信息
func (svc *UserServiceImpl) UpdateUserInfo(ctx context.Context, userID int, req dto.UserUpdateRequest) error {
	// 查询用户是否存在
	user, err := svc.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("查询用户失败: %v", err)
	}
	if user == nil {
		return fmt.Errorf("用户不存在")
	}

	// 构建更新字段映射
	updateFields := make(map[string]interface{})
	if req.Nickname != nil {
		updateFields["nickname"] = *req.Nickname
	}
	if req.AvatarURL != nil {
		updateFields["avatar_url"] = *req.AvatarURL
	}
	if req.Name != nil {
		updateFields["name"] = *req.Name
	}
	if req.Gender != nil {
		updateFields["gender"] = *req.Gender
	}
	if req.PhoneNumber != nil {
		updateFields["phone_number"] = *req.PhoneNumber
	}
	if req.Email != nil {
		updateFields["email"] = *req.Email
	}
	if req.Unit != nil {
		updateFields["unit"] = *req.Unit
	}
	if req.Department != nil {
		updateFields["department"] = *req.Department
	}
	if req.Position != nil {
		updateFields["position"] = *req.Position
	}
	if req.Industry != nil {
		updateFields["industry"] = *req.Industry
	}

	// 执行更新
	if len(updateFields) > 0 {
		if err := svc.userRepo.Update(ctx, userID, updateFields); err != nil {
			return fmt.Errorf("更新用户信息失败: %v", err)
		}
	}

	return nil
}

func (svc *UserServiceImpl) GetUserByID(ctx context.Context, userID int) (*model.User, error) {
	// 查询用户信息
	user, err := svc.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}
	if user == nil {
		return nil, fmt.Errorf("用户不存在")
	}

	return user, nil
}
