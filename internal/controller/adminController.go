package controller

import (
	"net/http"
	"news-release/internal/config"
	"news-release/internal/service"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type AdminController struct {
	adminService service.AdminUserService
	cfg          *config.Config
}

func NewAdminController(adminService service.AdminUserService, cfg *config.Config) *AdminController {
	return &AdminController{adminService: adminService, cfg: cfg}
}

// AdminLogin 管理系统登录接口
func (c *AdminController) AdminLogin(ctx *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	// 解析请求参数
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请输入账号和密码"})
		return
	}

	// 调用服务层登录
	admin, err := c.adminService.Login(ctx, req.Username, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 生成登录凭证（JWT）
	token, err := generateJWT(admin.ID, c.cfg.JWT.SecretKey, c.cfg.JWT.ExpirationHours)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"admin": gin.H{
			"id":       admin.ID,
			"username": admin.Username,
			"role":     admin.Role,
		},
		"token": token,
	})
}

// generateJWT 生成 JWT 令牌
func generateJWT(userID int, secretKey string, expirationHours int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * time.Duration(expirationHours)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
