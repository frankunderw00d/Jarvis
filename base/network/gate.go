package network

import (
	"errors"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	"jarvis/base/log"
	gRPC "jarvis/base/network/grpc"
	"net"
	"net/http"
)

type (
	// 入口定义
	Gate interface {
		// 入口名称
		Name() string

		// 初始化
		Initialize() error

		// 运行
		Running(func(Conn)) error

		// 销毁
		Destroy() error
	}

	// 入口基础
	baseGate struct {
		address string // 地址
	}

	// socket 入口实现
	socketGate struct {
		baseGate
		listener net.Listener // 监听者
	}

	// websocket 入口实现
	webSocketGate struct {
		baseGate
		server   *http.Server
		upgrader *websocket.Upgrader
		f        func(Conn)
	}

	// gRPC 入口实现
	gRPCGate struct {
		baseGate
		server   *grpc.Server
		listener net.Listener
		f        func(Conn)
	}
)

const (
	// 默认 WebSocket 路径
	DefaultWebSocketPath = "/ws"
	// 默认升级器读大小
	DefaultUpgraderReadBufferSize = 1024
	// 默认升级器写大小
	DefaultUpgraderWriteBufferSize = 1024
)

const (
	// Socket gate 名称
	SocketGateName = "Socket Gate"
	// WebSocket gate 名称
	WebSocketGateName = "WebSocket Gate"
	// gRPC gate 名称
	GRPCGateName = "GRPC Gate"
	// network
	DefaultNetwork = "tcp"
)

// 此变量组定义了 Gate 定义及实现中可能会发生的错误文本
const (
	ErrEmptyAddressText = "address is empty"
	ErrNilListenerText  = "listener is nil"
	ErrNilServerText    = "server is nil"
	ErrNilUpgraderText  = "upgrader is nil"
	ErrNilHookFuncText  = "hook function is nil"
)

// 此变量组定义了 Gate 定义及实现中可能会发生的错误
var (
	// 地址为空 错误
	ErrEmptyAddress = errors.New(ErrEmptyAddressText)
	// 监听者为 nil 错误
	ErrNilListener = errors.New(ErrNilListenerText)
	// server 为 nil 错误
	ErrNilServer = errors.New(ErrNilServerText)
	// Upgrader 为 nil 错误
	ErrNilUpgrader = errors.New(ErrNilUpgraderText)
	// hook function 为 nil 错误
	ErrNilHookFunc = errors.New(ErrNilHookFuncText)
)

// 实例化 socketGate 的 Gate 实现，返回表现形式为 Gate
func NewSocketGate(addr string) Gate {
	return &socketGate{
		baseGate: baseGate{address: addr},
		listener: nil,
	}
}

// 实例化 webSocketGate 的 Gate 实现，返回表现形式为 Gate
func NewWebSocketGate(addr string) Gate {
	return &webSocketGate{
		baseGate: baseGate{address: addr},
	}
}

// 实例化 gRPCGate 的 Gate 实现，返回表现形式为 Gate
func NewGRPCGate(addr string) Gate {
	return &gRPCGate{
		baseGate: baseGate{address: addr},
	}
}

// --------------------------------------------------- SocketGate ------------------------------------------------------
// 入口名称
func (sg *socketGate) Name() string {
	return SocketGateName
}

// 初始化
func (sg *socketGate) Initialize() error {
	if sg.address == "" {
		return ErrEmptyAddress
	}

	// 实例化监听
	l, err := net.Listen(DefaultNetwork, sg.address)
	if err != nil {
		return err
	}

	sg.listener = l

	return nil
}

// 运行
func (sg *socketGate) Running(function func(Conn)) error {
	if sg.listener == nil {
		return ErrNilListener
	}
	if function == nil {
		return ErrNilHookFunc
	}

	var e error
	for {
		c, err := sg.listener.Accept()
		if err != nil {
			e = err // 捕捉接收连接错误，反馈到 Service 中
			break
		}
		conn := NewSocketConn(c)
		function(conn) // 实例化 Conn ，下放到 Service.Manager 中管理
	}
	return e
}

