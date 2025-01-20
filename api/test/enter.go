package test

import (
	"FastGin/apiServer"
)

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UnpaidOrdersRequest 获取未支付订单请求参数
type UnpaidOrdersRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Page   int    `json:"page" binding:"required"`
	Size   int    `json:"size" binding:"required"`
}

// PaymentRequest 支付请求参数
type PaymentRequest struct {
	OrderID     string  `json:"order_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
	PaymentType string  `json:"payment_type" binding:"required"`
}

// HmfCiRequest 获取hmfCi请求参数
type HmfCiRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Token  string `json:"token" binding:"required"`
}

// Response 通用响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Handler struct {
	service LoginServiceInterface
}

// LoginServiceInterface 定义服务接口
type LoginServiceInterface interface {
	Login(params map[string]string) error
	GetUnpaidOrders(params map[string]string) error
	ProcessPayment(params map[string]string) error
	GetHmfCi(params map[string]string) error
}

// NewHandler 创建新的处理器实例
func NewHandler() *Handler {
	return &Handler{
		service: apiServer.NewLoginService(),
	}
}
