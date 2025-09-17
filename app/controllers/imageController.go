package controllers

import (
	"github.com/1917173927/WallOfLove/app/services"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/gin-gonic/gin"
)

// UploadImage 处理图片上传
func UploadImage(c *gin.Context) {
	uid, _ := c.Get("userID")
	UID := uid.(uint64)
	// 获取表单字段
	postID := c.PostForm("post_id")
	isAvatar := c.PostForm("is_avatar")

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		utils.JsonErrorResponse(c, 400, "获取文件失败")
		return
	}

	// 据userID查username
	user, err := services.GetUserDataByID(UID)
	if err != nil {
		utils.JsonErrorResponse(c, 400, "用户不存在")
		return
	}
	image, err := services.UploadImage(c, UID, user.Username, postID, isAvatar, file)
	if err != nil {
		utils.JsonErrorResponse(c, 500, err.Error())
		return
	}

	utils.JsonSuccessResponse(c, image)
}