package controllers

import (
	"github.com/1917173927/WallOfLove/app/apiException"
	"github.com/1917173927/WallOfLove/app/services"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/gin-gonic/gin"
)

type BlackListData struct {
	BlockedID uint64 `json:"blocked_id" binding:"required"`
}

func BlackUser(c *gin.Context) {
	uid, _ := c.Get("userID")
	userID := uid.(uint64)
	var req BlackListData
	if err := c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	if userID == req.BlockedID {
		apiException.AbortWithException(c, apiException.IllegalTarget, nil)
		return
	}
	if err := services.BlackUser(userID, req.BlockedID); err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	utils.JsonSuccessResponse(c, nil)
}

func UnblackUser(c *gin.Context) {
	uid, _ := c.Get("userID")
	userID := uid.(uint64)
	var req BlackListData
	if err := c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	if err := services.UnblackUser(userID, req.BlockedID); err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	utils.JsonSuccessResponse(c, nil)
}

func GetBlackList(c *gin.Context) {
	uid, _ := c.Get("userID")
	userID := uid.(uint64)
	users, err := services.GetBlackedUsers(userID)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	utils.JsonSuccessResponse(c, users)
}
