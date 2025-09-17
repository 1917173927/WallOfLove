package controllers

import (
	"github.com/1917173927/WallOfLove/app/services"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/1917173927/WallOfLove/conf/database"
	"github.com/1917173927/WallOfLove/app/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

// UploadImage 处理图片上传
func UploadImage(c *gin.Context) {
	// 获取表单字段
	userIDStr := c.PostForm("user_id")
	postID := c.PostForm("post_id")
	isAvatar := c.PostForm("is_avatar")

	// 验证必填字段
	if userIDStr == "" {
		utils.JsonErrorResponse(c, 400, "user_id 是必填字段")
		return
	}

	// 转换 userID 为 uint64
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		utils.JsonErrorResponse(c, 400, "user_id 必须是一个有效的数字")
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		utils.JsonErrorResponse(c, 400, "获取文件失败")
		return
	}

	// 调用服务层上传图片
	// 根据 userID 查询 username
	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		utils.JsonErrorResponse(c, 400, "用户不存在")
		return
	}
	image, err := services.UploadImage(database.DB, c, userID, user.Username, postID, isAvatar, file)
	if err != nil {
		utils.JsonErrorResponse(c, 500, err.Error())
		return
	}

	utils.JsonSuccessResponse(c, image)
}