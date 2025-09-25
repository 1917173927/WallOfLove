package services

import (
	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/conf/database"
)

//点赞
func Like(userID, postID,reviewID uint64) error {
    //存在就返回 nil
    var like int64
	database.DB.Model(&models.Like{}).Where("user_id = ? AND post_id = ? AND review_id = ?", userID, postID,reviewID).Count(&like)
    if like > 0 {
        return nil
    }
    //不存在就插入
    return database.DB.Create(&models.Like{UserID: userID, PostID: postID,ReviewID: reviewID}).Error
}
//取消点赞
func Unlike(userID, postID,reviewID uint64) error {
    return database.DB.
        Where("user_id = ? AND post_id = ?AND review_id = ?", userID, postID,reviewID).
        Delete(&models.Like{}).Error
}