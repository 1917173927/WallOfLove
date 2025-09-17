package controllers

import (
	"errors"
	"log"

	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/app/services"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)
// 创建帖子
type PostData struct{
	Content    string `json:"content" binding:"required"`
	Anonymous  bool   `json:"anonymous"`
	Visibility bool   `json:"visibility"`
}

func CreatePost(c *gin.Context) {
	var req PostData
	uid, _ := c.Get("userID")
	UID := uid.(uint64)   //jwt
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JsonErrorResponse(c, 400, "参数错误")
		return
	}

	if UID == 0 {
		utils.JsonErrorResponse(c, 400, "用户ID不能为空")
		return
	}
	
	if req.Content == "" {
		utils.JsonErrorResponse(c, 400, "表白内容不能为空")
		return
	}
	if err := services.CreatePost(&models.Post{
		UserID:    UID,
		Content:   req.Content,
		Anonymous: req.Anonymous,
		Visibility: req.Visibility,
	}); err != nil {
		utils.JsonErrorResponse(c, 500, "创建帖子失败")
		return
	}

	utils.JsonSuccessResponse(c, req)
}

// 更新帖子
type UpdatePostData struct {
	ID         uint64  `json:"id" binding:"required"`
	Content    string  `json:"content"`
	Anonymous  *bool   `json:"anonymous"`
	Visibility *bool   `json:"visibility"`
}
func UpdatePost(c *gin.Context) {
	uid, _ := c.Get("userID")
	UID := uid.(uint64)
	var req UpdatePostData
	err := c.ShouldBindJSON(&req); 
	if err != nil {
		utils.JsonErrorResponse(c, 501, "参数错误")
		return
	}
    
	post, err := services.GetPostDataByID(req.ID)
	if err != nil {
		c.Error(errors.New("表白不存在"))
		return
	}
    if UID!=post.UserID {
		utils.JsonErrorResponse(c, 512, "无权限")
		return
	}
	if req.Content == "" {
		req.Content=post.Content
	}
    var anonymous bool
    var visibility bool
	if req.Anonymous == nil {
		anonymous = post.Anonymous
	}else{
		anonymous = *req.Anonymous
	}
	
	if req.Visibility == nil {
		visibility = post.Visibility
	}else{
		visibility = *req.Visibility
	}

	updatedPost := models.Post{
		ID:            req.ID,
		UserID:        UID,
		Content:       req.Content,
		Anonymous:     anonymous,
		Visibility:    visibility,	
		Version:       post.Version,
	}

	err = services.UpdatePost(&updatedPost, post.Version)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Error(errors.New("数据已被其他会话修改，请重试"))
		} else {
			c.Error(errors.New("更新表白信息失败"))
		}
		return
	}

	c.JSON(200, map[string]any{"version": post.Version+1})
}

type DeletePostData struct {
	ID         uint64  `json:"post_id" binding:"required"`
}
func DeletePost(c *gin.Context) {
	uid, _ := c.Get("userID")
	UID := uid.(uint64)
	var req DeletePostData
	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.JsonErrorResponse(c,501, "参数错误")
		return
	}
	originalPost, err := services.GetPostDataByID(req.ID)
	if err != nil {
		utils.JsonErrorResponse(c, 508, "表白不存在")
		return
	}
	if originalPost.UserID != UID {
		utils.JsonErrorResponse(c, 512, "无权限")
		return
	}
	err = services.DeletePost(req.ID)
	if err != nil {
		utils.JsonErrorResponse(c, 511, "删除帖子失败")
		return
	}
	utils.JsonSuccessResponse(c, nil)
}

// GetVisiblePosts 获取未被拉黑的其他人发布的表白
func GetVisiblePosts(c *gin.Context) {
	uid, _ := c.Get("userID")
	UID := uid.(uint64)
	posts,total,err := services.GetVisiblePosts(UID,10,10)
	if err != nil {
		utils.JsonErrorResponse(c, 500, "获取帖子失败")
		return
	}
	log.Println(total)
	utils.JsonSuccessResponse(c, posts)
}