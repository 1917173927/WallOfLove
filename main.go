package main

import (
	"fmt"

	"github.com/1917173927/WallOfLove/app/middleware"
	"github.com/1917173927/WallOfLove/app/services"
	"github.com/1917173927/WallOfLove/conf/database"
	"github.com/1917173927/WallOfLove/conf/redis"
	"github.com/1917173927/WallOfLove/conf/route"
	"github.com/gin-gonic/gin"
)

func main() {
	database.Init()
	if err := redis.Init(); err != nil {
		fmt.Printf("Redis initialization failed: %v\n", err)
		fmt.Println("Continuing without Redis support")
	}
	services.StartScheduler()
	r := gin.Default()
	r.Use(middleware.ErrHandler())
	r.NoMethod(middleware.HandleNotFound)
	r.NoRoute(middleware.HandleNotFound)
	route.Init(r)
	r.Run(":8080")
}
