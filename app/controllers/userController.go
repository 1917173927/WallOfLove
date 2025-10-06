package controllers

import (
	"errors"

	"github.com/1917173927/WallOfLove/app/apiException"
	"github.com/1917173927/WallOfLove/app/middleware"
	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/app/services"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/gin-gonic/gin"

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
	var req RegisterData
	if err := c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	if err := services.CheckUsername(req.Username); err == nil {
		apiException.AbortWithException(c, apiException.UserAlreadyExisted, err)
		return
	} else if err != gorm.ErrRecordNotFound {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	hash, err := services.HashPassword(req.Password)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	req.Password = hash
	if req.AvatarPath == "" {
		req.AvatarPath = "images/default/default.jpg"
	}
	if err := services.Register(models.User{
		Username:   req.Username,
		Nickname:   req.Nickname,
		Password:   hash,
		Gender:     req.Gender,
		Signature:  "这个人很神秘，什么都没有写",
		AvatarPath: req.AvatarPath,
	}); err != nil {
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
	var req LoginData
	if err := c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	user, err := services.GetUser(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			apiException.AbortWithException(c, apiException.NotFindUser, err)
		} else {
			apiException.AbortWithException(c, apiException.ServerError, err)
		}
		return
	}
	if err := services.CompareHash(req.Password, user.Password); err != nil {
		apiException.AbortWithException(c, apiException.NoThatPasswordOrWrong, err)
		return
	}
	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	utils.JsonSuccessResponse(c, Logdata{
		ID:    uint64(user.ID),
		Token: token,
	})
}
