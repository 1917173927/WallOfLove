package controllers

import (
	"github.com/1917173927/WallOfLove/app/apiException"
	"github.com/1917173927/WallOfLove/app/services"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/gin-gonic/gin"
)

type LikeData struct {
	PostID   uint64 `json:"post_id" binding:"required"`
	ReviewID uint64 `json:"review_id"`
}
// 点赞
func LikePost(c *gin.Context) {
	uid, _ := c.Get("userID")
	UID := uid.(uint64)
	var req LikeData
	if err := c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c,apiException.ParamError,err)
		return
	}
	liked,err := services.Like(UID, req.PostID,req.ReviewID)
	if err != nil {
		apiException.AbortWithException(c,apiException.ServerError,err)
		return
	}
	utils.JsonSuccessResponse(c, liked)//false取消点赞，true点赞
}
