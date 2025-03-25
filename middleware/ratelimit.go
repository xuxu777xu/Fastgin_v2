package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimit 创建限流中间件
func RateLimit(maxRequests int, perSecond float64) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(perSecond), maxRequests)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"message": "请求太频繁，请稍后再试",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
