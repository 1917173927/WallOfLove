// Package controllers 包含所有 HTTP 请求处理逻辑，负责接收请求并调用服务层处理业务逻辑。
package controllers

import (
	"errors"

	"github.com/1917173927/WallOfLove/app/apiException"
	"github.com/1917173927/WallOfLove/app/services"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/gin-gonic/gin"
)

// UploadImage 处理图片上传请求，执行以下操作：
// 1. 从请求上下文中获取用户 ID 和表单字段（post_id 和 is_avatar）
// 2. 获取上传的文件并验证其有效性
// 3. 根据用户 ID 查询用户名
// 4. 调用服务层上传图片，并处理可能的错误（如文件大小超限、文件类型无效等）
// 5. 返回上传成功的响应
func UploadImage(c *gin.Context) {
	uid, _ := c.Get("userID")
	UID := uid.(uint64)
	// 获取表单字段
	postID := c.PostForm("post_id")
	isAvatar := c.PostForm("is_avatar")

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		apiException.AbortWithException(c,apiException.ServerError,err)
		return
	}

	// 据userID查username
	user, err := services.GetUserDataByID(UID)
	if err != nil {
		apiException.AbortWithException(c,apiException.ServerError,err)
		return
	}
	image, err := services.UploadImage(c, UID, user.Username, postID, isAvatar, file)
	if err != nil {
		if errors.Is(err, apiException.ImageSizeExceeded) {
			apiException.AbortWithException(c, apiException.FileSizeExceedError, err)
		} else if errors.Is(err, apiException.ImageTypeInvalid) {
			apiException.AbortWithException(c, apiException.ImageFormatError, err)
		} else if errors.Is(err, apiException.NotImage) {
			apiException.AbortWithException(c, apiException.FileNotImageError, err)
		} else {
			apiException.AbortWithException(c, apiException.UploadFileError, err)
		}
		return
	}

	utils.JsonSuccessResponse(c, image)
}

// DeleteImage 处理图片删除请求
// 1. 获取用户ID和图片ID
// 2. 调用服务层删除图片
// 3. 处理可能的错误
func DeleteImage(c *gin.Context) {
	uid, _ := c.Get("userID")
	UID := uid.(uint64)
	
	var req struct {
		ImageID uint64 `json:"image_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c,apiException.ParamError,err)
		return
	}

	if err := services.DeleteImage(req.ImageID, UID); err != nil {
		if errors.Is(err, apiException.ImageNotFound) {
			apiException.AbortWithException(c, apiException.TargetError, err)
		} else if errors.Is(err, apiException.NotPermission) {
			apiException.AbortWithException(c, apiException.NotPermission, err)
		} else {
			apiException.AbortWithException(c, apiException.ServerError, err)
		}
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

