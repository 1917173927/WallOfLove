package controllers

import (
	"errors"

	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/app/services"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type updateProfileData struct {
	Nickname         string  `json:"nickname"`
	Username         string  `json:"username"`
	OriginalPassword string  `json:"original_password"`
	Password         string  `json:"password" binding:"pwdmin"`
	AvatarID         *uint64 `json:"avatar_id"`
}

func UpdateProfile(c *gin.Context) {
	uid, _ := c.Get("userID")
	UID := uid.(uint64)

	var req updateProfileData
	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.JsonErrorResponse(c, 501, "参数错误")
		return
	}

	user, err := services.GetUserDataByID(UID)
	if err != nil {
		c.Error(errors.New("用户不存在"))
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
			c.Error(errors.New("密码错误"))
			return
		}
		newHash, err := services.HashPassword(req.Password)
		if err != nil {
			c.Error(errors.New("加密失败"))
			return
		}
		req.Password = newHash
	} else {
		req.Password = user.Password
	}

	if req.AvatarID == nil {
		req.AvatarID = user.AvatarImageID
	}

	updatedUser := models.User{
		ID:            UID,
		Nickname:      req.Nickname,
		Username:      req.Username,
		Password:      req.Password,
		AvatarImageID: req.AvatarID,
		Version:       user.Version,
	}

	err = services.UpdateProfile(&updatedUser, user.Version)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Error(errors.New("数据已被其他会话修改，请重试"))
		} else {
			c.Error(errors.New("更新用户信息失败"))
		}
		return
	}

	c.JSON(200, map[string]any{"version": user.Version + 1})
}
