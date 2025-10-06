package controllers

import (
	"github.com/1917173927/WallOfLove/app/apiException"
	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/app/services"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/gin-gonic/gin"
)

type updateProfileData struct {
	Nickname         string `json:"nickname"`
	Username         string `json:"username"`
	OriginalPassword string `json:"original_password"`
	Password         string `json:"password" binding:"omitempty,min=8,max=16"`
	Gender           *int   `json:"gender" binding:"omitempty,oneof=0 1 2"`
	Signature        string `json:"signature" binding:"max=80"`
	AvatarPath       string `json:"avatar_path"`
}

func UpdateProfile(c *gin.Context) {
	uid, _ := c.Get("userID")
	userID := uid.(uint64)
	var req updateProfileData
	if err := c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c, apiException.PwdOrParamError, err)
		return
	}
	user, err := services.GetAllUserDataByID(userID)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	if req.Nickname == "" {
		req.Nickname = user.Nickname
	}
	if req.Username == "" {
		req.Username = user.Username
	}
	if req.Password != "" {
		if err := services.CompareHash(req.OriginalPassword, user.Password); err != nil {
			apiException.AbortWithException(c, apiException.NoThatPasswordOrWrong, err)
			return
		}
		newHash, err := services.HashPassword(req.Password)
		if err != nil {
			apiException.AbortWithException(c, apiException.ServerError, err)
			return
		}
		req.Password = newHash
	} else {
		req.Password = user.Password
	}
	gender := user.Gender
	if req.Gender != nil {
		gender = *req.Gender
	}
	if req.Signature == "" {
		req.Signature = user.Signature
	}
	if req.AvatarPath == "" {
		req.AvatarPath = user.AvatarPath
	}
	updatedUser := models.User{
		ID:         userID,
		Nickname:   req.Nickname,
		Username:   req.Username,
		Password:   req.Password,
		Gender:     gender,
		Signature:  req.Signature,
		AvatarPath: req.AvatarPath,
	}
	if err := services.UpdateProfile(&updatedUser); err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	utils.JsonSuccessResponse(c, nil)
}

type getProfileData struct {
	ID uint64 `form:"id" binding:"required"`
}

type ProfileData struct {
	Profiles   []models.User `json:"profiles"`
	Permission bool          `json:"permission"`
}

func GetProfile(c *gin.Context) {
	uid, _ := c.Get("userID")
	userID := uid.(uint64)
	var req getProfileData
	if err := c.ShouldBind(&req); err != nil {
		apiException.AbortWithException(c, apiException.PwdOrParamError, err)
		return
	}
	permission := userID == req.ID
	user, err := services.GetUserDataByID(req.ID)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	out := ProfileData{
		Profiles:   []models.User{*user},
		Permission: permission,
	}
	utils.JsonSuccessResponse(c, out)
}
