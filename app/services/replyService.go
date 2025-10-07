package services

import (
	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/1917173927/WallOfLove/conf/database"
	"github.com/1917173927/WallOfLove/conf/redis"
)

// CreateReply 新增回复并增加对应评论的回复计数
func CreateReply(reply *models.Reply) error {
	redis.IncrReply(reply.ReviewID)
	return database.DB.Create(reply).Error
}

// GetRepliesByReviewID 获取评论的回复（过滤黑名单用户），附带用户昵称与头像
func GetRepliesByReviewID(reviewID, userID uint64, page, pageSize int) ([]models.ReplyWithNickname, int64, error) {
	sub, _ := utils.GetBlackListIDs(userID)
	var total int64
	countQ := database.DB.Model(&models.Reply{}).Where("review_id = ?", reviewID)
	if len(sub) > 0 {
		countQ = countQ.Where("user_id NOT IN (?)", sub)
	}
	countQ.Count(&total)

	var replies []models.Reply
	q := database.DB.Where("review_id = ?", reviewID)
	if len(sub) > 0 {
		q = q.Where("user_id NOT IN (?)", sub)
	}
	if err := q.Order("created_at desc").Scopes(utils.Paginate(page, pageSize)).Find(&replies).Error; err != nil {
		return nil, 0, err
	}

	// 批量抓取用户昵称与头像
	userInfos := make(map[uint64]struct {
		Nickname   string
		AvatarPath string
	})
	for _, r := range replies {
		if _, ok := userInfos[r.UserID]; ok { // 已缓存
			continue
		}
		if u, err := GetUserDataByID(r.UserID); err == nil && u != nil {
			userInfos[r.UserID] = struct {
				Nickname   string
				AvatarPath string
			}{Nickname: u.Nickname, AvatarPath: u.AvatarPath}
		}
	}

	list := make([]models.ReplyWithNickname, 0, len(replies))
	for _, r := range replies {
		info := userInfos[r.UserID]
		list = append(list, models.ReplyWithNickname{
			Reply:      r,
			Nickname:   info.Nickname,
			AvatarPath: info.AvatarPath,
		})
	}
	return list, total, nil
}

// DeleteReply 删除回复并减少计数
func DeleteReply(replyID uint64) error {
	redis.DecrReply(replyID)
	return database.DB.Delete(&models.Reply{}, replyID).Error
}

// GetReplyByReplyID 根据ID获取单条回复
func GetReplyByReplyID(replyID uint64) (*models.Reply, error) {
	var reply models.Reply
	if err := database.DB.Where("id = ?", replyID).First(&reply).Error; err != nil {
		return nil, err
	}
	return &reply, nil
}
