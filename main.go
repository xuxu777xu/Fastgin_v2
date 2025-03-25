package main

import (
	"FastGin/core"
	"FastGin/pkg/config"
	"FastGin/pkg/logg"
	"FastGin/router"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	//Todo: service里面的login错误处理   data的值应该改为nil还是*result 等待实践后做具体修改
	//初始化日志系统
	//使用默认配置
	Logconfig := logg.DefaultConfig()

	err := logg.InitLogger(Logconfig)
	if err != nil {
		panic(err)
	}

	//配置文件的读取
	config.RunSettingFile()
	// 读取配置
	cfg := core.ReadConfig(config.Options.File)
	// 这里添加您的应用程序逻辑
	logg.Info("应用程序启动成功，配置已加载")
	logg.Info("数据库配置: ", cfg.DB)

	// 设置 Gin 的运行模式
	gin.SetMode(cfg.Server.Mode)
	// 初始化路由
	r := router.InitRouter()
	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logg.Info("服务器启动在端口", addr)

	if err := r.Run(addr); err != nil {
		logg.Error("服务器启动失败:", err)
		panic(err)
	}
}
