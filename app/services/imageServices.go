package services

import (
	"errors"
	"fmt"
	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/conf/config"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"
	"os"
	"math/rand"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

//上传图片
func UploadImage(db *gorm.DB, c *gin.Context, userID uint64, username string, postID string, isAvatar string, file *multipart.FileHeader) (*models.Image, error) {
	// 检查图片大小
	maxSize := config.Config.GetInt64("image.max_size")
	if maxSize <= 0 {
		maxSize = 2 * 1024 * 1024 // 默认 2MB
	}
	if file.Size > maxSize {
		return nil, errors.New("图片大小超过限制")
	}
	
	// 检查文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".png" {
		return nil, errors.New("仅支持 JPG 和 PNG 格式的图片")
	}

	// 验证文件内容类型
	fileContent, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	defer fileContent.Close()

	buffer := make([]byte, 512)
	if _, err := fileContent.Read(buffer); err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}
	contentType := http.DetectContentType(buffer)
	if !strings.HasPrefix(contentType, "image/") {
		return nil, errors.New("文件内容不是有效的图片")
	}

	// 创建用户专属文件夹
	userFolder := fmt.Sprintf("%d-%s", userID, username)
	userFolderPath := filepath.Join("images", userFolder)
	if err := os.MkdirAll(userFolderPath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("创建用户文件夹失败: %v", err)
	}

	// 生成唯一文件名
	rand.Seed(time.Now().UnixNano())
	fileName := time.Now().Format("20060102150405") + "_" + strings.ToLower(randomString(8)) + ext
	dst := filepath.Join(userFolderPath, fileName)

	// 保存图片到服务器
	if err := c.SaveUploadedFile(file, dst); err != nil {
		return nil, fmt.Errorf("保存文件失败: %v", err)
	}

	// 保存图片信息到数据库
	var postIDUint *uint64
	if postID != "" {
		postIDParsed, err := strconv.ParseUint(postID, 10, 64)
		if err == nil {
			postIDUint = &postIDParsed
		}
	}
	image := &models.Image{
		UserID:   userID,
		PostID:   postIDUint,
		FilePath: dst,
		Size:     file.Size,
		CreatedAt: time.Now(),
	}
	if err := db.Create(image).Error; err != nil {
		return nil, fmt.Errorf("保存图片信息失败: %v", err)
	}

	return image, nil
}

// randomString 生成随机字符串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}