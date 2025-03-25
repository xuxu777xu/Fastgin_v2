package router

import (
	//"FastGin/api/example"
	"FastGin/api/test"
	"FastGin/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	// 使用中间件
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())        // 添加请求ID中间件
	r.Use(middleware.LoggerMiddleware()) // 添加日志中间件
	r.Use(middleware.Cors())
	//限流操作
	r.Use(middleware.RateLimit(300, 500)) // 每秒最多处理x个请求，突发最大y个

	// API 路由组
	apiGroup := r.Group("/api")
	// 创建处理器实例
	testHandler := test.NewHandler()
	{
		// test service 相关路由
		apiGroup.POST("/login", testHandler.Login) // 登录接口
	}

	return r
}
