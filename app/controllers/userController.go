package controllers

import (
	"errors"

	"github.com/1917173927/WallOfLove/app/middleware"
	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/app/services"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RegisterData struct {
	Username      string  `json:"username"        binding:"required"`
	Name          string  `json:"name"            binding:"required"`
	Password      string  `json:"password"        binding:"required,pwdmin"`
	AvatarImageID *uint64 `json:"avatar_image_id"`
}

// 注册
func Register(c *gin.Context) {
	var data RegisterData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, 501, "参数错误")
		return
	}
	//判断账号是否已经存在
	err = services.CheckUsername(data.Username)
	if err == nil {
		utils.JsonErrorResponse(c, 502, "账号已被注册")
		return
	} else if err != gorm.ErrRecordNotFound {
		utils.JsonInternalServerErrorResponse(c)
		return
	}
	//哈希加密密码
	hash, err := services.HashPassword(data.Password)
	if err != nil {
		utils.JsonErrorResponse(c, 503, "加密失败")
		return
	}
	data.Password = hash
	//若未上传头像，则用默认头像
	if data.AvatarImageID == nil {
		defaultID := uint64(1)
		data.AvatarImageID = &defaultID
	}
	//注册用户
	err = services.Register(models.User{
		Username:      data.Username,
		Nickname:      data.Name,
		Password:      hash,
		AvatarImageID: data.AvatarImageID,
	})
	if err != nil {
		utils.JsonInternalServerErrorResponse(c)
		return
	}
	utils.JsonSuccessResponse(c, nil)
}

// 登录
type LoginData struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type Logdata struct {
	ID    uint64 `json:"user_id"`
	Token string `json:"token"`
}

// 接收参数
func Login(c *gin.Context) {
	var data LoginData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, 501, "参数错误")
		return
	}
	user, err := services.GetUser(data.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.JsonErrorResponse(c, 504, "用户不存在")
		} else {
			utils.JsonInternalServerErrorResponse(c)
		}
		return
	}

	if err := services.CompareHash(data.Password, user.Password); err != nil {
		utils.JsonErrorResponse(c, 505, "密码错误")
		return
	}
	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		utils.JsonErrorResponse(c, 506, "生成token失败")
		return
	}
	logdata := Logdata{
		ID:    uint64(user.ID),
		Token: token,
	}
	utils.JsonSuccessResponse(c, logdata)
}
