package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mix-go/xutil/xenv"
	"net/http"
	"time"
)

type AuthController struct {
}

func (t *AuthController) Index(c *gin.Context) {
	// 检查用户登录代码
	// ...

	// 创建 token
	now := time.Now().Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "http://example.org",                                  // 签发人
		"iat": now,                                                   // 签发时间
		"exp": now + int64(7200),                                     // 过期时间
		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(), // 什么时间之前不可用
		"uid": 100008,
	})
	tokenString, err := token.SignedString([]byte(xenv.Getenv("HMAC_SECRET").String()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Creation of token fails",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       http.StatusOK,
		"message":      "ok",
		"access_token": tokenString,
		"expire_in":    7200,
	})
}
