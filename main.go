package main

import (
	"github.com/1917173927/WallOfLove/app/validator"
	"github.com/1917173927/WallOfLove/conf/database"
	"github.com/1917173927/WallOfLove/conf/route"
	"github.com/gin-gonic/gin"
)

func main() {
	database.Init()
	validator.Init()

	r := gin.Default()
	route.Init(r)
	r.Run(":8080")
}
