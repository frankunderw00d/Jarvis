package network

import (
	"errors"
	"sync"
)

type (
	// 路径处理函数定义
	RouteHandleFunc func(Context)

	// 路由定义
	Router interface {
		// 注册模块及所属路由
		RegisterRoute(string, map[string]CallLinkedList) error

		// 通过模块名和路径名映射寻找处理函数
		RouteHandleFun(string, string) (CallLinkedList, error)
	}

	// 路由定义实现
	router struct {
		mutex sync.Mutex                           // 路由锁
		route map[string]map[string]CallLinkedList // 路由映射，map[模块]map[路径]处理函数
	}
)

// 此常量组定义了 Router 定义及实现中可能会发生的错误文本
const (
	ErrEmptyModuleNameText = "module name is empty"
	ErrEmptyRouteNameText  = "route name is empty"
	ErrNilRouteText        = "route is nil"
	ErrModuleExistText     = "module exist"
	ErrModuleUnExistText   = "module doesn't exist"
	ErrRouteUnExistText    = "route doesn't exist"
)

// 此常量组定义了 Router 定义及实现中可能会发生的错误
var (
	// module 名为空 错误
	ErrEmptyModuleName = errors.New(ErrEmptyModuleNameText)
	// route 名为空 错误
	ErrEmptyRouteName = errors.New(ErrEmptyRouteNameText)
	// route 映射为 nil 错误
	ErrNilRoute = errors.New(ErrNilRouteText)
	// module 名路由已存在 错误
	ErrModuleExist = errors.New(ErrModuleExistText)
	// module 名路由不存在 错误
	ErrModuleUnExist = errors.New(ErrModuleUnExistText)
	// route 名路径不存在 错误
	ErrRouteUnExist = errors.New(ErrRouteUnExistText)
)

// 新建路由
func NewRouter() Router {
	return &router{
		mutex: sync.Mutex{},
		route: make(map[string]map[string]CallLinkedList),
	}
}

// 注册模块及所属路由
func (r *router) RegisterRoute(module string, route map[string]CallLinkedList) error {
	if module == "" {
		return ErrEmptyModuleName
	}
	if route == nil {
		return ErrNilRoute
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	// 模块名已存在
	if _, exist := r.route[module]; exist {
		return ErrModuleExist
	}

	r.route[module] = route
	return nil
}

// 通过模块名和路径名映射寻找处理函数
func (r *router) RouteHandleFun(module, route string) (CallLinkedList, error) {
	if module == "" {
		return nil, ErrEmptyModuleName
	}
	if route == "" {
		return nil, ErrEmptyRouteName
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	// 模块名不存在
	routeMap, exist := r.route[module]
	if !exist {
		return nil, ErrModuleUnExist
	}

	// 路径名不存在
	if handleFunc, exist := routeMap[route]; !exist {
		return nil, ErrRouteUnExist
	} else {
		return handleFunc, nil
	}
}
