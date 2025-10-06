package services

import (
	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/conf/database"
	"github.com/1917173927/WallOfLove/conf/redis"
)

//点赞
func Like(userID, postID,reviewID uint64) (liked bool, err error) {
	if redis.IsUserLiked(postID, userID,reviewID) {
		// 取消点赞
		redis.DelUserLiked(postID,reviewID,userID)
		redis.DecrPostLike(postID,reviewID)
		liked = false

		//同步删库
		if err:= database.DB.Where("user_id = ? AND post_id = ? AND review_id = ?", userID, postID, reviewID).Delete(&models.Like{}).Error; err!= nil {
			// 库删失败 → 回滚缓存
			redis.SetUserLiked(postID, reviewID, userID)
			redis.IncrPostLike(postID, reviewID)
			return false, err
		}
	} else {
		// 点赞
		redis.SetUserLiked(postID,reviewID,userID)
		redis.IncrPostLike(postID,reviewID)
		liked = true

		//同步写库
		if err:= database.DB.Create(&models.Like{UserID: userID, PostID: postID, ReviewID: reviewID}).Error; err!= nil {
			// 库插失败 → 回滚缓存
			redis.DelUserLiked(postID, reviewID, userID)
			redis.DecrPostLike(postID, reviewID)
			return true, err
		}
	}
	return liked, nil
}
