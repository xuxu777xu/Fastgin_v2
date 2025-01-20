package test

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Login 处理登录请求
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(fmt.Errorf("参数解析失败: %v", err))
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "参数错误",
			Data:    nil,
		})
		return
	}

	params := map[string]string{
		"username": req.Username,
		"password": req.Password,
	}

	if err := h.service.Login(params); err != nil {
		c.Error(fmt.Errorf("登录服务调用失败: %v", err))
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "登录失败",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "登录成功",
		Data:    nil,
	})
}

// GetUnpaidOrders 处理获取未支付订单请求
func (h *Handler) GetUnpaidOrders(c *gin.Context) {
	var req UnpaidOrdersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(fmt.Errorf("参数解析失败: %v", err))
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "参数错误",
			Data:    nil,
		})
		return
	}

	params := map[string]string{
		"user_id": req.UserID,
		"page":    strconv.Itoa(req.Page),
		"size":    strconv.Itoa(req.Size),
	}

	if err := h.service.GetUnpaidOrders(params); err != nil {
		c.Error(fmt.Errorf("获取未支付订单失败: %v", err))
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "获取未支付订单失败",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "获取成功",
		Data:    nil,
	})
}

// ProcessPayment 处理支付请求
func (h *Handler) ProcessPayment(c *gin.Context) {
	var req PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(fmt.Errorf("参数解析失败: %v", err))
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "参数错误",
			Data:    nil,
		})
		return
	}

	params := map[string]string{
		"order_id":     req.OrderID,
		"amount":       fmt.Sprintf("%.2f", req.Amount),
		"payment_type": req.PaymentType,
	}

	if err := h.service.ProcessPayment(params); err != nil {
		c.Error(fmt.Errorf("支付处理失败: %v", err))
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "支付处理失败",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "支付处理成功",
		Data:    nil,
	})
}

// GetHmfCi 处理获取hmfCi请求
func (h *Handler) GetHmfCi(c *gin.Context) {
	var req HmfCiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(fmt.Errorf("参数解析失败: %v", err))
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "参数错误",
			Data:    nil,
		})
		return
	}

	params := map[string]string{
		"user_id": req.UserID,
		"token":   req.Token,
	}

	if err := h.service.GetHmfCi(params); err != nil {
		c.Error(fmt.Errorf("获取hmfCi失败: %v", err))
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "获取hmfCi失败",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "获取hmfCi成功",
		Data:    nil,
	})
}
