package example

import (
	"FastGin/pkg/errcode"
	"FastGin/pkg/logg"
	"fmt"

	"github.com/gin-gonic/gin"
)

// GetUser 获取用户信息
func (h *Handler) GetUser(c *gin.Context) {
	// 获取URL参数
	userID := c.Query("id")

	// 记录请求参数
	logg.Info("获取用户信息请求", map[string]interface{}{
		"user_id": userID,
		"path":    c.Request.URL.Path,
		"method":  c.Request.Method,
	})

	if userID == "" {
		logg.Warn("获取用户信息失败：缺少用户ID", map[string]interface{}{
			"client_ip": c.ClientIP(),
			"headers":   c.Request.Header,
		})
		errcode.ParamError(c)
		return
	}

	// 模拟获取用户信息
	user := UserResponse{
		ID:       1,
		Username: "张三",
		Age:      25,
	}

	logg.Info("获取用户信息成功", map[string]interface{}{
		"user_id":   userID,
		"user_info": user,
	})
	errcode.Success(c, user)
}

// CreateUser 创建用户
func (h *Handler) CreateUser(c *gin.Context) {
	var req UserRequest

	// 记录开始处理创建用户请求
	logg.Info("开始处理创建用户请求", map[string]interface{}{
		"path":      c.Request.URL.Path,
		"method":    c.Request.Method,
		"client_ip": c.ClientIP(),
	})

	// 绑定JSON请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		logg.Error("创建用户失败：参数验证错误", map[string]interface{}{
			"error":     err.Error(),
			"client_ip": c.ClientIP(),
			"headers":   c.Request.Header,
		})
		errcode.ParamError(c)
		return
	}

	// 参数验证
	if req.Age < 0 || req.Age > 150 {
		errMsg := fmt.Sprintf("创建用户失败：年龄 %d 超出有效范围(0-150)", req.Age)
		logg.Warn(errMsg, map[string]interface{}{
			"request_body": req,
			"client_ip":    c.ClientIP(),
		})
		errcode.Error(c, 400, errMsg)
		return
	}

	// 模拟创建用户
	user := UserResponse{
		ID:       1,
		Username: req.Username,
		Age:      req.Age,
	}

	logg.Info("创建用户成功", map[string]interface{}{
		"request":  req,
		"response": user,
	})
	errcode.Success(c, user)
}
