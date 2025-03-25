package test

import (
	"FastGin/api/common"
	"FastGin/apiServer"
)

type Handler struct {
	service ExampleServiceInterface
}

// ExampleServiceInterface 定义服务接口
type ExampleServiceInterface interface {
	Login(params map[string]string) (*common.Response, error)
}

// NewHandler 创建新的处理器实例
func NewHandler() *Handler {
	return &Handler{
		service: apiServer.NewExampleService(),
	}
}
