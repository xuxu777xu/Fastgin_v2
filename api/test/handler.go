package test

import (
	"FastGin/api/common"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 处理登录请求
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(fmt.Errorf("参数解析失败: %v", err))
		c.JSON(http.StatusBadRequest, common.Response{
			Code:    400,
			Message: "参数错误",
			Data:    err.Error(),
		})
		return
	}

	params := map[string]string{
		"username": req.Username,
		"password": req.Password,
		"uid":      req.Username, // 使用username作为uid传递给服务
	}

	result, err := h.service.Login(params)
	if err != nil {
		c.Error(fmt.Errorf("登录服务调用失败: %v", err))
		c.JSON(http.StatusInternalServerError, common.Response{
			Code:    500,
			Message: "登录失败",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, *result)
}
