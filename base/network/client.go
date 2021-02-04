package network

import (
	oContext "context"
	"errors"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	gRPC "jarvis/base/network/grpc"
	"log"
	"net"
	"sync"
	"time"
)

type (
	// 客户端定义
	Client interface {
		// 初始化
		Initialize() error

		// 注册路由
		RegisterReply(string, func(Message)) error

		// 接收
		Receive() (Message, error)

		// 发送
		Send(Message) error

		// 关闭
		Close() error

		// 同步请求
		RequestSync(Message) (Message, error)
	}

	// 基础客户端结构
	baseClient struct {
		routeMap    map[string]chan Message // 请求
		mutex       sync.Mutex
		receiveChan chan Message
		address     string
		p           Packager
		e           Encrypter
		closed      bool
	}

	// socket 客户端结构
	socketClient struct {
		baseClient
		c Conn // 连接
	}

	// webSocket 客户端结构
	webSocketClient struct {
		baseClient
		c Conn // 连接
	}

	// webSocket 客户端结构
	gRPCClient struct {
		baseClient
		cc  *grpc.ClientConn
		ccc gRPC.Communicate_ConnectClient
	}
)

const (
	// 默认超时时间
	DefaultTimeout = time.Second * time.Duration(15)
)

const (
	ErrNilChannelText           = "channel is nil"
	ErrReceiveChannelClosedText = "receive channel closed"
	ErrClosedConnectionText     = "use of closed network connection"
	ErrClientAlreadyClosedText  = " already closed"
	ErrTimeoutText              = "timeout"
)

var (
	// channel 为 nil 错误
	ErrNilChannel = errors.New(ErrNilChannelText)
	// 接收管道已关闭 错误
	ErrReceiveChannelClosed = errors.New(ErrReceiveChannelClosedText)
	// 客户端已经关闭 错误
	ErrClientAlreadyClosed = errors.New(ErrClientAlreadyClosedText)
	// 超时 错误
	ErrTimeout = errors.New(ErrTimeoutText)
)

// 新建基础客户端
func newBaseClient(address string, packager Packager, encrypter Encrypter) baseClient {
	return baseClient{
		routeMap:    make(map[string]chan Message),
		mutex:       sync.Mutex{},
		receiveChan: make(chan Message),
		address:     address,
		p:           packager,
		e:           encrypter,
		closed:      false,
	}
}

// 新建 Socket 客户端
func NewSocketClient(address string, packager Packager, encrypter Encrypter) Client {
	return &socketClient{
		baseClient: newBaseClient(address, packager, encrypter),
		c:          nil,
	}
}

// 新建 WebSocket 客户端
func NewWebSocketClient(address string, packager Packager, encrypter Encrypter) Client {
	return &webSocketClient{
		baseClient: newBaseClient(address, packager, encrypter),
		c:          nil,
	}
}

// 新建 WebSocket 客户端
func NewGRPCClient(address string, packager Packager, encrypter Encrypter) Client {
	return &gRPCClient{
		baseClient: newBaseClient(address, packager, encrypter),
		cc:         nil,
		ccc:        nil,
	}
}

// -------------------------------------------------- Base Client ------------------------------------------------------
// 添加一条临时路由
func (bc *baseClient) AddRoute(route string, channel chan Message) error {
	if bc.closed {
		return ErrClientAlreadyClosed
	}
	if route == "" {
		return ErrEmptyRouteName
	}
	if channel == nil {
		return ErrNilChannel
	}

	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	if _, exist := bc.routeMap[route]; exist {
		return ErrRouteExist
	}

	bc.routeMap[route] = channel
	return nil
}

// 移除一条临时路由
func (bc *baseClient) RemoveRoute(route string) error {
	if bc.closed {
		return ErrClientAlreadyClosed
	}
	if route == "" {
		return ErrEmptyRouteName
	}

	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	if _, exist := bc.routeMap[route]; !exist {
		return ErrRouteUnExist
	}

	delete(bc.routeMap, route)

	return nil
}

// 将消息路由到对应的回复管道，并移除记录
func (bc *baseClient) Route(message Message) error {
	if bc.closed {
		return ErrClientAlreadyClosed
	}
	if message.Reply == "" {
		return ErrEmptyRouteName
	}

	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	if _, exist := bc.routeMap[message.Reply]; !exist {
		return ErrRouteUnExist
	}

	channel := bc.routeMap[message.Reply]
	if channel == nil {
		return nil
	}

	channel <- message

	return nil
}

// -------------------------------------------------- Socket Client ----------------------------------------------------
// 初始化
func (sc *socketClient) Initialize() error {
	if sc.baseClient.closed {
		return ErrClientAlreadyClosed
	}
	c, err := net.Dial(DefaultNetwork, sc.baseClient.address)
	if err != nil {
		return err
	}

	sc.c = NewSocketConn(c)

	go sc.run()

	return nil
}

