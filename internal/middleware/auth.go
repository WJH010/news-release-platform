package middleware

import (
	"errors"
	"net/http"
	"news-release/internal/config"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// 自定义Claims结构体，明确指定字段类型
type CustomClaims struct {
	OpenID             string `json:"openid"`
	UserID             int    `json:"userid"`
	jwt.StandardClaims        // 嵌入标准声明
}

// AuthMiddleware JWT认证中间件
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logrus.Warn("未提供认证信息")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未提供认证信息"})
			return
		}

		// 验证Authorization头格式
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logrus.Warn("认证格式错误")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "认证格式格式错误"})
			return
		}

		// 解析JWT令牌，使用自定义Claims
		claims, err := parseToken(cfg, parts[1])
		if err != nil {
			logrus.Warnf("解析JWT令牌失败: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "解析JWT令牌失败: " + err.Error()})
			return
		}

		// 将openID和userID存入上下文（此时已是正确类型）
		c.Set("openid", claims.OpenID)
		c.Set("userid", claims.UserID)

		c.Next()
	}
}

// 解析JWT令牌
func parseToken(cfg *config.Config, tokenString string) (*CustomClaims, error) {
	secret := []byte(cfg.JWT.JwtSecret)

	// 解析令牌时指定自定义Claims
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的签名方法")
		}
		return secret, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New("令牌已过期")
			}
		}
		return nil, err
	}

	// 验证令牌有效性并转换为自定义Claims
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的令牌")
}
