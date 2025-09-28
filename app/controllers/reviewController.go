package controllers

import (
	"time"

	"github.com/1917173927/WallOfLove/app/apiException"
	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/app/services"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/gin-gonic/gin"
)

type ReviewData struct {
	UserID    uint64    `json:"user_id"`
	PostID    uint64    `json:"post_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateReview(c *gin.Context) {
	var req ReviewData
	uid, _ := c.Get("userID")
	UID := uid.(uint64)
	if err := c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c,apiException.ParamError,err)
		return
	}
	err := services.GetReviewsByPostID(req.PostID)
	if err != nil {
		apiException.AbortWithException(c,apiException.TargetError,err)
		return
	}
	if req.Content == "" {
		apiException.AbortWithException(c,apiException.EmptyError,nil)
		return
	}

	if err := services.CreateReview(&models.Review{
		UserID:  UID,
		PostID:  req.PostID,
		Content: req.Content,
	}); err != nil {
		apiException.AbortWithException(c,apiException.ServerError,err)
		return
	}
	utils.JsonSuccessResponse(c, req)
}

type GetReviewsByPostIDData struct {
	PostID uint64 `form:"post_id"`
	Page int `form:"page"`
	PageSize int `form:"page_size"`
}
type ReviewList struct {
	Reviews []services.ReviewWithLike `json:"reviews"`
	Total   int64                     `json:"total"`
}
func GetReviewsByPostID(c *gin.Context) {
	var req GetReviewsByPostIDData
	uid, _ := c.Get("userID")
	UID := uid.(uint64)
	err := c.ShouldBind(&req)
	if err != nil {
		apiException.AbortWithException(c,apiException.ParamError,err)
		return
	}
	err = services.GetReviewsByPostID(req.PostID)
	if err != nil {
		apiException.AbortWithException(c,apiException.TargetError,err)
		return
	}
	reviews, total, err := services.GetVisibleReviews(req.PostID,UID, req.Page, req.PageSize)
	if err != nil {
		apiException.AbortWithException(c,apiException.ServerError,err)
		return
	}
	data := ReviewList{
		Reviews: reviews,
		Total:   total,
	}
	utils.JsonSuccessResponse(c, data)
}
