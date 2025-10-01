package services

import (
	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/1917173927/WallOfLove/conf/database"
	"github.com/1917173927/WallOfLove/conf/redis"
)

// 创建帖子
func CreatePost(post *models.Post) error {
	// 如果未设置发布时间，则立即发布
	if post.ScheduledAt == nil {
		post.IsPublished = true
	}
	return database.DB.Create(post).Error
}
func GetPostDataByID(postID uint64) (*models.Post, error) {
	var post models.Post
	result := database.DB.
		Where("id = ?", postID).
		First(&post)
	if result.Error != nil {
		return nil, result.Error
	}
	return &post, nil
}
func UpdatePost(post *models.Post) error {
	return database.DB.Model(post).
		Select("content", "anonymous", "visibility").
		Updates(post).Error
}

func DeletePost(postID uint64) error {
	return database.DB.Delete(&models.Post{}, "id = ?", postID).Error
}
//获取未被拉黑的其他人发布的表白
type PostWithLike struct {
	models.Post
	IsFull bool `json:"is_full"` //true:超过100字，false:未超过100字
	ShortContent string `json:"short_content"`
	LikeCount int64 `json:"like_count"`
	LikedByMe bool  `json:"liked_by_me"`
	ReviewsCount int64 `json:"reviews_count"`
}
func GetVisiblePosts(userID uint64, page, pageSize int) ([]PostWithLike, int64, error) {
	sub, _ := utils.GetBlackListIDs(userID)

	//总条数
	var total int64
	database.DB.Model(&models.Post{}).
		Where("visibility = ?", true).
		Where("user_id NOT IN (?)", sub).
		Count(&total)

	//拿帖子
	var posts []models.Post
	err := database.DB.
		Preload("Images").
		Where("visibility = ?", true).
		Where("user_id NOT IN (?)", sub).
		Order("created_at desc").
		Scopes(utils.Paginate(page, pageSize)).
		Find(&posts).Error
	if err != nil {
		return nil, 0, err
	}

	//点赞数 + 是否已赞（redis）
	list := make([]PostWithLike, 0, len(posts))
	for _, p := range posts {
		likeCount := redis.GetPostLikeCount(p.ID, 0)     // 帖子点赞 reviewID=0
		likedByMe := redis.IsUserLiked(p.ID, userID, 0) // 当前用户是否点赞
		reviewsCount := redis.GetPostReviewCount(p.ID) // 帖子评论数

		short:=p.Content
		isFull := false
	if len([]rune(p.Content)) > 100 {
		short = string([]rune(p.Content)[:100]) + "..."
		isFull = true
	}

		list = append(list, PostWithLike{
			Post:      p,
			IsFull: isFull,
			ShortContent: short,
			LikeCount: likeCount,
			LikedByMe: likedByMe,
			ReviewsCount: reviewsCount,
		})
	}

	return list, total, nil
}
//根据用户id获取用户发布的表白
func GetPostsByUserID(userID uint64, page, pageSize int) ([]PostWithLike, int64, error) {	
	//总条数
	var total int64
	database.DB.Model(&models.Post{}).
		Where("visibility = ?", true).
		Where("anonymous = ?",false).
		Where("user_id = ?", userID).
		Count(&total)
	//拿帖子
	var posts []models.Post
	err := database.DB.
		Preload("Images").
		Where("visibility = ?", true).
		Where("anonymous = ?",false).
		Where("user_id = ?", userID).
		Order("created_at desc").
		Scopes(utils.Paginate(page, pageSize)).
		Find(&posts).Error
	if err != nil {
		return nil, 0, err
	}

	//点赞数 + 是否已赞（redis）
	list := make([]PostWithLike, 0, len(posts))
	for _, p := range posts {
		likeCount := redis.GetPostLikeCount(p.ID, 0)     // 帖子点赞 reviewID=0
		likedByMe := redis.IsUserLiked(p.ID, userID, 0) // 当前用户是否点赞
		reviewsCount := redis.GetPostReviewCount(p.ID) // 帖子评论数	
		
		short:=p.Content
		isFull := false
	if len([]rune(p.Content)) > 100 {
		short = string([]rune(p.Content)[:100]) + "..."
		isFull = true
	}

		list = append(list, PostWithLike{
			Post:      p,
			IsFull: isFull,
			ShortContent: short,
			LikeCount: likeCount,
			LikedByMe: likedByMe,
			ReviewsCount: reviewsCount,
		})
	}

	return list, total, nil
}

//根据用户id获取用户发布的表白
func GetMyPosts(userID uint64, page, pageSize int) ([]PostWithLike, int64, error) {	
	//总条数
	var total int64
	database.DB.Model(&models.Post{}).
		Where("user_id = ?", userID).
		Count(&total)
	//拿帖子
	var posts []models.Post
	err := database.DB.
		Preload("Images").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Scopes(utils.Paginate(page, pageSize)).
		Find(&posts).Error
	if err != nil {
		return nil, 0, err
	}

	//点赞数 + 是否已赞（redis）
	list := make([]PostWithLike, 0, len(posts))
	for _, p := range posts {
		likeCount := redis.GetPostLikeCount(p.ID, 0)     // 帖子点赞 reviewID=0
		likedByMe := redis.IsUserLiked(p.ID, userID, 0) // 当前用户是否点赞
		reviewsCount := redis.GetPostReviewCount(p.ID) // 帖子评论数	
		
		short:=p.Content
		isFull := false
	if len([]rune(p.Content)) > 100 {
		short = string([]rune(p.Content)[:100]) + "..."
		isFull = true
	}

		list = append(list, PostWithLike{
			Post:      p,
			IsFull: isFull,
			ShortContent: short,
			LikeCount: likeCount,
			LikedByMe: likedByMe,
			ReviewsCount: reviewsCount,
		})
	}

	return list, total, nil
}
//获取单个表白
type SinglePost struct {
	models.Post
	LikeCount int64 `json:"like_count"`
	LikedByMe bool  `json:"liked_by_me"`
	ReviewsCount int64 `json:"reviews_count"`
}
func GetSinglePost(postID,userID uint64) (SinglePost, error) {
	var post models.Post
	redis.IncrView(postID)
	err := database.DB.
        Create(&models.View{PostID: postID, UserID: userID}).Error
	if err != nil {
		return SinglePost{},err
	}
	result := database.DB.
		Preload("Images").
		Where("id = ?", postID).
		First(&post)
	if result.Error != nil {
		return SinglePost{},result.Error
	}
	likeCount := redis.GetPostLikeCount(post.ID, 0) 
	likedByMe := redis.IsUserLiked(post.ID, userID, 0) 
	reviewsCount := redis.GetPostReviewCount(post.ID) 

	list:=SinglePost{
		Post: post,
		LikeCount: likeCount,
		LikedByMe: likedByMe,
		ReviewsCount: reviewsCount,
	}
	return list,nil
}