package utils

import (
	"gorm.io/gorm"
)

// Paginate 分页
func Paginate(page int, size int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		pageSize := size
		if pageSize > 20 {
			pageSize = 20
		} else if pageSize <= 0 {
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
