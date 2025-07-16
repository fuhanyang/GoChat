package inject

import (
	"sync"
)

var injectMutex sync.Mutex

// Inject 注入rpc服务以及其对应的handler的request和response实例
func Inject(serviceName string, initFunc func()) {
	injectMutex.Lock()
	ServicesName = append(ServicesName, serviceName)
	InitFunctions = append(InitFunctions, initFunc)
	injectMutex.Unlock()
}

var InitFunctions []func()

var ServicesName []string
