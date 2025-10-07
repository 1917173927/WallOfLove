package middleware

import (
	"errors"
	"net/http"

	"github.com/1917173927/WallOfLove/app/apiException"
	"github.com/1917173927/WallOfLove/app/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ErrHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			if err != nil {
				var apiErr *apiException.Error
				if !errors.As(err, &apiErr) {
					apiErr = apiException.ServerError
					zap.L().Error("Unknown Error Occurred", zap.Error(err))
				}
				utils.JsonErrorResponse(c, apiErr.Code, apiErr.Msg)
				return
			}
		}
	}
}

func HandleNotFound(c *gin.Context) {
	err := apiException.NotFound
	zap.L().Warn("404 Not Found",
		zap.String("path", c.Request.URL.Path),
		zap.String("method", c.Request.Method),
	)
	utils.JsonResponse(c, http.StatusNotFound, err.Code, err.Msg, nil)
}
