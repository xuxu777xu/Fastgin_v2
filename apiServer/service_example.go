package apiServer

import (
	"FastGin/api/common"
	"fmt"
	"github.com/imroc/req/v3"
	"net/http"
)

// ExampleService 示例服务结构体
type ExampleService struct {
	client *req.Client
}

// NewExampleService 创建示例服务实例
func NewExampleService() *ExampleService {
	return &ExampleService{
		client: req.C(),
	}
}

// Login 登录接口实现
func (s *ExampleService) Login(params map[string]string) (*common.Response, error) {
	// 创建新的HTTP客户端
	client := req.C()

	// 设置HTTP2
	client.EnableForceHTTP2()

	// 定义URL
	url := "https://capi.lkcoffee.com/resource/core/v1/order/create"

	// 准备表单数据
	formData := map[string]string{
		"sign": "630325412989152581553297576247734904",
		"uid":  params["uid"],
		"t":    "1742867714121",
		"cid":  "210101",
	}

	// 设置Cookie
	uidCookie := &http.Cookie{
		Name:  "uid",
		Value: "bf7dd9ba-3c02-4e4d-bc37-e06a46ddcd6b1742261779418",
	}

	// 发送请求
	resp, err := client.R().
		SetHeaders(map[string]string{
			"User-Agent":      "okhttp/4.9.3",
			"Accept-Encoding": "gzip",
			"Content-Type":    "application/x-www-form-urlencoded",
			"x-lk-akv":        "5205",
		}).
		SetCookies(uidCookie).
		SetFormData(formData).
		Post(url)

	// 记录响应
	fmt.Printf("响应体: %s\n", resp.String())
	result := resp.String()

	if err != nil {
		return &common.Response{
			Code:    500,
			Message: "请求有误",
			Data:    result,
		}, err
	}

	// 返回成功信息
	return &common.Response{
		Code:    200,
		Message: "请求成功",
		Data:    result,
	}, nil
}
