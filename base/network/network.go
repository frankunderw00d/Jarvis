package network

import (
	"jarvis/util/encrypt"
	"jarvis/util/rand"
)

type (
	// 内部唯一标识
	ID string
)

const ()

var (
	// 默认服务定义
	defaultService Service
)

func init() {
	// 默认服务
	defaultService = NewService(
		DefaultMaxConnection,
		DefaultIntoStreamSize,
	)
}

// 加密 ID，长度为8的前缀+ID.String()+长度为8的后缀，经 MD5 摘要后获得最终加密结果
func EncryptID(id ID) ID {
	return ID(encrypt.MD5(rand.RandomString(8) + id.String() + rand.RandomString(8)))
}

// native employment
func (id ID) String() string {
	return string(id)
}

// 向默认模块注册路由
// warning : 多线程注册默认模块路由可能引发竞态问题，默认模块注册路由未加锁，且此函数必须在 RegisterModule() 函数前调用
func RegisterRoute(route string, handleFunc ...RouteHandleFunc) error {
	return defaultModule.registerRoute(route, handleFunc...)
}

// 注册中间件
// 必须在注册模块之前调用
func UseMiddleware(middleware ...RouteHandleFunc) error {
	return defaultService.UseMiddleware(middleware...)
}

// 注册观察者
// 此函数必须在 Run() 前调用
func RegisterObserver(observer Observer) error {
	return defaultService.RegisterObserver(observer)
}

// 注册模块
func RegisterModule(modules ...Module) error {
	// 如果 module 列表不为空且有值，直接注册
	if modules != nil && len(modules) != 0 {
		return defaultService.RegisterModule(modules...)
	}

	// 否则直接注册默认模块
	return defaultService.RegisterModule(defaultModule)
}

// 运行
func Run(gates ...Gate) error {
	return defaultService.Run(gates...)
}
