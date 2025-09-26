package controllers

import (
	"github.com/1917173927/WallOfLove/app/apiException"
	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/app/services"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/gin-gonic/gin"
)

type ReplyData struct {
	UserID    uint64    `json:"user_id"`
	ReviewID  uint64    `json:"review_id"`
	Content   string    `json:"content"`
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
	err := services.GetReviewByReviewID(req.ReviewID)
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


type GetReplyData struct {
	ReviewID uint64 `form:"review_id"`
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
