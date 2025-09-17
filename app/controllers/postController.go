package controllers

import (
	"errors"

	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/app/services"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/1917173927/WallOfLove/conf/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	postService = &services.PostService{}
	errPostNotFound = errors.New("帖子不存在")
	errInvalidPostData = errors.New("无效的帖子数据")//错误统一处理还没做
	errUpdateFailed = errors.New("更新帖子失败")
)

func CreatePost(c *gin.Context) {
	var req models.Post
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JsonErrorResponse(c, 400, "参数错误")
		return
	}

	if req.UserID == 0 {
		utils.JsonErrorResponse(c, 400, "用户ID不能为空")
		return
	}
	if req.Content == "" {
		utils.JsonErrorResponse(c, 400, "帖子内容不能为空")
		return
	}
	if req.Visibility != "public" && req.Visibility != "private" {
		utils.JsonErrorResponse(c, 400, "visibility 必须是 'public' 或 'private'")
		return
	}

	if err := postService.CreatePost(&req); err != nil {
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

	if err := postService.UpdatePost(&req); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.JsonErrorResponse(c, 404, errPostNotFound.Error())
		} else {
			utils.JsonErrorResponse(c, 500, errUpdateFailed.Error())
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

	if err := postService.DeletePost(postID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.JsonErrorResponse(c, 404, errPostNotFound.Error())
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
	var posts []models.Post
	filter := services.FilterBlack(c, database.DB, UID)
	err := filter.Find(&posts)
	if err != nil {
		utils.JsonErrorResponse(c, 500, "获取表白失败")
		return
	}

	utils.JsonSuccessResponse(c, posts)
}