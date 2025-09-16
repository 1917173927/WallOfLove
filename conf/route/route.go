package route

import (
	"github.com/gin-gonic/gin"
	"github.com/1917173927/WallOfLove/app/controllers"
)



func Init(r *gin.Engine) {
	const pre = "/api"
	r.POST(pre+"/register", controllers.Register)
	r.POST(pre+"/login", controllers.Login)
	r.POST(pre+"/post", controllers.CreatePost)
	r.PUT(pre+"/post", controllers.UpdatePost)
	r.DELETE(pre+"/post/:id", controllers.DeletePost)
	r.POST(pre+"/blacklist", controllers.BlackUser)
	r.DELETE(pre+"/blacklist", controllers.UnblackUser)
}