// 销毁
func (sg *socketGate) Destroy() error {
	// 关闭监听并且置为 nil
	if err := sg.listener.Close(); err != nil {
		return err
	}

	sg.listener = nil

	return nil
}

// --------------------------------------------------- WebSocketGate ---------------------------------------------------
// 入口名称
func (wsg *webSocketGate) Name() string {
	return WebSocketGateName
}

// 初始化
func (wsg *webSocketGate) Initialize() error {
	if wsg.address == "" {
		return ErrEmptyAddress
	}

	// 实例化 server，处理器指定为自身，自身实现了 ServeHTTP(http.ResponseWriter,*http.Request)
	server := &http.Server{
		Addr:    wsg.address,
		Handler: wsg,
	}

	// 实例化升级器，用于升级链接
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  DefaultUpgraderReadBufferSize,
		WriteBufferSize: DefaultUpgraderWriteBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	wsg.server = server
	wsg.upgrader = upgrader

	return nil
}

// 运行
func (wsg *webSocketGate) Running(function func(Conn)) error {
	if wsg.server == nil {
		return ErrNilServer
	}
	if wsg.upgrader == nil {
		return ErrNilUpgrader
	}
	if function == nil {
		return ErrNilHookFunc
	}

	// 持有钩子函数，用于 ServeHTTP() 函数中使用
	wsg.f = function

	return wsg.server.ListenAndServe()
}

// 销毁
func (wsg *webSocketGate) Destroy() error {
	if wsg.server == nil {
		return ErrNilServer
	}
	if wsg.upgrader == nil {
		return ErrNilUpgrader
	}

	// 关闭 server
	return wsg.server.Close()
}

// 内部 HTTP 服务器函数
func (wsg *webSocketGate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 路径非 /ws
	if r.URL.Path != DefaultWebSocketPath {
		return
	}

	// 升级请求为 webSocket
	c, err := wsg.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.ErrorF("websocket Gate upgrade request error : %s", err.Error())
		return
	}

	// 实例化 Conn 交由 Service
	conn := NewWebSocketConn(c)
	if wsg.f != nil {
		wsg.f(conn)
	}
}

// --------------------------------------------------- GRPCGate --------------------------------------------------------
// 入口名称
func (gg *gRPCGate) Name() string {
	return GRPCGateName
}

// 初始化
func (gg *gRPCGate) Initialize() error {
	if gg.address == "" {
		return ErrEmptyAddress
	}

	// 实例化监听和 server
	l, err := net.Listen(DefaultNetwork, gg.address)
	if err != nil {
		return err
	}

	server := grpc.NewServer()

	gg.server = server
	gg.listener = l

	return nil
}

// 运行
func (gg *gRPCGate) Running(function func(Conn)) error {
	if gg.listener == nil {
		return ErrNilListener
	}
	if function == nil {
		return ErrNilHookFunc
	}

	// 持有钩子函数
	gg.f = function

	// 注册 server ，指定处理器为自身，自身实现了 CommunicateServer 接口
	gRPC.RegisterCommunicateServer(gg.server, gg)

	return gg.server.Serve(gg.listener)
}

// 销毁
func (gg *gRPCGate) Destroy() error {
	// 停止 server
	gg.server.Stop()

	// 停止监听
	return gg.listener.Close()
}

func (gg *gRPCGate) Connect(ccs gRPC.Communicate_ConnectServer) error {
	// 实例化 Conn ，且得到一个关闭的 channel
	conn, closeChannel := NewGRPCConn(ccs)

	// 将 Conn 发送到 Service
	if gg.f != nil {
		gg.f(conn)
	}

	// 阻塞当前调用，防止 ccs 关闭
	select {
	case <-closeChannel:
		return nil
	}
}
