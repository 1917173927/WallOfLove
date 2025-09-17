package controllers

import (
	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/app/services"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/gin-gonic/gin"
)

type Review2Data struct {
	UserID    uint64      `json:"user_id"`
	ReviewID  uint64      `json:"review_id"`
	Content   string    `json:"content"`
}

func CreateReview2(c *gin.Context) {
	var req Review2Data
	uid, _ := c.Get("userID")
	UID := uid.(uint64)
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JsonErrorResponse(c, 400, "参数错误")
		return
	}
	if req.ReviewID == 0 {
		utils.JsonErrorResponse(c, 400, "评论ID不能为空")
		return
	}
	if req.Content == "" {
		utils.JsonErrorResponse(c, 400, "回复内容不能为空")
		return
	}

	if err := services.CreateReview2(&models.Review2{
		UserID:  UID,
		ReviewID:  req.ReviewID,
		Content: req.Content,
	}); err != nil {
		utils.JsonErrorResponse(c, 500, "创建回复失败")
		return
	}
	utils.JsonSuccessResponse(c, req)
}