package network

import (
	"errors"
	"log"
)

type (
	// 服务定义
	Service interface {
		// 注册中间件
		// 必须在注册路由之前调用
		UseMiddleware(...RouteHandleFunc) error

		// 注册路由
		RegisterModule(...Module) error

		// 注册观察者
		// 此函数必须在 Run() 前调用
		RegisterObserver(Observer) error

		// 非阻塞运行，接收不定数量的 Gate ，每个 Gate 实例都会在同一个线程内串行执行 Initialize()-Running()-Destroy()
		// 因此，务必确保 Gate.Running() 函数是阻塞式的
		Run(...Gate) error
	}

	// 服务定义实现
	service struct {
		manager            Manager        // 端管理
		router             Router         // 路由管理
		packager           Packager       // 装包者
		IntoStream         chan Message   // 進入流
		rootCallLinkedList CallLinkedList // 根调用链
	}
)

const (
	// 默认进入流管道宽度
	DefaultIntoStreamSize = 5000
)

// 此常量组定义了 Service 定义及实现中可能会发生的错误文本
const (
	ErrEmptyGatesText          = "list of Gate is empty"
	ErrNilGatesText            = "list of Gate is nil"
	ErrEmptyModulesText        = "list of Module is empty"
	ErrNilModulesText          = "list of Module is nil"
	ErrNilObserverText         = "observer is nil"
	ErrNilMiddlewareListText   = "list of middleware is nil"
	ErrEmptyMiddlewareListText = "list of middleware is empty"
)

// 此常量组定义了 Service 定义及实现中可能会发生的错误
var (
	// Gate 列表为空 错误
	ErrEmptyGates = errors.New(ErrEmptyGatesText)
	// Gate 列表为 nil 错误
	ErrNilGates = errors.New(ErrNilGatesText)
	// Module 列表为空 错误
	ErrEmptyModules = errors.New(ErrEmptyModulesText)
	// Module 列表为 nil 错误
	ErrNilModules = errors.New(ErrNilModulesText)
	// Observer 列表为 nil 错误
	ErrNilObserver = errors.New(ErrNilObserverText)
	// 中间件列表为 nil 错误
	ErrNilMiddlewareList = errors.New(ErrNilMiddlewareListText)
	// 中间件列表为空 错误
	ErrEmptyMiddlewareList = errors.New(ErrEmptyMiddlewareListText)
)

// 新建服务
func NewService(max, intoStreamSize int64) Service {
	// 如果 intoStreamSize 小于等于零，则默认为 5000
	if intoStreamSize <= 0 {
		intoStreamSize = DefaultIntoStreamSize
	}

	return &service{
		manager:            NewManage(max),
		router:             NewRouter(),
		packager:           DefaultPackager(),
		IntoStream:         make(chan Message, intoStreamSize),
		rootCallLinkedList: NewCallLinkedList(),
	}
}

// 注册路由
// 此函数为串行调用，因此 Service.router 的注册可以不需要加锁，目前内部实现加锁是为了在并发路由时避免竞态
// Service.manager 的 HookObserveConnect() 和 HookObserveDisConnect() 都没有进行加锁
func (s *service) RegisterModule(modules ...Module) error {
	if modules == nil {
		return ErrNilModules
	}
	if len(modules) == 0 {
		return ErrEmptyModules
	}

	// 遍历所有注册的 Module
	for _, module := range modules {
		// 为每个 Module 实例化一个调用链映射表
		routeMap := make(map[string]CallLinkedList)
		// 遍历 Module 的路由
		for path, handleFuncList := range module.Route() {
			// 复制当前 Service 的根调用链，其产生的基础为插入中间件调用节点
			cll := s.rootCallLinkedList.Copy()
			// 遍历 Module 路由的调用函数组，根据顺序创建节点加入复制后的调用链中
			for _, handleFunc := range handleFuncList {
				cll.AddNode(handleFunc)
			}
			// 将完成的调用链根据 path 路径加入到调用链映射表
			routeMap[path] = cll
		}

		// 将 Module 和重建的调用链映射表注册到 Service 持有的路由器中
		if err := s.router.RegisterRoute(module.Name(), routeMap); err != nil {
			return err
		}
	}

	return nil
}

// 注册中间件
// 此函数必须在 RegisterModule() 前调用，才能得到包含完整的中间件的根调用链
func (s *service) UseMiddleware(middleware ...RouteHandleFunc) error {
	if middleware == nil {
		return ErrNilMiddlewareList
	}
	if len(middleware) == 0 {
		return ErrEmptyMiddlewareList
	}

	// 在根调用链上加入新节点
	for _, mw := range middleware {
		s.rootCallLinkedList.AddNode(mw)

	}
	return nil
}

// 注册观察者
// 此函数必须在 Run() 前调用
func (s *service) RegisterObserver(observer Observer) error {
	if observer == nil {
		return ErrNilObserver
	}

	// 在 Service 的 manager(端管理) 中注册观察者
	s.manager.RegisterObserver(observer)

	return nil
}

// 非阻塞运行，接收不定数量的 Gate ，每个 Gate 实例都会在同一个线程内串行执行 Initialize()-Running()-Destroy()
// 因此，务必确保 Gate.Running() 函数是阻塞式的
func (s *service) Run(gates ...Gate) error {
	if gates == nil {
		return ErrNilGates
	}
	if len(gates) == 0 {
		return ErrEmptyGates
	}

	// 开启服务接收进入流
	go s.receive()

	// 开启入口群
	for _, g := range gates {
		go func(gate Gate) {
			log.Printf("[%s] gate start", gate.Name())

			if err := gate.Initialize(); err != nil {
				log.Printf("[%s] gate initizlize error : %s", gate.Name(), err.Error())
				return
			}

			if err := gate.Running(s.acceptConn); err != nil {
				log.Printf("[%s] gate running error : %s", gate.Name(), err.Error())
				return
			}

			if err := gate.Destroy(); err != nil {
				log.Printf("[%s] gate destroy error : %s", gate.Name(), err.Error())
				return
			}

			log.Printf("[%s] gate done", gate.Name())
		}(g)
	}
	return nil
}

// 接收进入流
func (s *service) receive() {
	for {
		// 接收请求
		req, ok := <-s.IntoStream
		if !ok {
			log.Printf("service intoStream-channel closed")
			break
		}

		// 路由
		handleCLL, err := s.router.RouteHandleFun(req.Module, req.Route)
		if err != nil {
			log.Printf("route [%s]-[%s] error : %s", req.Module, req.Route, err.Error())
			continue
		}

		// 开线程去处理
		// 因此 Module 的 Route() 中注册的函数必然会发生竞态，如果有需要保护的数据，需要开发者自己维护
		go func(cll CallLinkedList, request Message) {
			// 新建上下文
			ctx := NewContext(request)
			// 上下文钩住当前 Service 的 manager(端管理) 的查找函数
			ctx.HookFind(s.manager.FindItem)
			// 调用链持有上下文开始按加入节点顺序调用
			cll.Run(ctx)
		}(handleCLL, req)
	}
}

// 接收 Gate 下放的连接
func (s *service) acceptConn(conn Conn) {
	i := NewItem(conn, s.packager.Clone(), DefaultEncrypter())

	if err := s.manager.ManageItem(i); err != nil {
		log.Printf("service manage new item [%s] error : %s", i.ID().String(), err.Error())
		// todo : send a message to the connection to tell it why
		i.Close()
		i = nil
		return
	}

	// 加入管理后再 Hook 和 开始接收消息
	i.PassiveCloseFeedback(s.manager.RemoveItem)
	go i.Receive(s.IntoStream)
}
