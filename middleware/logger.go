package middleware

import (
	"FastGin/pkg/logg"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	// 确保错误日志目录存在
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		logg.Error("创建日志目录失败", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// 打开或创建错误日志文件
	errorLogPath := filepath.Join(logDir, "error.log")
	errorLogFile, err := os.OpenFile(errorLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logg.Error("打开错误日志文件失败", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return func(c *gin.Context) {
		// 获取请求ID和追踪ID
		requestID, exists := c.Get("RequestID")
		if !exists {
			requestID = "unknown"
		}
		traceID, exists := c.Get("TraceID")
		if !exists {
			traceID = requestID
		}

		// 记录请求开始
		logg.Info("开始处理请求", map[string]interface{}{
			"request_id": requestID,
			"trace_id":   traceID,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		})

		// 开始时间
		startTime := time.Now()

		// 记录请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 使用自定义ResponseWriter记录响应
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)

		// 获取请求信息
		reqInfo := map[string]interface{}{
			"status_code":   c.Writer.Status(),
			"latency":       latencyTime.String(),
			"client_ip":     c.ClientIP(),
			"method":        c.Request.Method,
			"uri":           c.Request.RequestURI,
			"path":          c.Request.URL.Path,
			"query_params":  c.Request.URL.RawQuery,
			"headers":       c.Request.Header,
			"request_body":  string(requestBody),
			"response_body": blw.body.String(),
			"error":         c.Errors.String(),
			"user_agent":    c.Request.UserAgent(),
			"request_id":    requestID,
			"trace_id":      traceID,
		}

		// 根据状态码决定日志级别
		statusCode := c.Writer.Status()
		switch {
		case statusCode >= 500:
			logg.Error("请求处理失败", reqInfo)
			// 记录详细错误信息到错误日志文件
			writeErrorLog(errorLogFile, c, requestBody, blw.body.Bytes(), fmt.Sprint(requestID), fmt.Sprint(traceID), latencyTime)
		case statusCode >= 400:
			logg.Warn("请求参数错误", reqInfo)
			// 记录详细错误信息到错误日志文件
			writeErrorLog(errorLogFile, c, requestBody, blw.body.Bytes(), fmt.Sprint(requestID), fmt.Sprint(traceID), latencyTime)
		default:
			logg.Info("请求处理完成", reqInfo)
		}
	}
}

// compressJSON 压缩JSON字符串
func compressJSON(data []byte) string {
	// 如果不是JSON格式，直接返回原始字符串
	if !json.Valid(data) {
		return string(data)
	}

	// 解析JSON
	var obj interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return string(data)
	}

	// 重新编码为压缩格式
	compressed, err := json.Marshal(obj)
	if err != nil {
		return string(data)
	}

	return string(compressed)
}

// writeErrorLog 写入详细的错误日志
func writeErrorLog(file *os.File, c *gin.Context, requestBody, responseBody []byte, requestID, traceID string, latency time.Duration) {
	// 构建基础错误信息
	errorLog := map[string]interface{}{
		"t":         time.Now().Format("2006-01-02T15:04:05.000Z"), // 时间戳
		"requestID": requestID,                                     // 请求ID
		"traceID":   traceID,                                       // 追踪ID
		"st":        c.Writer.Status(),                             // 状态码
		"m":         c.Request.Method,                              // 请求方法
		"p":         c.Request.URL.Path,                            // 请求路径
		"q":         c.Request.URL.RawQuery,                        // 查询参数
		"ip":        c.ClientIP(),                                  // 客户端IP
		"ua":        c.Request.UserAgent(),                         // 用户代理
		"l":         latency.Milliseconds(),                        // 响应时间(ms)
		"e":         strings.Join(c.Errors.Errors(), "; "),         // 错误信息
	}

	// 添加请求头（仅添加关键header）
	headers := make(map[string]string)
	for _, key := range []string{"Content-Type", "Authorization", "X-Request-ID", "X-Trace-ID"} {
		if value := c.GetHeader(key); value != "" {
			headers[key] = value
		}
	}
	if len(headers) > 0 {
		errorLog["h"] = headers
	}

	// 添加请求体（如果存在且是JSON则压缩）
	if len(requestBody) > 0 {
		errorLog["req"] = compressJSON(requestBody)
	}

	// 添加响应体（如果存在且是JSON则压缩）
	if len(responseBody) > 0 {
		errorLog["res"] = compressJSON(responseBody)
	}

	// 将错误日志转换为压缩的JSON格式
	logJSON, _ := json.Marshal(errorLog)

	// 写入错误日志文件（单行格式）
	file.WriteString(fmt.Sprintf("%s\n", string(logJSON)))
}

// bodyLogWriter 用于记录响应体
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
