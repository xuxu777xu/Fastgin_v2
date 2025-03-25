package errcode

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse 定义错误结构
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewError 创建新的错误
func NewError(code int, msg string) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Message: msg,
	}
}

// 预定义错误
var (
	SuccessResponse     = NewError(200, "成功")
	InternalServerError = NewError(500, "服务内部错误")
	InvalidParams       = NewError(400, "参数错误")
	TooManyRequests     = NewError(429, "请求过多")
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    SuccessResponse.Code,
		Message: SuccessResponse.Message,
		Data:    data,
	})
}

// SendError 错误响应
func SendError(c *gin.Context, err *ErrorResponse) {
	c.JSON(http.StatusOK, Response{
		Code:    err.Code,
		Message: err.Message,
		Data:    nil,
	})
}

// ParamError 参数错误响应
func ParamError(c *gin.Context) {
	SendError(c, InvalidParams)
}

// ServerErrorResponse 服务器错误响应
func ServerErrorResponse(c *gin.Context) {
	SendError(c, InternalServerError)
}

// 参数验证错误
func Error(c *gin.Context, errcode int, msg string) {
	SendError(c, InvalidParams)
}
