package services

import (
	"slices"

	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/1917173927/WallOfLove/conf/database"
)

//创建评论
func CreateReview(review *models.Review) error {
	return database.DB.Create(review).Error
}
func GetReviewsByPostID(postID uint64) error {
	var reviews []models.Review
	return database.DB.Where("post_id = ?", postID).Find(&reviews).Error
}
//获取未被拉黑的其他人发布的评论
func GetVisibleReviews(postID uint64,userID uint64, page, pageSize int) ([]models.Review, int64, error) {
	sub, _ := utils.GetBlackListIDs(userID)
	var total int64
	database.DB.Model(&models.Review{}).
		Where("user_id NOT IN (?)", sub).
		Count(&total)

	var list []models.Review
	err := database.DB.Preload("Review2s").
	    Where("post_id = ?", postID).
		Where("user_id NOT IN (?)", sub).
		Order("created_at desc").
		Scopes(utils.Paginate(page, pageSize)).
		Find(&list).Error
		
	for i := range list {
		filterReview2s(&list[i], sub)
	}

	return list, total, err
}
//过滤被拉黑人的回复
func filterReview2s(review *models.Review, blackList []uint64) {
	if len(review.Review2s) == 0 || len(blackList) == 0 {
		return
	}
	// 把不在黑名单里的回复留下
	filtered := make([]models.Review2, 0, len(review.Review2s))
	for _, r2 := range review.Review2s {
		if !slices.Contains(blackList, r2.UserID) {
			filtered = append(filtered, r2)
		}
	}
	review.Review2s = filtered
}