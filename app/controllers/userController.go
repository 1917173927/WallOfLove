package controllers

import (
	"errors"
	"strings"

	"github.com/1917173927/WallOfLove/app/apiException"
	"github.com/1917173927/WallOfLove/app/middleware"
	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/app/services"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type RegisterData struct {
	Username   string `json:"username"        binding:"required"`
	Nickname   string `json:"nickname"        binding:"required"`
	Password   string `json:"password"        binding:"required,min=8,max=16"`
	Gender     int    `json:"gender"          binding:"min=0,max=2"` //0:男，1:女，2:保密
	AvatarPath string `json:"avatar_path"`
}

// 注册
func Register(c *gin.Context) {
	var data RegisterData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	//判断账号是否已经存在
	err = services.CheckUsername(data.Username)
	if err == nil {
		apiException.AbortWithException(c, apiException.UserAlreadyExisted, err)
		return
	} else if err != gorm.ErrRecordNotFound {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	//哈希加密密码
	hash, err := services.HashPassword(data.Password)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	data.Password = hash
	//若未上传头像，则用默认头像
	if data.AvatarPath == "" {
		defaultID := "images/default/default.jpg"
		data.AvatarPath = defaultID
	}
	//注册用户
	err = services.Register(models.User{
		Username:   data.Username,
		Nickname:   data.Nickname,
		Password:   hash,
		Gender:     data.Gender,
		Signature:  "这个人很神秘，什么都没有写",
		AvatarPath: data.AvatarPath,
	})
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
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
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	//检查是否有此用户
	user, err := services.GetUser(data.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			apiException.AbortWithException(c, apiException.NotFindUser, err)
		} else {
			apiException.AbortWithException(c, apiException.ServerError, err)
		}
		return
	}
	// 密码比对
	if err := services.CompareHash(data.Password, user.Password); err != nil {
		apiException.AbortWithException(c, apiException.NoThatPasswordOrWrong, err)
		return
	}
	//生成token
	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	logdata := Logdata{
		ID:    uint64(user.ID),
		Token: token,
	}
	utils.JsonSuccessResponse(c, logdata)
}
