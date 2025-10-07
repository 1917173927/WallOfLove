package services

import (
	"slices"

	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/1917173927/WallOfLove/conf/database"
	"github.com/1917173927/WallOfLove/conf/redis"
	"gorm.io/gorm"
)

func CreateReview(review *models.Review) error {
	redis.IncrReview(review.PostID)
	return database.DB.Create(review).Error
}

func DeleteReview(reviewID uint64) error {
	redis.DecrReview(reviewID)
	return database.DB.Delete(&models.Review{}, "id = ?", reviewID).Error
}
func GetReviewsByPostID(postID uint64) error {
	var reviews []models.Review
	return database.DB.Where("post_id = ?", postID).Find(&reviews).Error
}
func GetReviewByReviewID(ReviewID uint64) (*models.Review, error) {
	var review models.Review
	result := database.DB.
		Where("id = ?", ReviewID).
		First(&review)
	if result.Error != nil {
		return nil, result.Error
	}
	return &review, nil
}

// 获得所有评论and每条评论的前两条回复
type ReviewWithLike struct {
	models.Review
	LikeCount    int64                     `json:"like_count"`
	LikedByMe    bool                      `json:"liked_by_me"`
	Replies      []models.ReplyWithNickname `json:"replies"`
	RepliesCount int64               `json:"replies_count"`
	Nickname     string              `json:"nickname"`
	AvatarPath   string              `json:"avatarPath"`
}

func GetVisibleReviews(postID uint64, userID uint64, page, pageSize int) ([]ReviewWithLike, int64, error) {
	sub, _ := utils.GetBlackListIDs(userID)

	var total int64
	countDB := database.DB.Model(&models.Review{}).
		Where("post_id = ?", postID)
	if len(sub) > 0 {
		countDB = countDB.Where("user_id NOT IN (?)", sub)
	}
	countDB.Count(&total)

	var reviews []models.Review
	base := database.DB.
		Preload("Replies", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Limit(2) 
		}).
		Where("post_id = ?", postID)
	if len(sub) > 0 {
		base = base.Where("user_id NOT IN (?)", sub)
	}
	err := base.
		Order("created_at desc").
		Scopes(utils.Paginate(page, pageSize)).
		Find(&reviews).Error
	if err != nil {
		return nil, 0, err
	}

	userIDs := make([]uint64, 0, len(reviews))
	for _, r := range reviews {
		userIDs = append(userIDs, r.UserID)
	}

	userInfos := make(map[uint64]struct {
		Nickname   string
		AvatarPath string
	})
	for _, id := range userIDs {
		user, err := GetUserDataByID(id)
		if err == nil && user != nil {
			userInfos[id] = struct {
				Nickname   string
				AvatarPath string
			}{
				Nickname:   user.Nickname,
				AvatarPath: user.AvatarPath,
			}
		}
	}

	list := make([]ReviewWithLike, 0, len(reviews))
	for _, r := range reviews {
		likeCount := redis.GetPostLikeCount(postID, r.ID)    // 评论点赞 reviewID!=0
		likedByMe := redis.IsUserLiked(postID, userID, r.ID) // 当前用户是否点赞这条评论
		repliesCount := redis.GetReviewReplyCount(r.ID)      //评论回复数

		filterReplies(&r, sub)


		replyList := make([]models.ReplyWithNickname, 0, len(r.Replies))
		for _, reply := range r.Replies {
			replyList = append(replyList, models.ReplyWithNickname{
				Reply:    reply,
				Nickname: userInfos[reply.UserID].Nickname,
			})
		}

		list = append(list, ReviewWithLike{
			Review:       r,
			LikeCount:    likeCount,
			LikedByMe:    likedByMe,
			RepliesCount: repliesCount,
			Replies:      replyList,
			Nickname:     userInfos[r.UserID].Nickname,
			AvatarPath:   userInfos[r.UserID].AvatarPath,
		})
	}

	return list, total, nil
}

// 过滤被拉黑人的回复
func filterReplies(review *models.Review, blackList []uint64) {
	if len(review.Replies) == 0 || len(blackList) == 0 {
		return
	}

	filtered := make([]models.Reply, 0, len(review.Replies))
	for _, r2 := range review.Replies {
		if !slices.Contains(blackList, r2.UserID) {
			filtered = append(filtered, r2)
		}
	}
	review.Replies = filtered
}
