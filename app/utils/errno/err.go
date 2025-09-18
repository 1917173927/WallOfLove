package errno

import "errors"

var (
	ErrImageSizeExceeded = errors.New("图片大小超出限制")
	ErrImageTypeInvalid = errors.New("图片类型无效")
	ErrImageUploadFailed = errors.New("图片上传失败")
	ErrNotImage = errors.New("不是图片")
)