package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"news-release/internal/config"
	msgsvc "news-release/internal/message/service"
	"news-release/internal/user/dto"
	"news-release/internal/user/model"
	"news-release/internal/user/repository"
	"news-release/internal/utils"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/argon2"
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
	GetUserByID(ctx context.Context, userID int) (*dto.UserInfoResponse, error)
	ListAllUsers(ctx context.Context, page, pageSize int, req dto.ListUsersRequest) ([]*dto.ListUsersResponse, int64, error)
	// CreateAdminUser 新增管理员
	CreateAdminUser(ctx context.Context, req dto.CreateAdminRequest) error
	// BgLogin 后台登录
	BgLogin(ctx context.Context, req dto.BgLoginRequest) (string, error)
}

// UserServiceImpl 用户服务实现
type UserServiceImpl struct {
	userRepo repository.UserRepository
	msgSvc   msgsvc.MsgGroupService
	cfg      *config.Config
}

// Argon2参数配置
const (
	// 内存成本：哈希过程中使用的内存量（字节）
	argonMemory uint32 = 65536 // 64MB
	// 时间成本：计算迭代次数
	argonTime uint32 = 3
	// 并行度：使用的CPU核心数
	argonThreads uint8 = 4
	// 生成的哈希长度（字节）
	argonKeyLen uint32 = 32
	// 盐值长度（字节）
	argonSaltLen uint32 = 16
)

// NewUserService 创建用户服务实例
func NewUserService(userRepo repository.UserRepository, msgSvc msgsvc.MsgGroupService, cfg *config.Config) UserService {
	return &UserServiceImpl{userRepo: userRepo, msgSvc: msgSvc, cfg: cfg}
}

// Login 微信登录逻辑
func (svc *UserServiceImpl) Login(ctx context.Context, code string) (string, error) {
	// 调用微信接口
	wxResp, err := svc.getFromWechat(code)
	if err != nil {
		return "", err
	}

	// 查找或创建用户
	userID, userRole, err := svc.findOrCreateUser(ctx, wxResp.OpenID, wxResp.SessionKey, wxResp.UnionID)
	if err != nil {
		return "", fmt.Errorf("处理用户信息失败: %v", err)
	}

	// 生成登录状态 Token
	token, err := svc.generateToken(wxResp.OpenID, userID, userRole)
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
func (svc *UserServiceImpl) findOrCreateUser(ctx context.Context, openID, sessionKey, unionID string) (int, int, error) {
	// 查找用户
	user, err := svc.userRepo.GetUserByOpenID(ctx, openID)
	if err != nil {
		return 0, 0, err
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
			return user.UserID, user.Role, err
		}
		// 新用户创建成功后加入全体成员的群组
		svc.msgSvc.AddUserToAllUserGroups(ctx, user.UserID)
	} else {
		// 如果用户存在，更新session_key和登录时间
		if err := svc.userRepo.UpdateSessionAndLoginTime(ctx, user.UserID, sessionKey); err != nil {
			return 0, 0, err
		}
	}

	return user.UserID, user.Role, nil
}

