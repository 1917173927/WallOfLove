package middleware

import (
	"errors"
	"net/http"
	"time"

	"github.com/1917173927/WallOfLove/app/utils"
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
			utils.JsonErrorResponse(c, 401, "缺少token")
			c.Abort()
			return
		}
		// 去掉 "Bearer "
		if len(auth) < 7 || auth[:7] != "Bearer " {
			utils.JsonErrorResponse(c, 401, "token格式错误")
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
		      utils.JsonErrorResponse(c,402, "token已过期")
	      } else {
		      utils.JsonErrorResponse(c,401, "token无效")
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
			utils.JsonErrorResponse(c, 401, "token claims错误")
			c.Abort()
		}
	}
}

//访问不存在地址时404
func HandleNotFound(c *gin.Context) {
	utils.JsonResponse(c, 404, 200404, http.StatusText(http.StatusNotFound), nil)
}