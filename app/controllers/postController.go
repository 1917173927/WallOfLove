package controllers

import (
	"errors"

	"github.com/1917173927/WallOfLove/app/apiException"
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
		apiException.AbortWithException(c,apiException.ParamError,err)
		return
	}
	user, err := services.GetUserDataByID(UID)
	if err != nil {
		apiException.AbortWithException(c,apiException.ServerError,err)
		return
	}
	if req.Content == "" {
		apiException.AbortWithException(c,apiException.EmptyError,nil)
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
		apiException.AbortWithException(c,apiException.ServerError,err)
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
		apiException.AbortWithException(c,apiException.ParamError,err)
		return
	}
	post, err := services.GetPostDataByID(req.ID)
	if err != nil {
		apiException.AbortWithException(c,apiException.TargetError,err)
		return
	}
	// 检查是否有权限修改帖子
	if UID != post.UserID {
		apiException.AbortWithException(c,apiException.NotPermission,nil)
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
			apiException.AbortWithException(c,apiException.ConflictError,err)
			return
		} else {
			apiException.AbortWithException(c,apiException.ServerError,err)
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
		apiException.AbortWithException(c,apiException.ParamError,err)
		return
	}
	originalPost, err := services.GetPostDataByID(req.ID)
	if err != nil {
		apiException.AbortWithException(c,apiException.TargetError,err)
		return
	}
	// 检查是否有权限删除帖子
	if originalPost.UserID != UID {
		apiException.AbortWithException(c,apiException.NotPermission,nil)
		return
	}
	err = services.DeletePost(req.ID)
	if err != nil {
		apiException.AbortWithException(c,apiException.ServerError,err)
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
		apiException.AbortWithException(c,apiException.ParamError,err)
		return
	}
	posts, total, err := services.GetVisiblePosts(UID, req.PageNum, req.PageSize)
	if err != nil {
		apiException.AbortWithException(c,apiException.ServerError,err)
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