func (sc *socketClient) RegisterReply(route string, function func(Message)) error {
	c := make(chan Message)
	go func(channel chan Message) {
		msg := <-channel
		function(msg)
	}(c)

	return sc.baseClient.AddRoute(route, c)
}

// 接收
func (sc *socketClient) Receive() (Message, error) {
	if sc.baseClient.closed {
		return Message{}, ErrClientAlreadyClosed
	}
	message, ok := <-sc.baseClient.receiveChan
	if !ok {
		return Message{}, ErrReceiveChannelClosed
	}

	return message, nil
}

// 发送
func (sc *socketClient) Send(message Message) error {
	if sc.baseClient.closed {
		return ErrClientAlreadyClosed
	}
	data, err := message.Marshal()
	if err != nil {
		return err
	}

	return sc.c.Write(sc.baseClient.p.Pack(sc.baseClient.e.Encrypt(data)))
}

// 关闭
func (sc *socketClient) Close() error {
	if sc.baseClient.closed {
		return ErrClientAlreadyClosed
	}
	sc.baseClient.closed = true
	close(sc.baseClient.receiveChan)
	sc.baseClient.receiveChan = nil

	return sc.c.Close()
}

// 同步请求
func (sc *socketClient) RequestSync(request Message) (Message, error) {
	if sc.baseClient.closed {
		return Message{}, ErrClientAlreadyClosed
	}

	responseChannel := make(chan Message)

	if err := sc.baseClient.AddRoute(request.Reply, responseChannel); err != nil {
		return Message{}, err
	}

	if err := sc.Send(request); err != nil {
		_ = sc.baseClient.RemoveRoute(request.Reply)
		return Message{}, err
	}

	select {
	case response := <-responseChannel:
		_ = sc.baseClient.RemoveRoute(request.Reply)
		return response, nil
	case <-time.After(DefaultTimeout):
		_ = sc.baseClient.RemoveRoute(request.Reply)
		return Message{}, ErrTimeout
	}
}

// 监听消息且转发
func (sc *socketClient) run() {
	var e error
	for {
		d, err := sc.c.Read()
		if err != nil {
			e = err
			if sc.baseClient.closed {
				e = nil
			}
			break
		}

		for _, data := range sc.baseClient.p.Unpack(d) {
			response := Message{}
			if err := response.Unmarshal(sc.baseClient.e.Decrypt(data)); err != nil {
				log.Printf("Socket unmarshal data error : %s", err.Error())
				continue
			}

			if response.Reply == "" {
				sc.baseClient.receiveChan <- response
			} else {
				if err := sc.baseClient.Route(response); err != nil {
					log.Printf("Socket route response error : %s", err.Error())
				}
			}
		}
	}

	if e != nil {
		log.Printf("Socket  read error : %s", e.Error())

		if err := sc.Close(); err != nil {
			log.Printf("socket  close error : %s", err.Error())
		}
	}
}

// -------------------------------------------------- WebSocket Client -------------------------------------------------
// 初始化
func (wsc *webSocketClient) Initialize() error {
	if wsc.baseClient.closed {
		return ErrClientAlreadyClosed
	}

	c, _, err := websocket.DefaultDialer.Dial(wsc.baseClient.address, nil)
	if err != nil {
		return err
	}

	wsc.c = NewWebSocketConn(c)

	go wsc.run()

	return nil
}

func (wsc *webSocketClient) RegisterReply(route string, function func(Message)) error {
	c := make(chan Message)
	go func(channel chan Message) {
		msg := <-channel
		function(msg)
	}(c)

	return wsc.baseClient.AddRoute(route, c)
}

// 接收
func (wsc *webSocketClient) Receive() (Message, error) {
	if wsc.baseClient.closed {
		return Message{}, ErrClientAlreadyClosed
	}
	message, ok := <-wsc.baseClient.receiveChan
	if !ok {
		return Message{}, ErrReceiveChannelClosed
	}

	return message, nil
}

// 发送
func (wsc *webSocketClient) Send(message Message) error {
	if wsc.baseClient.closed {
		return ErrClientAlreadyClosed
	}
	data, err := message.Marshal()
	if err != nil {
		return err
	}

	return wsc.c.Write(wsc.baseClient.p.Pack(wsc.baseClient.e.Encrypt(data)))
}

// 关闭
func (wsc *webSocketClient) Close() error {
	if wsc.baseClient.closed {
		return ErrClientAlreadyClosed
	}
	wsc.baseClient.closed = true
	close(wsc.baseClient.receiveChan)
	wsc.baseClient.receiveChan = nil

	return wsc.c.Close()
}

// 同步请求
func (wsc *webSocketClient) RequestSync(request Message) (Message, error) {
	if wsc.baseClient.closed {
		return Message{}, ErrClientAlreadyClosed
	}

	responseChannel := make(chan Message)

	if err := wsc.baseClient.AddRoute(request.Reply, responseChannel); err != nil {
		return Message{}, err
	}

	if err := wsc.Send(request); err != nil {
		_ = wsc.baseClient.RemoveRoute(request.Reply)
		return Message{}, err
	}

	select {
	case response := <-responseChannel:
		_ = wsc.baseClient.RemoveRoute(request.Reply)
		return response, nil
	case <-time.After(DefaultTimeout):
		_ = wsc.baseClient.RemoveRoute(request.Reply)
		return Message{}, ErrTimeout
	}
}

