package logg

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Config 定义日志系统配置

type Config struct {
	LogDir      string
	MaxSize     int // MB
	MaxBackups  int
	MaxAge      int // days
	Compress    bool
	BufferSize  int // bytes
	EnableColor bool
}

// DefaultConfig 返回默认配置

func DefaultConfig() Config {
	return Config{
		LogDir:      "logs",
		MaxSize:     1024,
		MaxBackups:  7,
		MaxAge:      7,
		Compress:    true,
		BufferSize:  8 * 1024, // 8KB
		EnableColor: true,
	}
}

var (
	Log               *logrus.Logger
	currentDate       string
	fileWriter        *lumberjack.Logger
	writerPool        sync.Pool
	mu                sync.RWMutex
	onRotateCallbacks []func()
	isInitialized     bool
	config            Config
)

// CustomFormatter 自定义日志格式化
// EnableColors 是否开启终端颜色

type CustomFormatter struct {
	EnableColors bool
}

// LogData 统一的日志数据结构
type LogData struct {
	Timestamp  string                 `json:"timestamp"`
	Level      string                 `json:"level"`
	Message    string                 `json:"message"`
	Caller     string                 `json:"caller,omitempty"`
	Function   string                 `json:"function,omitempty"`
	Line       int                    `json:"line,omitempty"`
	TraceID    string                 `json:"trace_id,omitempty"`
	RequestID  string                 `json:"request_id,omitempty"`
	Method     string                 `json:"method,omitempty"`
	Path       string                 `json:"path,omitempty"`
	ClientIP   string                 `json:"client_ip,omitempty"`
	StatusCode int                    `json:"status_code,omitempty"`
	Latency    string                 `json:"latency,omitempty"`
	UserAgent  string                 `json:"user_agent,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty"`
	Error      string                 `json:"error,omitempty"`
}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	level := strings.ToUpper(entry.Level.String())
	message := entry.Message

	// 分离基础字段和自定义字段
	baseFields := make(map[string]interface{})
	contextFields := make(map[string]interface{})

	for k, v := range entry.Data {
		switch k {
		case "file", "func", "line", "request_id", "trace_id":
			baseFields[k] = v
		default:
			contextFields[k] = v
		}
	}

	var output string
	if f.EnableColors {
		// 控制台输出：简洁模式
		var colorCode string
		switch entry.Level {
		case logrus.DebugLevel:
			colorCode = "\033[36m" // 青色
		case logrus.InfoLevel:
			colorCode = "\033[32m" // 绿色
		case logrus.WarnLevel:
			colorCode = "\033[33m" // 黄色
		case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
			colorCode = "\033[31m" // 红色
		}

		// 基础日志信息（简洁版）
		output = fmt.Sprintf("[%s] [%s%s\033[0m] %s",
			timestamp,
			colorCode,
			level,
			message)

		// 只在控制台显示关键的上下文数据
		if len(contextFields) > 0 {
			// 选择性显示重要字段
			importantFields := make(map[string]interface{})
			for k, v := range contextFields {
				// 只显示关键字段
				if k == "error" || k == "path" || k == "method" || k == "status_code" {
					importantFields[k] = v
				}
			}
			if len(importantFields) > 0 {
				contextJSON, _ := json.Marshal(importantFields)
				output += fmt.Sprintf(" %s", string(contextJSON))
			}
		}
	} else {
		// 文件输出：详细模式，但保持易读格式
		output = fmt.Sprintf("[%s] [%s] %s", timestamp, level, message)

		// 添加详细的上下文数据
		if len(contextFields) > 0 {
			contextJSON, _ := json.Marshal(contextFields)
			output += fmt.Sprintf(" %s", string(contextJSON))
		}

		//// 添加调用者信息
		//if funcName, ok := baseFields["func"].(string); ok {
		//	if idx := strings.LastIndex(funcName, "/"); idx != -1 {
		//		funcName = funcName[idx+1:]
		//	}
		//	output += fmt.Sprintf(" [%s]", funcName)
		//}

		//// 添加文件和行号
		//if file, ok := baseFields["file"].(string); ok {
		//	if line, exists := baseFields["line"].(int); exists {
		//		output += fmt.Sprintf(" [%s:%d]", path.Base(file), line)
		//	}
		//}

		// 添加请求ID和追踪ID
		if requestID, ok := baseFields["request_id"].(string); ok && requestID != "unknown" {
			output += fmt.Sprintf(" [REQ:%s]", requestID)
		}
		if traceID, ok := baseFields["trace_id"].(string); ok && traceID != "unknown" {
			output += fmt.Sprintf(" [TRACE:%s]", traceID)
		}

	}

	return []byte(output + "\n"), nil
}

