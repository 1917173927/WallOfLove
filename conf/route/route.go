package route

import (
	"github.com/1917173927/WallOfLove/app/controllers"
	"github.com/1917173927/WallOfLove/app/middleware"
	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {
	r.MaxMultipartMemory = 64 << 20 // 64MB
	r.Static("/images", "./images")
	// 全局前缀
	api := r.Group("/api")
	{
		// 无需 JWT
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)
		// 需要 JWT
		auth := api.Group("")
		auth.Use(middleware.JWT())
		{
			auth.PUT("/profile", controllers.UpdateProfile)
			auth.GET("/profile", controllers.GetProfile)

			auth.POST("/post", controllers.CreatePost)
			auth.PUT("/post", controllers.UpdatePost)
			auth.GET("/post/list", controllers.GetVisiblePosts)
			auth.GET("/post", controllers.GetSinglePost)
			auth.GET("/post/user", controllers.GetPostsByUserID)
			auth.DELETE("/post", controllers.DeletePost)
			auth.GET("/popranking", controllers.PopRanking)

			auth.POST("/review", controllers.CreateReview)
			auth.GET("/review", controllers.GetReviewsByPostID)
			auth.DELETE("/review", controllers.DeleteReview)
			auth.POST("/reply", controllers.CreateReply)
			auth.GET("/reply", controllers.GetRepliesByReviewID)
			auth.DELETE("/reply", controllers.DeleteReply)

			auth.POST("/blacklist", controllers.BlackUser)
			auth.DELETE("/blacklist", controllers.UnblackUser)
			auth.GET("/blacklist", controllers.GetBlackList)

			auth.POST("/image", controllers.UploadImage)
			auth.DELETE("/image", controllers.DeleteImage)

			auth.POST("/like", controllers.LikePost)
		}
	}
}