// 监听消息且转发
func (wsc *webSocketClient) run() {
	var e error
	for {
		d, err := wsc.c.Read()
		if err != nil {
			e = err
			if wsc.baseClient.closed {
				e = nil
			}
			break
		}

		for _, data := range wsc.baseClient.p.Unpack(d) {
			response := Message{}
			if err := response.Unmarshal(wsc.baseClient.e.Decrypt(data)); err != nil {
				log.Printf("WebSocket unmarshal data error : %s", err.Error())
				continue
			}

			if response.Reply == "" {
				wsc.baseClient.receiveChan <- response
			} else {
				if err := wsc.baseClient.Route(response); err != nil {
					log.Printf("%+v", response)
					log.Printf("Socket route response error : %s", err.Error())
				}
			}
		}
	}

	if e != nil {
		log.Printf("Socket  read error : %s", e.Error())

		if err := wsc.Close(); err != nil {
			log.Printf("socket  close error : %s", err.Error())
		}
	}
}

// -------------------------------------------------- gRPC Client ------------------------------------------------------
// 初始化
func (gc *gRPCClient) Initialize() error {
	if gc.baseClient.closed {
		return ErrClientAlreadyClosed
	}

	cc, err := grpc.Dial(gc.baseClient.address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}

	cClient := gRPC.NewCommunicateClient(cc)
	ccc, err := cClient.Connect(oContext.Background())
	if err != nil {
		return err
	}

	gc.cc = cc
	gc.ccc = ccc

	go gc.run()

	return nil
}

func (gc *gRPCClient) RegisterReply(route string, function func(Message)) error {
	c := make(chan Message)
	go func(channel chan Message, f func(Message)) {
		msg := <-channel
		f(msg)
	}(c, function)

	return gc.baseClient.AddRoute(route, c)
}

// 接收
func (gc *gRPCClient) Receive() (Message, error) {
	if gc.baseClient.closed {
		return Message{}, ErrClientAlreadyClosed
	}
	message, ok := <-gc.baseClient.receiveChan
	if !ok {
		return Message{}, ErrReceiveChannelClosed
	}

	return message, nil
}

// 发送
func (gc *gRPCClient) Send(message Message) error {
	if gc.baseClient.closed {
		return ErrClientAlreadyClosed
	}
	data, err := message.Marshal()
	if err != nil {
		return err
	}

	return gc.ccc.Send(&gRPC.Message{Data: gc.baseClient.p.Pack(gc.baseClient.e.Encrypt(data))})
}

// 关闭
func (gc *gRPCClient) Close() error {
	if gc.baseClient.closed {
		return ErrClientAlreadyClosed
	}
	gc.baseClient.closed = true
	close(gc.baseClient.receiveChan)
	gc.baseClient.receiveChan = nil

	if err := gc.cc.Close(); err != nil {
		return err
	}

	return gc.ccc.CloseSend()
}

// 同步请求
func (gc *gRPCClient) RequestSync(request Message) (Message, error) {
	if gc.baseClient.closed {
		return Message{}, ErrClientAlreadyClosed
	}

	responseChannel := make(chan Message)

	if err := gc.baseClient.AddRoute(request.Reply, responseChannel); err != nil {
		return Message{}, err
	}

	if err := gc.Send(request); err != nil {
		_ = gc.baseClient.RemoveRoute(request.Reply)
		return Message{}, err
	}

	select {
	case response := <-responseChannel:
		_ = gc.baseClient.RemoveRoute(request.Reply)
		return response, nil
	case <-time.After(DefaultTimeout):
		_ = gc.baseClient.RemoveRoute(request.Reply)
		return Message{}, ErrTimeout
	}
}

// 监听消息且转发
func (gc *gRPCClient) run() {
	var e error
	for {
		d, err := gc.ccc.Recv()
		if err != nil {
			e = err
			if gc.baseClient.closed {
				e = nil
			}
			break
		}

		for _, data := range gc.baseClient.p.Unpack(d.Data) {
			response := Message{}
			if err := response.Unmarshal(gc.baseClient.e.Decrypt(data)); err != nil {
				log.Printf("gRPC unmarshal data error : %s", err.Error())
				continue
			}

			if response.Reply == "" {
				gc.baseClient.receiveChan <- response
			} else {
				if err := gc.baseClient.Route(response); err != nil {
					log.Printf("%+v", response)
					log.Printf("Socket route response error : %s", err.Error())
				}
			}
		}
	}

	if e != nil {
		log.Printf("Socket  read error : %s", e.Error())

		if err := gc.Close(); err != nil {
			log.Printf("socket  close error : %s", err.Error())
		}
	}
}
