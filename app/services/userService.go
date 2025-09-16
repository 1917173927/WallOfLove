package services

import (
	"fmt"
	"log"

	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/conf/database"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CheckUsername(username string) error {
	result := database.DB.Where("username=?", username).First(&models.User{})
	return result.Error
}
func GetUser(username string) (*models.User, error) {
	var user models.User
	result := database.DB.Where("username=?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
func Register(user models.User) error {
	result := database.DB.Create(&user)
	return result.Error
}
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
func CompareHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}
func GetUserDataByID(userID uint64) (*models.User, error) {
	var user models.User
	result := database.DB.
		Where("id = ?", userID).
		First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func UpdateProfile(user *models.User, oldVersion uint) error {
	tx := database.DB.Model(&models.User{}).
		Where("id = ? AND version = ?", user.ID, oldVersion).
		Updates(map[string]any{
			"nickname":        user.Nickname,
			"password":        user.Password,
			"avatar_image_id": user.AvatarImageID,
			"version":         gorm.Expr("version + 1"),
		})
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return tx.Error
}

func BlackUser(db *gorm.DB, userID, blockedID uint64) error {
	if userID == blockedID {
		return gorm.ErrInvalidData
	}
	return db.Create(&models.Blacklist{UserID: userID, BlockedID: blockedID}).Error
}
func UnblackUser(db *gorm.DB, userID, blockedID uint64) error {
	return db.Where("user_id = ? AND blocked_id = ?", userID, blockedID).Delete(&models.Blacklist{}).Error
}
func GetBlackListID(db *gorm.DB, userID uint64) (myBlack []uint64, blackMe []uint64, err error) {
	err = db.Model(&models.Blacklist{}).Where("user_id = ?", userID).Pluck("blocked_id", &myBlack).Error
	if err != nil {
		log.Println("获取拉黑列表失败:", err)
		return
	}
	err = db.Model(&models.Blacklist{}).Where("blocked_id = ?", userID).Pluck("user_id", &blackMe).Error
	if err != nil {
		log.Println("获取被拉黑列表失败:", err)
		return
	}
	return
}
func FilterBlack(c *gin.Context,db *gorm.DB, userID uint64) *gorm.DB {
	key := fmt.Sprintf("black:%d", userID)
	// 看本请求有没有缓存
	if v, ok := c.Get(key); ok {
		all := v.([]uint64)
		if len(all) == 0 {
			return db.Where("1=1")
		}
		return db.Where("user_id NOT IN ?", all)
	}
	// 没有库
	myBlack, blackMe, _ := GetBlackListID(db, userID)
	all := append(myBlack, blackMe...)
	c.Set(key, all) // 缓存
	if len(all) == 0 {
		return db.Where("1=1")
	}
	return db.Where("user_id NOT IN ?", all)
}
	//在获取时使用黑名单过滤加入：filter:=services.FilterBlack(database.DB,UID)  后面用filter.Find(&posts)或者其他即可
