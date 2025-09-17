package middleware

import (
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			switch err.Err.Error() {
			case "参数错误":
				utils.JsonErrorResponse(c, 501, "参数错误")
			case "用户不存在":
				utils.JsonErrorResponse(c, 504, "用户不存在")
			case "密码错误":
				utils.JsonErrorResponse(c, 505, "密码错误")
			case "数据已被其他会话修改，请重试":
				utils.JsonErrorResponse(c, 507, "数据已被其他会话修改，请重试")
			case "更新用户信息失败":
				utils.JsonErrorResponse(c, 508, "更新用户信息失败")
			default:
				utils.JsonInternalServerErrorResponse(c)
			}
		}
	}
}
