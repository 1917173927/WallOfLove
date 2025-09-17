package controllers

import (
	"github.com/1917173927/WallOfLove/app/services"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/gin-gonic/gin"
)

type BlackListData struct {
	BlockedID       uint64 `json:"blocked_id"`
}
//拉黑
func BlackUser(c *gin.Context) {
	uid,_:=c.Get("userID")
	UID:=uid.(uint64)
	var req BlackListData
	err := c.ShouldBindJSON(&req); 
	if err != nil {
		utils.JsonErrorResponse(c, 501, "参数错误")
		return
	}
	if UID==req.BlockedID {
		utils.JsonErrorResponse(c, 511, "不能拉黑自己")
		return
	}
	err = services.BlackUser(UID, req.BlockedID)
	if err != nil {
		utils.JsonErrorResponse(c, 512, "拉黑失败")
		return
	}
	utils.JsonSuccessResponse(c, nil)
}

//取消拉黑
func UnblackUser(c *gin.Context) {
	uid,_:=c.Get("userID")
	UID:=uid.(uint64)
	var req BlackListData
	err := c.ShouldBindJSON(&req); 
	if err != nil {
		utils.JsonErrorResponse(c, 501, "参数错误")
		return
	}
	err = services.UnblackUser(UID, req.BlockedID)
	if  err != nil {
		utils.JsonErrorResponse(c, 513, "取消拉黑失败")
		return
	}
	utils.JsonSuccessResponse(c, nil)
}

//获取拉黑列表
func GetBlackList(c *gin.Context) {
	uid,_:=c.Get("userID")
	UID:=uid.(uint64)
	blackList,err:=services.GetBlackListID(UID)
	if err!=nil {
		utils.JsonErrorResponse(c, 514, "获取拉黑列表失败")
		return
	}
	utils.JsonSuccessResponse(c, blackList)
}