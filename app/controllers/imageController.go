package controllers

import (
	"errors"

	"github.com/1917173927/WallOfLove/app/apiException"
	"github.com/1917173927/WallOfLove/app/services"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/gin-gonic/gin"
)

func UploadImage(c *gin.Context) {
	uid, _ := c.Get("userID")
	userID := uid.(uint64)
	postID := c.PostForm("post_id")
	isAvatar := c.PostForm("is_avatar")
	file, err := c.FormFile("file")
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	user, err := services.GetUserDataByID(userID)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	image, err := services.UploadImage(c, userID, user.Username, postID, isAvatar, file)
	if err != nil {
		switch {
		case errors.Is(err, apiException.ImageSizeExceeded):
			apiException.AbortWithException(c, apiException.FileSizeExceedError, err)
		case errors.Is(err, apiException.ImageTypeInvalid):
			apiException.AbortWithException(c, apiException.ImageFormatError, err)
		case errors.Is(err, apiException.NotImage):
			apiException.AbortWithException(c, apiException.FileNotImageError, err)
		default:
			apiException.AbortWithException(c, apiException.UploadFileError, err)
		}
		return
	}
	utils.JsonSuccessResponse(c, image)
}

func DeleteImage(c *gin.Context) {
	uid, _ := c.Get("userID")
	userID := uid.(uint64)
	var req struct {
		ImageID uint64 `json:"image_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	if err := services.DeleteImage(req.ImageID, userID); err != nil {
		switch {
		case errors.Is(err, apiException.ImageNotFound):
			apiException.AbortWithException(c, apiException.TargetError, err)
		case errors.Is(err, apiException.NotPermission):
			apiException.AbortWithException(c, apiException.NotPermission, err)
		default:
			apiException.AbortWithException(c, apiException.ServerError, err)
		}
		return
	}
	utils.JsonSuccessResponse(c, nil)
}
