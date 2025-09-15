package models

import "time"

type Image struct {
    ID         uint64     
    OwnerID    uint64            // 上传者
    PostID     *uint64       // 关联 post（若作为 post 图）
    IsAvatar   bool               // 是否头像
    FilePath   string     // 存储相对 path 而非外部 URL
    ThumbPath  string      // 缩略图相对 path
    Mime       string      // e.g. "image/png"
    Width      int
    Height     int
    Size       int64
    Checksum   string      // 去重
    OrderIndex int                     // 在 post 中的顺序（1..9）
    CreatedAt  time.Time
}