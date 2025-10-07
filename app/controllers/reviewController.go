package controllers

import (
	"github.com/1917173927/WallOfLove/app/apiException"
	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/app/services"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/gin-gonic/gin"
)

// 创建评论
type ReviewData struct {
	PostID  uint64 `json:"post_id" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func CreateReview(c *gin.Context) {
	var req ReviewData
	uid, _ := c.Get("userID")
	userID := uid.(uint64)
	if err := c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	if _, err := services.GetPostDataByID(req.PostID); err != nil {
		apiException.AbortWithException(c, apiException.TargetError, err)
		return
	}
	if req.Content == "" {
		apiException.AbortWithException(c, apiException.EmptyError, nil)
		return
	}
	if err := services.CreateReview(&models.Review{
		UserID:  userID,
		PostID:  req.PostID,
		Content: req.Content,
	}); err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	utils.JsonSuccessResponse(c, req)
}

// 删除评论
type DeleteReviewData struct {
	ReviewID uint64 `json:"review_id" binding:"required"`
}

func DeleteReview(c *gin.Context) {
	uid, _ := c.Get("userID")
	userID := uid.(uint64)
	var req DeleteReviewData
	if err := c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	review, err := services.GetReviewByReviewID(req.ReviewID)
	if err != nil {
		apiException.AbortWithException(c, apiException.TargetError, err)
		return
	}
	if review.UserID != userID {
		apiException.AbortWithException(c, apiException.NotPermission, nil)
		return
	}
	if err := services.DeleteReview(req.ReviewID); err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	utils.JsonSuccessResponse(c, nil)
}

// 获取评论列表
type GetReviewsByPostIDData struct {
	PostID   uint64 `form:"post_id" binding:"required"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}
type ReviewList struct {
	Reviews []services.ReviewWithLike `json:"reviews"`
	Total   int64                     `json:"total"`
}

func GetReviewsByPostID(c *gin.Context) {
	var req GetReviewsByPostIDData
	uid, _ := c.Get("userID")
	userID := uid.(uint64)
	if err := c.ShouldBind(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if err := services.GetReviewsByPostID(req.PostID); err != nil {
		apiException.AbortWithException(c, apiException.TargetError, err)
		return
	}
	reviews, total, err := services.GetVisibleReviews(req.PostID, userID, req.Page, req.PageSize)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	data := ReviewList{
		Reviews: reviews,
		Total:   total,
	}
	utils.JsonSuccessResponse(c, data)
}
