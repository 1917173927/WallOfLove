// Package controllers 包含所有 HTTP 请求处理逻辑，负责接收请求并调用服务层处理业务逻辑。
package controllers

import (
	"github.com/1917173927/WallOfLove/app/apiException"
	"github.com/1917173927/WallOfLove/app/services"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/gin-gonic/gin"
)

// BlackListData 定义拉黑请求的数据结构
type BlackListData struct {
	BlockedID uint64 `json:"blocked_id" binding:"required"` // 被拉黑的用户 ID
}

// BlackUser 处理用户拉黑请求，执行以下操作：
// 1. 从请求上下文中获取用户 ID
// 2. 解析请求体并验证参数
// 3. 检查是否尝试拉黑自己（非法操作）
// 4. 调用服务层执行拉黑操作
// 5. 返回操作结果
func BlackUser(c *gin.Context) {
	uid, _ := c.Get("userID")
	UID := uid.(uint64)
	var req BlackListData
	err := c.ShouldBindJSON(&req)
	if err != nil {
		apiException.AbortWithException(c,apiException.ParamError,err)
		return
	}
	if UID == req.BlockedID {
		apiException.AbortWithException(c,apiException.IllegalTarget,nil)
		return
	}
	err = services.BlackUser(UID, req.BlockedID)
	if err != nil {
		apiException.AbortWithException(c,apiException.ServerError,err)
		return
	}
	utils.JsonSuccessResponse(c, nil)
}

// UnblackUser 处理取消拉黑请求，执行以下操作：
// 1. 从请求上下文中获取用户 ID
// 2. 解析请求体并验证参数
// 3. 调用服务层执行取消拉黑操作
// 4. 返回操作结果
func UnblackUser(c *gin.Context) {
	uid, _ := c.Get("userID")
	UID := uid.(uint64)
	var req BlackListData
	err := c.ShouldBindJSON(&req)
	if err != nil {
		apiException.AbortWithException(c,apiException.ParamError,err)
		return
	}
	err = services.UnblackUser(UID, req.BlockedID)
	if err != nil {
		apiException.AbortWithException(c,apiException.ServerError,err)
		return
	}
	utils.JsonSuccessResponse(c, nil)
}

// GetBlackList 处理获取拉黑列表请求，执行以下操作：
// 1. 从请求上下文中获取用户 ID
// 2. 调用服务层获取拉黑用户列表
// 3. 返回拉黑列表
func GetBlackList(c *gin.Context) {
	uid, _ := c.Get("userID")
	UID := uid.(uint64)
	users, err := services.GetBlackedUsers(UID)
if err != nil {
	apiException.AbortWithException(c,apiException.ServerError,err)
	return
}
utils.JsonSuccessResponse(c,users)
}
