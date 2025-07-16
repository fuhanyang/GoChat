package chain

import (
	"common/routinePool"
	"common/zap"
	"context"
	"fmt"
	"log"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
)

type MyHandlerFunc func(c *MyContext) error

type stack []int

type MyContext struct {
	mutex    sync.Mutex //一个context需要互斥访问
	KV       map[string]interface{}
	Handlers []MyHandlerFunc
	stack    stack
	ctx      context.Context
	// 熔断相关字段
	isAborted    bool
	Error        error
	ErrorHandler func(error, *MyContext)                                 // 错误处理函数
	Logger       func(level string, format string, value ...interface{}) // 日志函数
	StartTime    time.Time
	HandlerNames []string // 用于记录每个Handler的名称
}

func (c *MyContext) Push(v int) {
	c.stack = append(c.stack, v)
}
func (c *MyContext) Pop() int {
	if len(c.stack) == 0 { // 直接操作 c.stack
		return 0
	}
	// 获取最后一个元素
	v := c.stack[len(c.stack)-1]
	c.stack = c.stack[:len(c.stack)-1]
	return v
}
func (l stack) len() int {
	return len(l)
}

// LoadHandlers 加载处理器，支持命名
func LoadHandlers(logger func(string, string, ...interface{}), timer MyHandlerFunc, handlers ...MyHandlerFunc) *MyContext {
	//TODO: 使用池化技术优化性能
	c := MyContext{}
	c.Handlers = make([]MyHandlerFunc, 0, len(handlers)+2)
	c.HandlerNames = make([]string, 0, len(handlers)+2)
	if timer != nil {
		c.Handlers = append(c.Handlers, timer)
	} else {
		c.Handlers = append(c.Handlers, DefaultTimer())
	}

	for _, handler := range handlers {
		c.Handlers = append(c.Handlers, handler)
		//默认名称是函数自己名称
		c.HandlerNames = append(c.HandlerNames, getFunctionName(handler))
	}
	c.KV = make(map[string]interface{})
	c.stack = make(stack, 0, len(handlers))
	c.Logger = logger
	if c.Logger == nil {
		c.Logger = DefaultLogger // 默认使用标准日志
	}

	return &c
}

// 辅助函数：获取函数名称
func getFunctionName(fn interface{}) string {
	if fn == nil {
		return "nil"
	}
	// 获取函数指针
	v := reflect.ValueOf(fn)
	if v.Kind() != reflect.Func {
		return "not a function"
	}
	// 通过函数指针获取函数信息
	ptr := v.Pointer()
	funcInfo := runtime.FuncForPC(ptr)
	if funcInfo == nil {
		return "unknown"
	}
	// 提取函数名（去除包路径前缀）
	fullName := funcInfo.Name()
	// 例如："your-package-name.F1" → "F1"
	lastSlash := strings.LastIndex(fullName, "/")
	if lastSlash == -1 {
		lastSlash = 0
	} else {
		lastSlash++
	}
	dotIndex := strings.Index(fullName[lastSlash:], ".")
	if dotIndex == -1 {
		return fullName
	}
	return fullName[lastSlash+dotIndex+1:]
}

// LoadHandlersWithNames 加载处理器并指定名称
func LoadHandlersWithNames(logger func(string, string, ...interface{}), timer MyHandlerFunc, handlerMap map[string]MyHandlerFunc) *MyContext {
	c := MyContext{}
	c.Handlers = make([]MyHandlerFunc, 0, len(handlerMap)+2)
	c.HandlerNames = make([]string, 0, len(handlerMap)+2)

	if timer != nil {
		c.Handlers = append(c.Handlers, timer)
	}

	for name, handler := range handlerMap {
		c.Handlers = append(c.Handlers, handler)
		c.HandlerNames = append(c.HandlerNames, name)
	}

	c.KV = make(map[string]interface{})
	c.stack = make(stack, 0, len(handlerMap))
	c.StartTime = time.Now()
	c.Logger = logger
	if c.Logger == nil {
		c.Logger = DefaultLogger // 默认使用标准日志
	}
	return &c
}
func (c *MyContext) Set(key string, value interface{}) {
	c.KV[key] = value
}
func (c *MyContext) SetWithMap(kv map[string]interface{}) {
	for k, v := range kv {
		c.Set(k, v)
	}
}
func (c *MyContext) Get(key string) interface{} {
	return c.KV[key]
}
func GetToType[T any](c *MyContext, key string) (T, error) {
	v, exists := c.KV[key]
	if !exists || v == nil {
		var zero T
		return zero, fmt.Errorf("key %q not found or value is nil", key)
	}

	// 尝试类型断言
	if result, ok := v.(T); ok {
		return result, nil
	}

	// 类型不匹配
	var zero T
	return zero, fmt.Errorf("value for key %q is not of type %T", key, zero)
}

