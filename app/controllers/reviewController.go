package controllers

import (


	"time"
	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/app/services"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/gin-gonic/gin"
)

type ReviewData struct {
	UserID    uint64      `json:"user_id"`
	PostID    uint64      `json:"post_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateReview(c *gin.Context) {
	var req ReviewData
	uid, _ := c.Get("userID")
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
	if req.PostID == 0 {
		utils.JsonErrorResponse(c, 400, "帖子ID不能为空")
		return
	}
	if req.Content == "" {
		utils.JsonErrorResponse(c, 400, "评论内容不能为空")
		return
	}

	if err := services.CreateReview(&models.Review{
		UserID:  UID,
		PostID:  req.PostID,
		Content: req.Content,
	}); err != nil {
		utils.JsonErrorResponse(c, 500, "创建评论失败")
		return
	}
	utils.JsonSuccessResponse(c, req)
}

type GetReviewsByPostIDData struct {
	PostID uint64 `json:"post_id"`
	Page int `json:"page"`
	PageSize int `json:"page_size"`
}

func GetReviewsByPostID(c *gin.Context) {
	var req GetReviewsByPostIDData
	uid, _ := c.Get("userID")
	UID := uid.(uint64)
	err := c.ShouldBindJSON(&req)
	if err != nil {
        utils.JsonErrorResponse(c, 404, "未找到该评论")
        return
    }
	if err != nil {
		utils.JsonErrorResponse(c, 501, "参数错误")
		return
	}
	exists := services.GetReviewsByPostID(req.PostID)
	if exists != nil {
		utils.JsonErrorResponse(c, 400, "无效的帖子ID")
		return
	}
	if req.PostID == 0 {
		utils.JsonErrorResponse(c, 400, "帖子ID不能为空")
		return
	}
	reviews, total, err := services.GetVisibleReviews(UID, req.Page, req.PageSize)
	if err != nil {
		utils.JsonErrorResponse(c, 500, "获取帖子失败")
		return
	}
	data := map[string]any{
		"reviews": reviews,
		"total": total,
	}
	utils.JsonSuccessResponse(c, data)
}
