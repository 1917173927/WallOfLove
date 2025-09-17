package controllers

import (
	"log"
	"time"

	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/app/services"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/gin-gonic/gin"
)

type Review2Data struct {
	UserID    uint64      `json:"user_id"`
	ReviewID  uint64      `json:"review_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateReview2(c *gin.Context) {
	var req Review2Data
	uid, exists := c.Get("userID")
	if !exists {
		utils.JsonErrorResponse(c, 400, "用户ID未提供")
		return
	}
	UID, ok := uid.(uint64)
	if !ok {
		utils.JsonErrorResponse(c, 400, "无效的用户ID类型")
		return
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JsonErrorResponse(c, 400, "参数错误")
		return
	}
	if UID == 0 {
		utils.JsonErrorResponse(c, 400, "用户ID不能为空")
		return
	}
	if req.ReviewID == 0 {
		utils.JsonErrorResponse(c, 400, "评论ID不能为空")
		return
	}
	if req.Content == "" {
		utils.JsonErrorResponse(c, 400, "评论内容不能为空")
		return
	}

	if err := services.CreateReview2(&models.Review2{
		UserID:  UID,
		ReviewID:  req.ReviewID,
		Content: req.Content,
	}); err != nil {
		log.Printf("创建评论失败: %v", err)
		utils.JsonErrorResponse(c, 500, "创建评论失败")
		return
	}
	utils.JsonSuccessResponse(c, req)
}

type GetReviews2ByPostIDData struct {
	ReviewID uint64 `json:"review_id"`
	Page int `json:"page"`
	PageSize int `json:"page_size"`
}

func GetReviews2ByPostID(c *gin.Context) {
	var req GetReviews2ByPostIDData
	uid, _ := c.Get("userID")
	UID := uid.(uint64)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.JsonErrorResponse(c, 501, "参数错误")
		return
	}
	exists := services.GetReview2sByPostID(req.ReviewID)
	if exists != nil {
		utils.JsonErrorResponse(c, 400, "无效的评论ID")
		return
	}
	if req.ReviewID == 0 {
		utils.JsonErrorResponse(c, 400, "评论ID不能为空")
		return
	}
	reviews, total, err := services.GetVisibleReviews(UID, req.Page, req.PageSize)
	if err != nil {
		utils.JsonErrorResponse(c, 500, "获取评论失败")
		return
	}
	data := map[string]any{
		"reviews": reviews,
		"total": total,
	}
	utils.JsonSuccessResponse(c, data)
}