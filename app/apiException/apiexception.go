package apiException

import (
	"net/http"

	"github.com/1917173927/WallOfLove/app/utils/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Error 表示自定义错误，包括状态码、消息和日志级别。
type Error struct {
	Code  int
	Msg   string
	Level log.Level
}

// Error 表示自定义的错误类型
var (
	ServerError           = NewError(200500, log.LevelError,"系统异常，请稍后重试!")
	FileNotImageError     = NewError(200502, log.LevelInfo, "上传的文件不是图片")
	ParamError            = NewError(200501, log.LevelInfo, "参数错误")
	UserAlreadyExisted    = NewError(200503, log.LevelInfo, "该用户已被激活")
	TargetError           = NewError(200504, log.LevelInfo, "找不到该目标")
	EmptyError            = NewError(200505, log.LevelInfo, "内容不能为空")
	PwdOrParamError       = NewError(200506, log.LevelInfo, "密码不能小于8位,或参数错误")
	IllegalTarget         = NewError(200507, log.LevelInfo, "不能拉黑自己")
	NoThatPasswordOrWrong = NewError(200508, log.LevelInfo, "密码错误")
	NotLogin              = NewError(200509, log.LevelInfo, "未登录")
	NotPermission         = NewError(200510, log.LevelInfo, "该用户无权限")
	NotFindUser           = NewError(200511, log.LevelInfo, "未注册")
	ConflictError         = NewError(200512, log.LevelInfo, "数据已被其他会话修改，请重试")
	TokenExpired          = NewError(200513, log.LevelInfo, "登录已过期")
	UploadFileError       = NewError(200514, log.LevelInfo, "上传文件失败")
	FileSizeExceedError   = NewError(200515, log.LevelInfo, "文件大小超限")
	ImageFormatError      = NewError(200516, log.LevelInfo, "仅支持 JPG 和 PNG 格式的图片")

	NotFound = NewError(200404, log.LevelWarn, http.StatusText(http.StatusNotFound))
)

// Error 方法实现了 error 接口，返回错误的消息内容
func (e *Error) Error() string {
	return e.Msg
}

// NewError 创建并返回一个新的自定义错误实例
func NewError(code int, level log.Level, msg string) *Error {
	return &Error{
		Code:  code,
		Msg:   msg,
		Level: level,
	}
}

// AbortWithException 用于返回自定义错误信息
func AbortWithException(c *gin.Context, apiError *Error, err error) {
	logError(c, apiError, err)
	_ = c.AbortWithError(200, apiError) //nolint:errcheck
}

// logError 记录错误日志
func logError(c *gin.Context, apiErr *Error, err error) {
	// 构建日志字段
	logFields := []zap.Field{
		zap.Int("error_code", apiErr.Code),
		zap.String("path", c.Request.URL.Path),
		zap.String("method", c.Request.Method),
		zap.String("ip", c.ClientIP()),
		zap.Error(err), // 记录原始错误信息
	}
	log.GetLogFunc(apiErr.Level)(apiErr.Msg, logFields...)
}