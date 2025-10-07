package services

import (
	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/1917173927/WallOfLove/conf/database"
	"golang.org/x/crypto/bcrypt"
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
	return database.DB.Create(&user).Error
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CompareHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func GetUserDataByID(userID uint64) (*models.User, error) {
	var user models.User
	result := database.DB.
		Select("id", "username", "nickname", "avatar_path", "created_at", "gender", "signature").
		Where("id = ?", userID).
		First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
func GetAllUserDataByID(userID uint64) (*models.User, error) {
	var user models.User
	result := database.DB.
		Where("id = ?", userID).
		First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
func UpdateProfile(user *models.User) error {
	return database.DB.Model(user).
		Select("nickname", "username", "password", "avatar_path", "gender", "signature").
		Updates(user).Error
}

func BlackUser(userID, blockedID uint64) error {
	return database.DB.Create(&models.Blacklist{UserID: userID, BlockedID: blockedID}).Error
}

func UnblackUser(userID, blockedID uint64) error {
	return database.DB.Where("user_id = ? AND blocked_id = ?", userID, blockedID).Delete(&models.Blacklist{}).Error
}

type BlackedUser struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

func GetBlackedUsers(userID uint64) ([]BlackedUser, error) {
	ids, err := utils.GetBlackListIDs(userID)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return []BlackedUser{}, nil
	}
	var list []BlackedUser
	err = database.DB.
		Table("Users").
		Select("id as user_id, username, nickname").
		Where("id IN ?", ids).
		Scan(&list).Error
	return list, err
}
