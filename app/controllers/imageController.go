package controllers

import (
	"errors"

	"github.com/1917173927/WallOfLove/app/apiException"
	"github.com/1917173927/WallOfLove/app/services"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/1917173927/WallOfLove/app/utils/errno"
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
		apiException.AbortWithException(c,apiException.UploadFileError,err)
		return
	}else if errors.Is(err,errno.ErrImageSizeExceeded){
		apiException.AbortWithException(c,apiException.FileSizeExceedError,err)
		return
	}else if errors.Is(err,errno.ErrImageTypeInvalid){
		apiException.AbortWithException(c,apiException.ImageFormatError,err)
		return
	}else if errors.Is(err,errno.ErrNotImage){
		apiException.AbortWithException(c,apiException.FileNotImageError,err)
		return
	}

	utils.JsonSuccessResponse(c, image)
}