// 生成JWT Token
func (svc *UserServiceImpl) generateToken(openID string, userID int, userRole int) (string, error) {
	// 创建令牌声明
	claims := jwt.MapClaims{
		"openid":    openID,
		"userid":    userID,
		"user_role": userRole,
		"exp":       time.Now().Add(time.Hour * 24).Unix(), // 令牌有效期24小时
		"iat":       time.Now().Unix(),
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

func (svc *UserServiceImpl) GetUserByID(ctx context.Context, userID int) (*dto.UserInfoResponse, error) {
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

// ListAllUsers 分页查询用户列表
func (svc *UserServiceImpl) ListAllUsers(ctx context.Context, page, pageSize int, req dto.ListUsersRequest) ([]*dto.ListUsersResponse, int64, error) {
	return svc.userRepo.ListAllUsers(ctx, page, pageSize, req)
}

// CreateAdminUser 新增管理员
func (svc *UserServiceImpl) CreateAdminUser(ctx context.Context, req dto.CreateAdminRequest) error {
	var avatar string
	if req.AvatarURL == "" {
		avatar = "http://47.113.194.28:9000/news-platform/images/202508/1754126743005963551.webp"
	} else {
		avatar = req.AvatarURL
	}

	// 对密码进行哈希处理
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return utils.NewSystemError(fmt.Errorf("密码加密失败: %w", err))
	}

	// 创建数据
	user := &model.User{
		Nickname:      req.Nickname,
		Name:          req.Name,
		AvatarURL:     avatar,
		PhoneNumber:   req.PhoneNumber,
		Email:         req.Email,
		Role:          req.Role,
		Password:      hashedPassword,
		LastLoginTime: time.Now(),
	}

	if err := svc.userRepo.Create(ctx, user); err != nil {
		return err
	}
	return nil
}

// 生成密码哈希
func hashPassword(password string) (string, error) {
	// 生成随机盐值
	salt := make([]byte, argonSaltLen)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	// 使用Argon2id变体进行哈希（推荐用于密码哈希）
	hash := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)

	// 组合盐值和哈希值，并进行Base64编码以便存储
	// 格式: $argon2id$v=19$m=65536,t=3,p=4$<salt>$<hash>
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// 包含算法参数以便验证时使用
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, argonMemory, argonTime, argonThreads, b64Salt, b64Hash)

	return encodedHash, nil
}

// BgLogin 后台登录
func (svc *UserServiceImpl) BgLogin(ctx context.Context, req dto.BgLoginRequest) (string, error) {
	// 从数据库中根据手机号查询密码
	password, err := svc.userRepo.GetPasswordByPhone(ctx, req.PhoneNumber)
	if err != nil {
		return "", err
	}

	// 验证密码
	ok, err := verifyPassword(password, req.Password)
	if err != nil {
		return "", utils.NewSystemError(fmt.Errorf("验证密码失败: %w", err))
	}
	if !ok {
		return "", utils.NewBusinessError(utils.ErrCodeAuthFailed, "密码错误")
	}

	// 密码验证成功，生成JWT Token
	token, err := svc.generateToken(req.PhoneNumber, 0, utils.RoleAdmin)
	if err != nil {
		return "", utils.NewSystemError(fmt.Errorf("生成Token失败: %w", err))
	}

	return token, nil
}

// 验证密码
func verifyPassword(encodedHash, password string) (bool, error) {
	// 解析格式: $argon2id$v=19$m=65536,t=3,p=4$<salt>$<hash>
	// 按 $ 分割字符串，得到各部分
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, fmt.Errorf("哈希格式错误")
	}
	// 验证算法是否为 argon2id
	if parts[1] != "argon2id" {
		return false, fmt.Errorf("不支持的算法: %s", parts[1])
	}

	// 解析版本号（如 v=19）
	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return false, fmt.Errorf("解析版本失败: %v", err)
	}

	// 解析参数（m=内存, t=时间, p=并行度）
	var memory, time uint32
	var threads uint8
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads); err != nil {
		return false, fmt.Errorf("解析参数失败: %v", err)
	}

	// 提取盐值和哈希值（直接从分割结果中获取，避免解析错误）
	saltStr := parts[4]
	hashStr := parts[5]

	// 解码盐值和哈希值
	saltBytes, err := base64.RawStdEncoding.DecodeString(saltStr)
	if err != nil {
		return false, err
	}

	hashBytes, err := base64.RawStdEncoding.DecodeString(hashStr)
	if err != nil {
		return false, err
	}

	// 使用相同的参数计算输入密码的哈希
	inputHash := argon2.IDKey([]byte(password), saltBytes, time, memory, threads, uint32(len(hashBytes)))

	// 比较计算出的哈希和存储的哈希
	return constantTimeCompare(inputHash, hashBytes), nil
}

// 常量时间比较函数，防止时序攻击
func constantTimeCompare(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	result := 0
	for i := range a {
		result |= int(a[i] ^ b[i])
	}
	return result == 0
}
