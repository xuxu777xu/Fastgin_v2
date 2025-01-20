package example

// api example UserRequest 用户请求结构体
type UserRequest struct {
	Username string `json:"username" binding:"required"`
	Age      int    `json:"age" binding:"required,gte=0,lte=150"`
}

// api example UserResponse 用户响应结构体
type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Age      int    `json:"age"`
}

// Handler 处理器结构体
type Handler struct{}

// NewHandler 创建处理器实例
func NewHandler() *Handler {
	return &Handler{}
}
