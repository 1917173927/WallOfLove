package controllers

import (
	"errors"

	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/app/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreatePost(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.Error(errors.New("参数错误"))
		return
	}

	if err := (&services.PostService{}).CreatePost(&post); err != nil {
		if errors.Is(err, gorm.ErrInvalidData) {
			c.Error(errors.New("图片数量超过限制"))
		} else {
			c.Error(errors.New("创建帖子失败"))
		}
		return
	}

	c.JSON(200, post)
}

func UpdatePost(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.Error(errors.New("参数错误"))
		return
	}

	if err := (&services.PostService{}).UpdatePost(&post); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Error(errors.New("帖子不存在"))
		} else {
			c.Error(errors.New("更新帖子失败"))
		}
		return
	}

	c.JSON(200, post)
}

func DeletePost(c *gin.Context) {
	postID := c.Param("id")
	if err := (&services.PostService{}).DeletePost(postID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Error(errors.New("帖子不存在"))
		} else {
			c.Error(errors.New("删除帖子失败"))
		}
		return
	}

	c.JSON(200, nil)
}