// SetErrorHandler 设置错误处理函数
func (c *MyContext) SetErrorHandler(handler func(error, *MyContext)) {
	c.ErrorHandler = handler
}

// SetLogger 设置日志函数
func (c *MyContext) SetLogger(logger func(string, string, ...interface{})) {
	c.Logger = logger
}

func (c *MyContext) Next() {
	for {
		if c.isAborted {
			c.Logger("info", "Chain aborted, stopping execution")
			break
		}

		index := c.Pop()
		if index >= len(c.Handlers) {
			break
		}

		c.Push(index + 1)

		// 记录当前执行的Handler名称
		handlerName := "Unknown"
		if index < len(c.HandlerNames) {
			handlerName = c.HandlerNames[index]
		}

		// 执行Handler并捕获错误
		err := c.Handlers[index](c)
		if err != nil {
			c.Error = err
			c.Logger("Handler %s failed with error: %v", handlerName, err)

			// 调用错误处理函数
			if c.ErrorHandler != nil {
				c.ErrorHandler(err, c)
			}

			//熔断：停止执行后续Handler
			c.isAborted = true
			c.Abort()
			break
		}
		//执行结束，如果是最后一个则结束这层循环，向前回调
		if index >= len(c.Handlers)-1 {
			c.Pop()
			break
		}
		if c.stack.len() == 0 {
			break
		}
	}
}

func (c *MyContext) Abort() {
	c.isAborted = true
	c.stack = c.stack[:0]
}

func (c *MyContext) Apply() error {
	if len(c.Handlers) == 0 {
		return nil
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.StartTime = time.Now()
	//c.ZapLogger("Starting chain execution with %d handlers", len(c.Handlers))
	c.Next()
	return c.GetError()
}

func (c *MyContext) AppendHandler(handler MyHandlerFunc, name string) {
	c.Handlers = append(c.Handlers, handler)
	if name == "" {
		c.HandlerNames = append(c.HandlerNames, getFunctionName(handler))
	}
}

func (c *MyContext) Reset() {
	c.KV = make(map[string]interface{})
	c.stack = make(stack, 0, len(c.Handlers))
	c.isAborted = false
	c.Error = nil
	c.StartTime = time.Now()
}

func (c *MyContext) Done() {
	c.ctx.Done()
}

// GetExecutionTime 获取执行时间
func (c *MyContext) GetExecutionTime() time.Duration {
	return time.Since(c.StartTime)
}

// GetError 获取错误信息
func (c *MyContext) GetError() error {
	return c.Error
}

func DefaultTimer() MyHandlerFunc {
	return func(c *MyContext) error {
		c.Next()
		duration := c.GetExecutionTime()
		if c.Error != nil {
			c.Logger("error", "Chain execution failed after %v: %v", duration, c.Error)
			return c.Error
		} else {
			c.Logger("info", "Chain execution completed successfully in %v", duration)
			return nil
		}
	}
}
func DefaultLogger(level string, format string, value ...interface{}) {
	log.Printf(format, value...)
}
func ZapLogger(level string, format string, value ...interface{}) {
	sugarLogger := zap.GetSugarLogger()
	switch level {
	case "info":
		sugarLogger.Info(fmt.Sprintf(format, value...))
	case "error":
		sugarLogger.Error(fmt.Sprintf(format, value...))
	}
}
func CreateTask(ctx *MyContext) routinePool.Task {
	var task routinePool.Task
	task = func() error {
		return ctx.Apply()
	}
	return task
}