// InitLogger 初始化日志系统

func InitLogger(cfg Config) error {
	mu.Lock()
	defer mu.Unlock()

	if isInitialized {
		return fmt.Errorf("logger already initialized")
	}

	config = cfg

	if err := os.MkdirAll(config.LogDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	// 创建两个Logger实例
	Log = logrus.New()
	fileLog := logrus.New()

	// 控制台logger配置
	Log.SetReportCaller(true)
	Log.SetFormatter(&CustomFormatter{EnableColors: config.EnableColor})
	Log.SetOutput(os.Stdout)

	// 文件Logger配置
	fileLog.SetReportCaller(true)
	fileLog.SetFormatter(&CustomFormatter{EnableColors: false})

	if err := initFileWriter(); err != nil {
		return fmt.Errorf("failed to initialize file writer: %w", err)
	}

	// 设置文件输出
	fileLog.SetOutput(fileWriter)

	// 添加Hook将日志同时写入文件
	Log.AddHook(&writeHook{fileLog})

	// 添加caller hook
	Log.Hooks.Add(&callerHook{})
	fileLog.Hooks.Add(&callerHook{})

	go checkLogFileDaily()

	isInitialized = true
	return nil
}

// writeHook 用于将日志同时写入文件
type writeHook struct {
	fileLogger *logrus.Logger
}

func (h *writeHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *writeHook) Fire(entry *logrus.Entry) error {
	// 创建新的Entry，保留所有字段
	newEntry := logrus.NewEntry(h.fileLogger)
	newEntry.Data = make(logrus.Fields, len(entry.Data))

	// 复制所有字段
	for k, v := range entry.Data {
		newEntry.Data[k] = v
	}

	// 设置消息和日志级别
	newEntry.Message = entry.Message
	newEntry.Level = entry.Level

	// 如果有调用信息，也复制过来
	if entry.Caller != nil {
		newEntry.Caller = entry.Caller
	}

	// 使用新Entry记录日志
	newEntry.Log(entry.Level, entry.Message)
	return nil
}

// callerHook 用于获取正确的调用者信息
type callerHook struct{}

func (h *callerHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *callerHook) Fire(entry *logrus.Entry) error {
	// 获取调用栈信息
	pcs := make([]uintptr, 10)
	n := runtime.Callers(6, pcs) // 调整调用栈深度
	if n == 0 {
		return nil
	}

	frames := runtime.CallersFrames(pcs[:n])
	// 寻找第一个非日志库的调用者
	for {
		frame, more := frames.Next()
		// 跳过日志库内部的调用
		if !strings.Contains(frame.File, "logg/") &&
			!strings.Contains(frame.File, "logrus") &&
			!strings.Contains(frame.Function, "logrus.") {
			entry.Caller = &runtime.Frame{
				Function: frame.Function,
				File:     frame.File,
				Line:     frame.Line,
			}
			break
		}
		if !more {
			break
		}
	}
	return nil
}

// initFileWriter 初始化文件写入器

func initFileWriter() error {
	currentDate = time.Now().Format("2006-01-02")
	fileWriter = &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s.log", config.LogDir, currentDate),
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}
	return nil
}

// checkLogFileDaily 检查并轮转日志文件

