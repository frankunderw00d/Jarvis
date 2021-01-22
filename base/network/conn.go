// Conn 是对 socket-conn、websocket-conn、grpc-conn 的一种封装，将这种长连接的读取、写入、关闭等功能统一，便于上层处理
// Conn 的读取，获取的是统一的字节组数据，每次读取大小为 BufferLength (1024 byte)，不涉及解包、解密
// Conn 的写入，发送的是统一的字节组数据，大小不得超过底层文件句柄限制，即 1 << 30 == 1073741824 byte,且不涉及打包、加密
// Conn 的关闭，不可重复关闭，内部 net.Conn 实例为 nil 时报错，会在加互斥锁的情况下将 closed 状态值为 true
// Conn 的关闭查询，在加互斥锁的情况下返回 closed 的值
// Conn 通过调用 UniqueSymbol() string 向外返回一个独一无二的标识，这个标识用于 Item 基于此值构造内部唯一标识
package network

import (
	"errors"
	"github.com/gorilla/websocket"
	"jarvis/base/network/grpc"
	"jarvis/util/rand"
	"net"
	"sync"
)

type (
	// 连接定义
	Conn interface {
		// 读取，一次性阻塞读取
		Read() ([]byte, error)

		// 写入
		Write([]byte) error

		// 关闭，重复关闭会报错
		Close() error

		// 询问是否已关闭，已处理与 Close() error 函数存在的 closed 多线程竞态问题
		IsClosed() bool

		// 唯一标识
		UniqueSymbol() string
	}

	// 基础通用结构
	baseConn struct {
		closed bool // 是否已关闭
	}

	// socket 连接实现
	socketConn struct {
		baseConn
		c     net.Conn   // 底层连接
		mutex sync.Mutex // 对关闭状态的多线程竞态加锁
	}

	// webSocket 连接实现
	webSocketConn struct {
		baseConn
		c          *websocket.Conn // 底层连接
		mutex      sync.Mutex      // 对关闭状态的多线程竞态加锁
		writeMutex sync.Mutex      // 当前使用的 "github.com/gorilla/websocket" 库不支持并发写，因此加入锁
	}

	// gRPC 连接接口实现结构
	gRPCConn struct {
		baseConn
		ccs       grpc.Communicate_ConnectServer // gRPC server 实例
		mutex     sync.Mutex                     // 对关闭状态的多线程竞态加锁
		closeChan chan struct{}                  // 关闭通知通道
	}
)

const (
	// 缓存默认大小，默认为 1024 Byte , 即 1KB
	BufferLength = 1024
)

// 此常量组定义了 Conn 定义及实现中可能会发生的错误文本
const (
	ErrNilConnText    = "connection is nil"
	ErrConnClosedText = "connection already closed"
)

var (
	// 连接为空错误，此错误下建议不使用任何 Conn 函数
	ErrNilConn = errors.New(ErrNilConnText)
	// 连接已关闭错误，此错误下建议不使用任何 Conn 函数
	ErrConnClosed = errors.New(ErrConnClosedText)
)

// 实例化 socketConn 的 Conn 实现，返回表现形式为 Conn
func NewSocketConn(c net.Conn) Conn {
	return &socketConn{
		baseConn: baseConn{},
		c:        c,
		mutex:    sync.Mutex{},
	}
}

// 实例化 webSocketConn 的 Conn 实现，返回表现形式为 Conn
func NewWebSocketConn(c *websocket.Conn) Conn {
	return &webSocketConn{
		baseConn:   baseConn{},
		c:          c,
		mutex:      sync.Mutex{},
		writeMutex: sync.Mutex{},
	}
}

// 实例化 gRPCConn 的 Conn 实现，返回表现形式为 Conn
func NewGRPCConn(ccs grpc.Communicate_ConnectServer) (Conn, chan struct{}) {
	closeChan := make(chan struct{})

	return &gRPCConn{
		baseConn:  baseConn{},
		ccs:       ccs,
		closeChan: closeChan,
		mutex:     sync.Mutex{},
	}, closeChan
}

