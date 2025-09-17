package route

import (
	"github.com/1917173927/WallOfLove/app/controllers"
	"github.com/1917173927/WallOfLove/app/middleware"
	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {
	r.MaxMultipartMemory = 64 << 20 // 64MB
	const pre = "/api"
	r.POST(pre+"/register", controllers.Register)
	r.POST(pre+"/login", controllers.Login)
	r.PUT(pre+"/user", middleware.JWT(), controllers.UpdateProfile)
	r.POST(pre+"/post", middleware.JWT(), controllers.CreatePost)
	r.PUT(pre+"/post", middleware.JWT(), controllers.UpdatePost)
	r.POST(pre+"/review", middleware.JWT(), controllers.CreateReview)
	r.GET(pre+"/review/:id", middleware.JWT(), controllers.GetReviewsByPostID)
	r.POST(pre+"/review2", middleware.JWT(), controllers.CreateReview2)
	r.GET(pre+"/post/:id", middleware.JWT(), controllers.GetVisiblePosts)
	r.DELETE(pre+"/post/:id", middleware.JWT(), controllers.DeletePost)
	r.POST(pre+"/blacklist", middleware.JWT(), controllers.BlackUser)
	r.DELETE(pre+"/blacklist", middleware.JWT(), controllers.UnblackUser)
	r.POST(pre+"/uploadimage", middleware.JWT(), controllers.UploadImage)
}
