package middleware

import (
	"errors"
	"time"

	"github.com/1917173927/WallOfLove/app/apiException"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("ohyeahmambo")

func GenerateToken(userID uint64) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if len(auth) < 7 || auth[:7] != "Bearer " {
			apiException.AbortWithException(c, apiException.NotLogin, nil)
			c.Abort()
			return
		}
		tokenStr := auth[7:]
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			if errors.Is(err, jwt.ErrTokenExpired) {
				apiException.AbortWithException(c, apiException.TokenExpired, err)
			} else {
				apiException.AbortWithException(c, apiException.NotLogin, nil)
			}
			c.Abort()
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userID := uint64(claims["userID"].(float64))
			c.Set("userID", userID)
			c.Next()
		} else {
			apiException.AbortWithException(c, apiException.NotLogin, nil)
			c.Abort()
		}
	}
}
