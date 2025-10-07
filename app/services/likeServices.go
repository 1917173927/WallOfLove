package services

import (
	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/conf/database"
	"github.com/1917173927/WallOfLove/conf/redis"
)

func Like(userID, postID, reviewID uint64) (bool, error) {
	if redis.IsUserLiked(postID, userID, reviewID) {
		redis.DelUserLiked(postID, reviewID, userID)
		redis.DecrPostLike(postID, reviewID)
		if err := database.DB.Where("user_id = ? AND post_id = ? AND review_id = ?", userID, postID, reviewID).Delete(&models.Like{}).Error; err != nil {
			redis.SetUserLiked(postID, reviewID, userID)
			redis.IncrPostLike(postID, reviewID)
			return false, err
		}
		return false, nil
	}
	redis.SetUserLiked(postID, reviewID, userID)
	redis.IncrPostLike(postID, reviewID)
	if err := database.DB.Create(&models.Like{UserID: userID, PostID: postID, ReviewID: reviewID}).Error; err != nil {
		redis.DelUserLiked(postID, reviewID, userID)
		redis.DecrPostLike(postID, reviewID)
		return true, err
	}
	return true, nil
}
