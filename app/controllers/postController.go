package controllers

import (
	"time"

	"github.com/1917173927/WallOfLove/app/apiException"
	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/app/services"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/gin-gonic/gin"
)

// 创建帖子
type PostData struct {
	Content     string     `json:"content" binding:"required"`
	Anonymous   bool       `json:"anonymous"`
	Visibility  bool       `json:"visibility"`
	ScheduledAt *time.Time `json:"scheduled_at"`
}

func CreatePost(c *gin.Context) {
	var req PostData

	uid, _ := c.Get("userID")
	userID := uid.(uint64)
	if err := c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	user, err := services.GetUserDataByID(userID)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	if req.Content == "" {
		apiException.AbortWithException(c, apiException.EmptyError, nil)
		return
	}
	post := &models.Post{
		UserID:       userID,
		Content:      req.Content,
		Anonymous:    req.Anonymous,
		Visibility:   req.Visibility,
		UserNickname: user.Nickname,
		AvatarPath:   user.AvatarPath,
		ScheduledAt:  req.ScheduledAt,
	}
	if err := services.CreatePost(post); err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	utils.JsonSuccessResponse(c, gin.H{"post_id": post.ID})
}

// 更新帖子
type UpdatePostData struct {
	ID          uint64 `json:"id" binding:"required"`
	Content     string `json:"content"`
	Anonymous   *bool  `json:"anonymous"`  //bool值不能为空，只能通过传指针来判断
	Visibility  *bool  `json:"visibility"` //同上
	IsPublished *bool  `json:"is_published"`
}

func UpdatePost(c *gin.Context) {
	uid, _ := c.Get("userID")
	userID := uid.(uint64)
	var req UpdatePostData
	if err := c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	post, err := services.GetPostDataByID(req.ID)
	if err != nil {
		apiException.AbortWithException(c, apiException.TargetError, err)
		return
	}
	if userID != post.UserID {
		apiException.AbortWithException(c, apiException.NotPermission, nil)
		return
	}
	content := req.Content
	if content == "" {
		content = post.Content
	}
	anonymous := post.Anonymous
	if req.Anonymous != nil {
		anonymous = *req.Anonymous
	}
	visibility := post.Visibility
	if req.Visibility != nil {
		visibility = *req.Visibility
	}
	updatedPost := models.Post{
		ID:           req.ID,
		UserID:       userID,
		Content:      content,
		Anonymous:    anonymous,
		Visibility:   visibility,
		UserNickname: post.UserNickname,
		AvatarPath:   post.AvatarPath,
		IsPublished:  true,
	}
	if err := services.UpdatePost(&updatedPost); err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	utils.JsonSuccessResponse(c, nil)
}

// 删除帖子
type DeletePostData struct {
	ID uint64 `json:"post_id" binding:"required"`
}

func DeletePost(c *gin.Context) {
	uid, _ := c.Get("userID")
	userID := uid.(uint64)
	var req DeletePostData
	if err := c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	post, err := services.GetPostDataByID(req.ID)
	if err != nil {
		apiException.AbortWithException(c, apiException.TargetError, err)
		return
	}
	if post.UserID != userID {
		apiException.AbortWithException(c, apiException.NotPermission, nil)
		return
	}
	if err := services.DeletePost(req.ID); err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	utils.JsonSuccessResponse(c, nil)
}

// 获取未被拉黑的其他人发布的表白
type PageData struct {
	PageSize int `form:"page_size" json:"page_size"`
	PageNum  int `form:"page" json:"page"`
}
type PostWithPaths struct {
	Post       services.PostWithLike
	ImagePaths []string `json:"image_paths"`
	Total      int64    `json:"total"`
}

func GetVisiblePosts(c *gin.Context) {
	uid, _ := c.Get("userID")
	userID := uid.(uint64)
	var req PageData
	if err := c.ShouldBindQuery(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	posts, total, err := services.GetVisiblePosts(userID, req.PageNum, req.PageSize)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	for i := range posts {
		if posts[i].Anonymous {
			posts[i].UserID = 0
			posts[i].UserNickname = "匿名用户"
			posts[i].AvatarPath = "images/default/anonymous.png"
		}
	}
	out := make([]PostWithPaths, 0, len(posts))
	for _, p := range posts {
		paths := make([]string, 0, len(p.Images))
		for _, img := range p.Images {
			paths = append(paths, img.FilePath)
		}
		out = append(out, PostWithPaths{
			Post:       p,
			ImagePaths: paths,
			Total:      total,
		})
	}
	utils.JsonSuccessResponse(c, out)
}

// 获取指定用户发布的表白
type GetPostsByUserIDData struct {
	UserID   uint64 `form:"user_id" binding:"required"`
	PageSize int    `form:"page_size"`
	PageNum  int    `form:"page_num"`
}

func GetPostsByUserID(c *gin.Context) {
	uid, _ := c.Get("userID")
	userID := uid.(uint64)
	var req GetPostsByUserIDData
	if err := c.ShouldBind(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if _, err := services.GetUserDataByID(req.UserID); err != nil {
		apiException.AbortWithException(c, apiException.TargetError, err)
		return
	}
	var posts []services.PostWithLike
	var total int64
	if req.UserID == userID {
		var err error
		posts, total, err = services.GetMyPosts(req.UserID, req.PageNum, req.PageSize)
		if err != nil {
			apiException.AbortWithException(c, apiException.ServerError, err)
			return
		}
	} else {
		var err error
		posts, total, err = services.GetPostsByUserID(req.UserID, req.PageNum, req.PageSize)
		if err != nil {
			apiException.AbortWithException(c, apiException.ServerError, err)
			return
		}
	}
	out := make([]PostWithPaths, 0, len(posts))
	for _, p := range posts {
		paths := make([]string, 0, len(p.Images))
		for _, img := range p.Images {
			paths = append(paths, img.FilePath)
		}
		out = append(out, PostWithPaths{
			Post:       p,
			ImagePaths: paths,
			Total:      total,
		})
	}
	utils.JsonSuccessResponse(c, out)
}

// 获得单个帖子
type GetSinglePostData struct {
	ID uint64 `form:"post_id" binding:"required"`
}
type SinglePost struct {
	Post       services.SinglePost `json:"post"`
	ImagePaths []string            `json:"image_paths"`
}

func GetSinglePost(c *gin.Context) {
	uid, _ := c.Get("userID")
	userID := uid.(uint64)
	var req GetSinglePostData
	if err := c.ShouldBind(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	post, err := services.GetSinglePost(req.ID, userID)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	if post.Anonymous {
		post.UserID = 0
		post.UserNickname = "匿名用户"
		post.AvatarPath = "images/default/anonymous.png"
	}
	paths := make([]string, 0, len(post.Images))
	for _, img := range post.Images {
		paths = append(paths, img.FilePath)
	}
	out := SinglePost{
		Post:       post,
		ImagePaths: paths,
	}
	utils.JsonSuccessResponse(c, out)
}

type ConfirmPostData struct {
	PostID      uint64 `form:"post_id" binding:"required"`
	IsPublished bool   `form:"is_published" binding:"required"`
}
