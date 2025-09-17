package route

import (
	"github.com/gin-gonic/gin"
	"github.com/1917173927/WallOfLove/app/controllers"
	"github.com/1917173927/WallOfLove/app/middleware"
)



func Init(r *gin.Engine) {
	r.MaxMultipartMemory = 64 << 20 // 64MB
	const pre = "/api"
	r.POST(pre+"/register", controllers.Register)
	r.POST(pre+"/login", controllers.Login)
	r.PUT(pre+"/user", middleware.JWT(), controllers.UpdateProfile)
	r.POST(pre+"/post", middleware.JWT(), controllers.CreatePost)
	r.PUT(pre+"/post", middleware.JWT(), controllers.UpdatePost)
	r.DELETE(pre+"/post/:id", middleware.JWT(), controllers.DeletePost)
	r.POST(pre+"/blacklist", middleware.JWT(), controllers.BlackUser)
	r.DELETE(pre+"/blacklist", middleware.JWT(), controllers.UnblackUser)
	r.POST(pre+"/uploadimage", middleware.JWT(), controllers.UploadImage)
}