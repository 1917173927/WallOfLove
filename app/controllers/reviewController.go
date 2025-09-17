package controllers

import (
	"log"
	"time"

	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/app/services"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

type ReviewData struct {
	ID        uint64      `json:"id" gorm:"primaryKey"`
	UserID    uint64      `json:"user_id"`
	PostID    uint64      `json:"post_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateReview(c *gin.Context) {
	var req ReviewData
	uid, exists := c.Get("userID")
	UID := uid.(uint64)
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
		log.Printf("创建评论失败: %v", err)
		utils.JsonErrorResponse(c, 500, "创建评论失败")
		return
	}
	utils.JsonSuccessResponse(c, req)
}

func GetReviewsByPostID(c *gin.Context) {
	postID := c.Param("postID")
	PostID, err := strconv.ParseUint(postID, 10, 64)
	if err != nil {
		log.Printf("获取评论失败: %v", err)
		utils.JsonErrorResponse(c, 400, "无效的帖子ID")
		return
	}
	if PostID == 0 {
		log.Printf("获取评论失败: %v", err)
		utils.JsonErrorResponse(c, 400, "帖子ID不能为空")
		return
	}
	reviews, err := services.GetReviewsByPostID(PostID)
	if err != nil {
		log.Printf("获取评论失败: %v", err)
		utils.JsonErrorResponse(c, 500, "获取评论失败")
		return
	}
	utils.JsonSuccessResponse(c, reviews)
}