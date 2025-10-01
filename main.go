// Package main 是应用程序的入口文件，负责初始化核心组件并启动服务。
package main

import (
	"github.com/1917173927/WallOfLove/app/middleware"
	"github.com/1917173927/WallOfLove/app/services"
	"github.com/1917173927/WallOfLove/conf/database"
	"github.com/1917173927/WallOfLove/conf/route"
	"github.com/gin-gonic/gin"
)

// main 是应用程序的入口函数，执行以下操作：
// 1. 初始化数据库连接
// 2. 初始化验证器
// 3. 创建 Gin 引擎实例并配置中间件
// 4. 注册路由
// 5. 启动服务，监听 8080 端口
func main() {
	database.Init()  // 初始化数据库连接
    services.StartScheduler() // 启动定时任务服务

	r := gin.Default() // 创建 Gin 引擎实例
	r.Use(middleware.ErrHandler()) // 注册全局错误处理中间件
	r.NoMethod(middleware.HandleNotFound) // 处理未实现的 HTTP 方法
    r.NoRoute(middleware.HandleNotFound) // 处理未定义的路由
	route.Init(r) // 注册所有路由
	r.Run(":8080") // 启动服务，监听 8080 端口
}
