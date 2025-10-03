package controllers

import (
	"github.com/1917173927/WallOfLove/app/apiException"
	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/app/services"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/gin-gonic/gin"
)

type ReplyData struct {
	ReviewID  uint64    `json:"review_id" binding:"required"`
	Content   string    `json:"content" binding:"required"`
}
//创建回复
func CreateReply(c *gin.Context) {
	var req ReplyData
	uid, _ := c.Get("userID")
	UID := uid.(uint64)
	if err := c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c,apiException.ParamError,err)
		return
	}
	_,err := services.GetReviewByReviewID(req.ReviewID)
	if err != nil {
		apiException.AbortWithException(c,apiException.TargetError,err)
		return
	}
	if req.Content == "" {
		apiException.AbortWithException(c,apiException.EmptyError,nil)
		return
	}

	if err := services.CreateReply(&models.Reply{
		UserID:  UID,
		ReviewID:  req.ReviewID,
		Content: req.Content,
	}); err != nil {
		apiException.AbortWithException(c,apiException.ServerError,err)
		return
	}
	utils.JsonSuccessResponse(c, req)
}
//删除回复
type DeleteReplyData struct {
	ReplyID uint64 `json:"reply_id" binding:"required"`
}
func DeleteReply(c *gin.Context) {
	uid, _ := c.Get("userID")
	UID := uid.(uint64)
	var req DeleteReplyData
	err := c.ShouldBindJSON(&req)
	if err != nil {
		apiException.AbortWithException(c,apiException.ParamError,err)
		return
	}
	originalReply, err := services.GetReplyByReplyID(req.ReplyID)
	if err != nil {
		apiException.AbortWithException(c,apiException.TargetError,err)
		return
	}
	// 检查是否有权限删除评论
	if originalReply.UserID != UID {
		apiException.AbortWithException(c,apiException.NotPermission,nil)
		return
	}
	err = services.DeleteReply(req.ReplyID)
	if err != nil {
		apiException.AbortWithException(c,apiException.ServerError,err)
		return
	}
	utils.JsonSuccessResponse(c, nil)
}

//获得评论的回复列表
type GetReplyData struct {
	ReviewID uint64 `form:"review_id" binding:"required"`
	Page     int    `form:"page"`
}
type ReplyList struct {
	Replies []models.Reply `json:"replies"`
	Total   int64          `json:"total"`
}
func GetRepliesByReviewID(c *gin.Context) {
	const pageSize = 10 
	var req GetReplyData
	uid, _ := c.Get("userID")
	UID := uid.(uint64)
	err := c.ShouldBind(&req)
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	list, total, err := services.GetRepliesByReviewID(req.ReviewID, UID, req.Page, pageSize)
	if err != nil {
		apiException.AbortWithException(c, apiException.TargetError, err)
		return
	}

	data := ReplyList{
		Replies: list,
		Total:   total,
	}

	utils.JsonSuccessResponse(c, data)
}
