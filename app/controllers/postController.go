package controllers

import (
	"errors"
	"log"

	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/app/services"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostData struct{
	Content    string `json:"content" binding:"required"`
	Anonymous  bool   `json:"anonymous"`
	Visibility bool `json:"visibility" binding:"required,oneof=public private"`
}

func CreatePost(c *gin.Context) {
	var req PostData
	uid, _ := c.Get("userID")
	UID := uid.(uint64)   //jwt
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JsonErrorResponse(c, 400, "参数错误")
		return
	}

	if UID == 0 {
		utils.JsonErrorResponse(c, 400, "用户ID不能为空")
		return
	}
	
	if req.Content == "" {
		utils.JsonErrorResponse(c, 400, "帖子内容不能为空")
		return
	}
	if err := services.CreatePost(&models.Post{
		UserID:    UID,
		Content:   req.Content,
		Anonymous: req.Anonymous,
		Visibility: req.Visibility,
	}); err != nil {
		utils.JsonErrorResponse(c, 500, "创建帖子失败")
		return
	}

	utils.JsonSuccessResponse(c, req)
}

func UpdatePost(c *gin.Context) {
	var req models.Post
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JsonErrorResponse(c, 400, "参数错误")
		return
	}

	if req.ID == 0 {
		utils.JsonErrorResponse(c, 400, "帖子ID不能为空")
		return
	}
	if req.UserID == 0 {
		utils.JsonErrorResponse(c, 400, "用户ID不能为空")
		return
	}

	if err := services.UpdatePost(&req); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.JsonErrorResponse(c, 404, "帖子不存在")
		} else {
			utils.JsonErrorResponse(c, 500, "更新帖子失败")
		}
		return
	}

	utils.JsonSuccessResponse(c, req)
}

func DeletePost(c *gin.Context) {
	postID := c.Param("id")
	if postID == "" {
		utils.JsonErrorResponse(c, 400, "无效的帖子ID")
		return
	}

	if err := services.DeletePost(postID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.JsonErrorResponse(c, 404, "帖子不存在")
		} else {
			utils.JsonErrorResponse(c, 500, "删除帖子失败")
		}
		return
	}

	utils.JsonSuccessResponse(c, gin.H{"message": "删除成功"})
}

// GetVisiblePosts 获取未被拉黑的其他人发布的表白
func GetVisiblePosts(c *gin.Context) {
	uid, _ := c.Get("userID")
	UID := uid.(uint64)
	posts,total,err := services.GetVisiblePosts(UID,10,10)
	if err != nil {
		utils.JsonErrorResponse(c, 500, "获取帖子失败")
		return
	}
	log.Println(total)
	utils.JsonSuccessResponse(c, posts)
}