package controllers

import (

	"github.com/1917173927/WallOfLove/app/apiException"
	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/app/services"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/gin-gonic/gin"
)
//更新用户信息
type updateProfileData struct {
	Nickname         string  `json:"nickname"`
	Username         string  `json:"username"`
	OriginalPassword string  `json:"original_password"`
	Password         string  `json:"password" binding:"omitempty,min=8,max=16"`
	Gender           *int    `json:"gender" binding:"omitempty,oneof=0 1 2"`
	Signature        string  `json:"signature" binding:"max=80"`
	AvatarPath       string  `json:"avatar_path"`
}
func UpdateProfile(c *gin.Context) {
	uid, _ := c.Get("userID")
	UID := uid.(uint64)

	var req updateProfileData
	err := c.ShouldBindJSON(&req)
	if err != nil {
		apiException.AbortWithException(c,apiException.PwdOrParamError,err)
		return
	}
	//获得用户原信息
	user, err := services.GetAllUserDataByID(UID)
	if err != nil {
		apiException.AbortWithException(c,apiException.ServerError,err)
		return
	}
	//若未填写昵称，则用原值
	if req.Nickname == "" {
		req.Nickname = user.Nickname
	}
	//若未填写用户名，则用原值
	if req.Username == "" {
		req.Username = user.Username
	}
	//若未填写密码，则用原值，若要更改密码，则需填写原密码，并验证原密码是否正确，若正确，则用新密码，否则报错
	if req.Password != "" {
		if err := services.CompareHash(req.OriginalPassword, user.Password); err != nil {
			apiException.AbortWithException(c,apiException.NoThatPasswordOrWrong,err)
			return
		}
		newHash, err := services.HashPassword(req.Password)
		if err != nil {
			apiException.AbortWithException(c,apiException.ServerError,err)
			return
		}
		req.Password = newHash
	} else {
		req.Password = user.Password
	}
	//若未填写性别，则用原值
	gender := user.Gender
	if req.Gender != nil {
		gender = *req.Gender
	}
	//若未填写签名，则用原值
	if req.Signature == "" {
		req.Signature = user.Signature
	}
	//若未填写头像，则用原值
	if req.AvatarPath == "" {
		req.AvatarPath = user.AvatarPath
	}

	updatedUser := models.User{
		ID:            UID,
		Nickname:      req.Nickname,
		Username:      req.Username,
		Password:      req.Password,
		Gender:        gender,
		Signature:     req.Signature,
		AvatarPath:    req.AvatarPath,
	}
	//更新用户信息
	err = services.UpdateProfile(&updatedUser)
	if err != nil {
			apiException.AbortWithException(c,apiException.ServerError,err)
			return
		}
    utils.JsonSuccessResponse(c,nil)
}
//获取用户信息
type getProfileData struct {
	ID uint64 `form:"id" binding:"required"`
}
type ProfileData struct {
	Profiles []models.User `json:"profiles"`
    Permission bool `json:"permission"`
}
func GetProfile(c *gin.Context) {
	uid, _ := c.Get("userID")
	UID := uid.(uint64)
	var req getProfileData
	err := c.ShouldBind(&req)
	if err != nil {
		apiException.AbortWithException(c,apiException.PwdOrParamError,err)
		return
	}
	var permission bool
	if UID==req.ID{
		permission=true
	}else{
		permission=false
	}//ture:有权限，false:无权限
	user, err := services.GetUserDataByID(req.ID)
	if err != nil {
		apiException.AbortWithException(c,apiException.ServerError,err)
		return
	}
	out:=ProfileData{
		Profiles: []models.User{*user},
		Permission: permission,
	}
	utils.JsonSuccessResponse(c, out)
}
