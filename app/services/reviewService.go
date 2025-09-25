package services

import (
	"slices"

	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/1917173927/WallOfLove/conf/database"
	"gorm.io/gorm"
)

//创建评论
func CreateReview(review *models.Review) error {
	return database.DB.Create(review).Error
}
func GetReviewsByPostID(postID uint64) error {
	var reviews []models.Review
	return database.DB.Where("post_id = ?", postID).Find(&reviews).Error
}
func GetReviewByReviewID(ReviewID uint64) error {
	var reviews []models.Review
	return database.DB.Where("id = ?", ReviewID).First(&reviews).Error
}
//获取未被拉黑的其他人发布的评论,现在是获得所有评论and每条评论的前两条回复
func GetVisibleReviews(reviewID uint64,userID uint64, page, pageSize int) ([]models.Review, int64, error) {
	sub, _ := utils.GetBlackListIDs(userID)
	var total int64
	database.DB.Model(&models.Review{}).
		Where("user_id NOT IN (?)", sub).
		Count(&total)

	var list []models.Review
	err := database.DB.Table("review").
		Select(`review.*,(SELECT COUNT(*) FROM likes WHERE likes.review_id = review.id) AS like_count,
		                  EXISTS(SELECT 1 FROM likes WHERE likes.review_id = review.id AND likes.user_id = ?) AS liked_by_me`,userID).
	    Preload("Replies", func(db *gorm.DB) *gorm.DB {
		    return db.Order("created_at DESC").Limit(2)
	    }).
		Where("id = ?", reviewID).
		Where("user_id NOT IN (?)", sub).
		Order("created_at desc").
		Scopes(utils.Paginate(page, pageSize)).
		Find(&list).Error
	
	for i := range list {
		filterReplies(&list[i], sub)
	}
	return list, total, err
}
//过滤被拉黑人的回复
func filterReplies(review *models.Review, blackList []uint64) {
	if len(review.Replies) == 0 || len(blackList) == 0 {
		return
	}
	// 把不在黑名单里的回复留下
	filtered := make([]models.Reply, 0, len(review.Replies))
	for _, r2 := range review.Replies {
		if !slices.Contains(blackList, r2.UserID) {
			filtered = append(filtered, r2)
		}
	}
	review.Replies = filtered
}

