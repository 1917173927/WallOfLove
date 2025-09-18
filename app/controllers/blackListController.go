package controllers

import (
	"github.com/1917173927/WallOfLove/app/apiException"
	"github.com/1917173927/WallOfLove/app/services"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/gin-gonic/gin"
)

type BlackListData struct {
	BlockedID uint64 `json:"blocked_id"`
}

// 拉黑
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

// 取消拉黑
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

// 获取拉黑列表
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