func checkLogFileDaily() {
	for {
		now := time.Now()
		nextDay := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		time.Sleep(time.Until(nextDay))

		mu.Lock()
		newDate := time.Now().Format("2006-01-02")
		if newDate != currentDate {
			currentDate = newDate
			fileWriter.Filename = fmt.Sprintf("%s/%s.log", config.LogDir, currentDate)
			if err := fileWriter.Rotate(); err != nil {
				Log.Errorf("Failed to rotate log file: %v", err)
			}
			// 执行回调
			for _, callback := range onRotateCallbacks {
				callback()
			}
		}
		mu.Unlock()
	}
}

// asyncWriter 异步写入器
type asyncWriter struct{}

func (w *asyncWriter) Write(p []byte) (n int, err error) {
	// 复制数据以减少锁定时间
	data := make([]byte, len(p))
	copy(data, p)

	mu.RLock()
	bufWriter := writerPool.Get().(*bufio.Writer)
	mu.RUnlock()

	defer writerPool.Put(bufWriter)
	bufWriter.Reset(fileWriter)

	n, err = bufWriter.Write(data)
	if err != nil {
		return n, fmt.Errorf("buffer write failed: %w", err)
	}

	if err = bufWriter.Flush(); err != nil {
		return n, fmt.Errorf("buffer flush failed: %w", err)
	}

	return n, nil
}

// AddRotateCallback 添加日志轮转回调函数

func AddRotateCallback(callback func()) {
	mu.Lock()
	defer mu.Unlock()
	onRotateCallbacks = append(onRotateCallbacks, callback)
}

// Shutdown 优雅关闭日志系统

func Shutdown() error {
	mu.Lock()
	defer mu.Unlock()

	if !isInitialized {
		return nil
	}

	// 获取一个缓冲写入器并刷新
	bufWriter := writerPool.Get().(*bufio.Writer)
	defer writerPool.Put(bufWriter)

	if err := bufWriter.Flush(); err != nil {
		return fmt.Errorf("failed to flush buffer on shutdown: %w", err)
	}

	if err := fileWriter.Close(); err != nil {
		return fmt.Errorf("failed to close file writer: %w", err)
	}

	isInitialized = false
	return nil
}

// logWithContext 函数修改
func logWithContext(level logrus.Level, args ...interface{}) {
	fields := logrus.Fields{}

	// 获取调用者信息
	pc, file, line, ok := runtime.Caller(2)
	if ok {
		// 获取完整的函数名
		if fn := runtime.FuncForPC(pc); fn != nil {
			fields["func"] = fn.Name()
		}
		fields["file"] = file
		fields["line"] = line
	}

	var msg string
	var contextData map[string]interface{}

	// 处理参数
	if len(args) > 0 {
		// 第一个参数作为消息
		msg = fmt.Sprint(args[0])

		// 如果有第二个参数且是map，作为上下文数据
		if len(args) > 1 {
			if data, ok := args[len(args)-1].(map[string]interface{}); ok {
				contextData = data
				// 合并上下文数据到fields
				for k, v := range contextData {
					fields[k] = v
				}
			}
		}
	}

	// 确保关键信息不为空
	if fields["request_id"] == nil {
		fields["request_id"] = "unknown"
	}
	if fields["trace_id"] == nil {
		fields["trace_id"] = "unknown"
	}

	entry := Log.WithFields(fields)
	switch level {
	case logrus.DebugLevel:
		entry.Debug(msg)
	case logrus.InfoLevel:
		entry.Info(msg)
	case logrus.WarnLevel:
		entry.Warn(msg)
	case logrus.ErrorLevel:
		entry.Error(msg)
	case logrus.FatalLevel:
		entry.Fatal(msg)
	}
}

// 更新日志函数
func Debug(args ...interface{}) {
	logWithContext(logrus.DebugLevel, args...)
}

func Info(args ...interface{}) {
	logWithContext(logrus.InfoLevel, args...)
}

func Warn(args ...interface{}) {
	logWithContext(logrus.WarnLevel, args...)
}

func Error(args ...interface{}) {
	logWithContext(logrus.ErrorLevel, args...)
}

func Fatal(args ...interface{}) {
	logWithContext(logrus.FatalLevel, args...)
}
