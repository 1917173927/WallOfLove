package middleware

import (
	"errors"
	"time"

	"github.com/1917173927/WallOfLove/app/apiException"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("ohyeahmambo")

// 生成token
func GenerateToken(userID uint64) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// JWT 中间件，获取并检验
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 拿 token
		auth := c.GetHeader("Authorization")
		if auth == "" {
			apiException.AbortWithException(c,apiException.NotLogin,nil)
			c.Abort()
			return
		}
		// 去掉 "Bearer "
		if len(auth) < 7 || auth[:7] != "Bearer " {
			apiException.AbortWithException(c,apiException.NotLogin,nil)
			c.Abort()
			return
		}
		tokenStr := auth[7:]

		// 解析 + 验签
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			//先区分是不是过期
			if errors.Is(err, jwt.ErrTokenExpired) {
				apiException.AbortWithException(c,apiException.TokenExpired,err)
			} else {
				apiException.AbortWithException(c,apiException.NotLogin,nil)
			}
			c.Abort()
			return
		}

		// 取出 userID 写进上下文
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userID := uint64(claims["userID"].(float64))
			c.Set("userID", userID)
			c.Next()
		} else {
			apiException.AbortWithException(c,apiException.NotLogin,nil)
			c.Abort()
		}
	}
}
