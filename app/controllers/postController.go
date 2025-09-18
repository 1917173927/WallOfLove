package controllers

import (
	"errors"

	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/app/services"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 创建帖子
type PostData struct {
	Content    string `json:"content" binding:"required"`
	Anonymous  bool   `json:"anonymous"`
	Visibility bool   `json:"visibility"`
}

func CreatePost(c *gin.Context) {
	var req PostData
	uid, _ := c.Get("userID")
	UID := uid.(uint64)
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JsonErrorResponse(c, 400, "参数错误")
		return
	}
	user, err := services.GetUserDataByID(UID)
	if err != nil {
		utils.JsonErrorResponse(c, 400, "用户不存在")
		return
	}
	if req.Content == "" {
		utils.JsonErrorResponse(c, 400, "表白内容不能为空")
		return
	}
	if err := services.CreatePost(&models.Post{
		UserID:        UID,
		Content:       req.Content,
		Anonymous:     req.Anonymous,
		Visibility:    req.Visibility,
		UserNickname:  user.Nickname,
		AvatarImageID: user.AvatarImageID,
	}); err != nil {
		utils.JsonErrorResponse(c, 500, "创建帖子失败")
		return
	}

	utils.JsonSuccessResponse(c, req)
}

// 更新帖子
type UpdatePostData struct {
	ID         uint64 `json:"id" binding:"required"`
	Content    string `json:"content"`
	Anonymous  *bool  `json:"anonymous"`//bool值不能为空，只能通过传指针来判断
	Visibility *bool  `json:"visibility"`//同上
}

func UpdatePost(c *gin.Context) {
	uid, _ := c.Get("userID")
	UID := uid.(uint64)
	var req UpdatePostData
	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.JsonErrorResponse(c, 501, "参数错误")
		return
	}
	post, err := services.GetPostDataByID(req.ID)
	if err != nil {
		c.Error(errors.New("表白不存在"))
		return
	}
	// 检查是否有权限修改帖子
	if UID != post.UserID {
		utils.JsonErrorResponse(c, 512, "无权限")
		return
	}
	//若未填写内容，则用原值
	if req.Content == "" {
		req.Content = post.Content
	}
	//转指针
	var anonymous bool
	var visibility bool
	//若未填写匿名，则用原值
	if req.Anonymous == nil {
		anonymous = post.Anonymous
	} else {
		anonymous = *req.Anonymous
	}
	//若未填写可见性，则用原值
	if req.Visibility == nil {
		visibility = post.Visibility
	} else {
		visibility = *req.Visibility
	}

	updatedPost := models.Post{
		ID:            req.ID,
		UserID:        UID,
		Content:       req.Content,
		Anonymous:     anonymous,
		Visibility:    visibility,
		Version:       post.Version,
		UserNickname:  post.UserNickname,
		AvatarImageID: post.AvatarImageID,
	}
	//更新表白信息，乐观锁
	err = services.UpdatePost(&updatedPost, post.Version)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Error(errors.New("数据已被其他会话修改，请重试"))
		} else {
			c.Error(errors.New("更新表白信息失败"))
		}
		return
	}

	c.JSON(200, map[string]any{"version": post.Version + 1})
}

// 删除帖子
type DeletePostData struct {
	ID uint64 `json:"post_id" binding:"required"`
}

func DeletePost(c *gin.Context) {
	uid, _ := c.Get("userID")
	UID := uid.(uint64)
	var req DeletePostData
	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.JsonErrorResponse(c, 501, "参数错误")
		return
	}
	originalPost, err := services.GetPostDataByID(req.ID)
	if err != nil {
		utils.JsonErrorResponse(c, 508, "表白不存在")
		return
	}
	// 检查是否有权限删除帖子
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
type PageData struct {
	PageSize int `json:"page_size"`
	PageNum  int `json:"page_num"`
}

func GetVisiblePosts(c *gin.Context) {
	uid, _ := c.Get("userID")
	UID := uid.(uint64)
	var req PageData
	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.JsonErrorResponse(c, 501, "参数错误")
		return
	}
	posts, total, err := services.GetVisiblePosts(UID, req.PageNum, req.PageSize)
	if err != nil {
		utils.JsonErrorResponse(c, 500, "获取帖子失败")
		return
	}
	//匿名
	for i := range posts {
		if posts[i].Anonymous {
			posts[i].UserID = 0
			posts[i].UserNickname = "?"
			posts[i].AvatarImageID = nil
		}
	}
	data := map[string]any{
		"posts": posts,
		"total": total,
	}
	utils.JsonSuccessResponse(c, data)
}
