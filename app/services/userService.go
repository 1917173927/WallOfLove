package services

import (
	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/conf/database"
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

func BlackUser( userID, blockedID uint64) error {
	return database.DB.Create(&models.Blacklist{UserID: userID, BlockedID: blockedID}).Error
}
func UnblackUser(userID, blockedID uint64) error {
	return database.DB.Where("user_id = ? AND blocked_id = ?", userID, blockedID).Delete(&models.Blacklist{}).Error
}
func GetBlackListID(userID uint64) (myBlack []uint64, err error) {
	err = database.DB.Model(&models.Blacklist{}).Where("user_id = ?", userID).Pluck("blocked_id", &myBlack).Error
	return
}