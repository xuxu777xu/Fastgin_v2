package apiServer

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

// LoginService 登录服务结构体
type LoginService struct {
	client *resty.Client
}

// NewLoginService 创建登录服务实例
func NewLoginService() *LoginService {
	return &LoginService{
		client: resty.New(),
	}
}

// Login 登录
func (s *LoginService) Login(params map[string]string) error {
	//TODO: 从参数中获取信息,运行,获取的信息不同来判断是否继续执行 GetUnpaidOrders 函数  成功继续执行，失败返回请求错误信息，日志记录
	//TODO: 根据账号密码的信息组成map[string]string
	fmt.Println("登录")
	return nil
}

// GetUnpaidOrders 获取未支付订单列表
func (s *LoginService) GetUnpaidOrders(params map[string]string) error {
	//TODO: 从参数中获取信息,运行,获取的信息不同来判断是否继续执行 ProcessPayment 函数  成功继续执行，失败返回请求错误信息，日志记录
	fmt.Println("获取未支付订单列表")
	return nil
}

// ProcessPayment 支付流程
func (s *LoginService) ProcessPayment(params map[string]string) error {
	//TODO: 从参数中获取信息,运行,获取的信息不同来判断是否继续执行 GetHmfCi 函数  成功继续执行，失败返回请求错误信息，日志记录
	fmt.Println("支付流程")
	return nil
}

// GetHmfCi 获取hmfCi
func (s *LoginService) GetHmfCi(params map[string]string) error {
	//TODO: 从参数中获取信息,运行,获取的信息不同来判断是否继续执行 GetHmfCi 函数  成功返回给请求一些信息，失败返回请求错误信息，日志记录
	fmt.Println("获取hmfCi")
	return nil
}