// --------------------------------------------------- socketConn ------------------------------------------------------
// 读取，一次性阻塞读取
// 当 err != nil 时，[]byte 为 nil
func (sc *socketConn) Read() ([]byte, error) {
	if sc.c == nil {
		return nil, ErrNilConn
	}
	if sc.closed {
		return nil, ErrConnClosed
	}

	buffer := make([]byte, BufferLength)
	length, err := sc.c.Read(buffer)
	if err != nil {
		return nil, err
	}

	return buffer[:length], nil
}

// 写入
func (sc *socketConn) Write(data []byte) error {
	if sc.c == nil {
		return ErrNilConn
	}
	if sc.closed {
		return ErrConnClosed
	}

	_, err := sc.c.Write(data)

	return err
}

// 关闭，重复关闭会报错
func (sc *socketConn) Close() error {
	if sc.c == nil {
		return ErrNilConn
	}
	if sc.closed {
		return ErrConnClosed
	}

	sc.mutex.Lock()
	sc.closed = true
	sc.mutex.Unlock()

	return sc.c.Close()
}

// 询问是否已关闭
func (sc *socketConn) IsClosed() bool {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	return sc.closed
}

// 唯一标识
func (sc *socketConn) UniqueSymbol() string {
	if sc.c == nil || sc.closed {
		return ""
	}

	return sc.c.RemoteAddr().String()
}

// --------------------------------------------------- webSocketConn ---------------------------------------------------
// 读取，一次性阻塞读取
// 当 err != nil 时，[]byte 为 nil
func (wsc *webSocketConn) Read() ([]byte, error) {
	if wsc.c == nil {
		return nil, ErrNilConn
	}
	if wsc.closed {
		return nil, ErrConnClosed
	}

	_, data, err := wsc.c.ReadMessage()
	if err != nil {
		return nil, err
	}

	return data, nil
}

// 写入
func (wsc *webSocketConn) Write(data []byte) error {
	if wsc.c == nil {
		return ErrNilConn
	}
	if wsc.closed {
		return ErrConnClosed
	}
	wsc.writeMutex.Lock()
	defer wsc.writeMutex.Unlock()

	return wsc.c.WriteMessage(websocket.TextMessage, data)
}

// 关闭，重复关闭会报错
func (wsc *webSocketConn) Close() error {
	if wsc.c == nil {
		return ErrNilConn
	}
	if wsc.closed {
		return ErrConnClosed
	}

	wsc.mutex.Lock()
	wsc.closed = true
	wsc.mutex.Unlock()

	return wsc.c.Close()
}

// 询问是否已关闭
func (wsc *webSocketConn) IsClosed() bool {
	wsc.mutex.Lock()
	defer wsc.mutex.Unlock()
	return wsc.closed
}

// 唯一标识
func (wsc *webSocketConn) UniqueSymbol() string {
	if wsc.c == nil || wsc.closed {
		return ""
	}

	return wsc.c.RemoteAddr().String()
}

// --------------------------------------------------- gRPCConn --------------------------------------------------------
// 读取，一次性阻塞读取
// 当 err != nil 时，[]byte 为 nil
func (gc *gRPCConn) Read() ([]byte, error) {
	if gc.ccs == nil {
		return nil, ErrNilConn
	}
	if gc.closed {
		return nil, ErrConnClosed
	}

	message, err := gc.ccs.Recv()
	if err != nil {
		return nil, err
	}

	return message.Data, nil
}

// 写入
func (gc *gRPCConn) Write(data []byte) error {
	if gc.ccs == nil {
		return ErrNilConn
	}
	if gc.closed {
		return ErrConnClosed
	}

	return gc.ccs.Send(&grpc.Message{Data: data})
}

// 关闭，重复关闭会报错
func (gc *gRPCConn) Close() error {
	if gc.ccs == nil {
		return ErrNilConn
	}
	if gc.closed {
		return ErrConnClosed
	}

	gc.mutex.Lock()
	gc.closed = true
	gc.mutex.Unlock()

	// 向 gate 报告连接关闭
	gc.closeChan <- struct{}{}

	return nil
}

// 询问是否已关闭
func (gc *gRPCConn) IsClosed() bool {
	gc.mutex.Lock()
	defer gc.mutex.Unlock()
	return gc.closed
}

// 唯一标识
func (gc *gRPCConn) UniqueSymbol() string {
	if gc.ccs == nil || gc.closed {
		return ""
	}

	return rand.RandomString(8)
}
