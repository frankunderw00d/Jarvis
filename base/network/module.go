package network

import "errors"

type (
	// 模块定义
	Module interface {
		// 模块名
		Name() string

		// 模块路由
		Route() map[string][]RouteHandleFunc
	}

	// 模块定义实现
	module struct {
		route map[string][]RouteHandleFunc // 路由簇
	}
)

const (
	// 默认路由名
	DefaultModuleName = "default"
)

// 此常量组定义了 Module 定义及实现中可能会发生的错误文本
const (
	ErrNilHandleFuncText   = "route list is nil"
	ErrEmptyHandleFuncText = "route list is empty"
	ErrRouteExistText      = "route exist"
)

// 此常量组定义了 Module 定义及实现中可能会发生的错误
var (
	// HandleFunc 为 nil 错误
	ErrNilHandleFunc = errors.New(ErrNilHandleFuncText)
	// HandleFunc 列表为空 错误
	ErrEmptyHandleFunc = errors.New(ErrEmptyHandleFuncText)
	// route 名路径存在 错误
	ErrRouteExist = errors.New(ErrRouteExistText)
)

var (
	// 默认内部模块
	defaultModule = &module{
		route: make(map[string][]RouteHandleFunc),
	}
)

// 模块名
func (m *module) Name() string {
	return DefaultModuleName
}

// 模块路由
func (m *module) Route() map[string][]RouteHandleFunc {
	return m.route
}

// 模块内部注册路由
func (m *module) registerRoute(route string, handleFunc ...RouteHandleFunc) error {
	if route == "" {
		return ErrEmptyRouteName
	}
	if handleFunc == nil {
		return ErrNilHandleFunc
	}
	if len(handleFunc) == 0 {
		return ErrEmptyHandleFunc
	}

	// 默认模块通过 network.go 中的 RegisterRoute() 函数对外开放注册路由
	// 因此所有通过 RegisterRoute() 注册的路由，模块名都为 "default"
	if _, exist := m.route[route]; exist {
		return ErrRouteExist
	}

	m.route[route] = handleFunc

	return nil
}
