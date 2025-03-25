package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestID 中间件用于生成请求ID和追踪ID
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {

		requestID := uuid.New().String()
		traceID := requestID

		// 设置到上下文
		c.Set("RequestID", requestID)
		c.Set("TraceID", traceID)

		// 设置响应头
		c.Header("X-Request-ID", requestID)
		c.Header("X-Trace-ID", traceID)

		c.Next()
	}
}
