// Item 是对 Conn 的业务包装，统一负责读取、写入、关闭、断开逆反馈
// Item 在 Conn.IsClosed() == false 的情况下，于单个线程内阻塞式读取消息，并将读取的字节组输入到上层统一的 Packager 中进行解包，解密
// Item 在写入数据时，通过 Packager 进行加密、打包成字节组，再调用 Conn 的 Write() 函数进行发送到客户端
// Item 共有5种状态，创建、运行、主动关闭、被动关闭、未知关闭，被动关闭和未知关闭状态下会调用一次 Close() 关闭 Conn 并上报反馈，服务端主动
// 关闭 Item 的情况下，直接调用 Item.Close() 也会触发上报反馈
// Item 默认 Hook 了 上层 Manager 的 RemoveItem() 函数，因此 Close() 的时候会调用此函数将自己从管理中移除
package network

import (
	"jarvis/base/log"
	"strings"
)

type (
	// 端 状态
	ItemState int

	// 客户端断开反馈函数签名
	PassiveCloseFeedbackFunc func(string) error

	// 端定义
	Item interface {
		// 获取唯一标识
		ID() ID

		// 接收消息
		Receive(chan<- Message)

		// 发送消息
		Send(Message)

		// 关闭
		Close()

		// 客户端断开反馈
		PassiveCloseFeedback(PassiveCloseFeedbackFunc)
	}

	// 端定义实现
	item struct {
		id        ID                       // 内部唯一标识
		conn      Conn                     // 连接
		state     ItemState                // 状态
		packager  Packager                 // 装包者
		encrypter Encrypter                // 加密器
		FbFunc    PassiveCloseFeedbackFunc // 客户端断开反馈函数
	}
)

// 此常量组定义了 Item 存在的状态
const (
	ItemStateCreate          ItemState = iota // 创建
	ItemStateRunning                          // 运行
	ItemStateInitiativeClose                  // 服务端主动关闭
	ItemStatePassiveClose                     // 服务端被动关闭，即客户端主动断开
	ItemStateUnKnowClose                      // 服务端接收到未知错误，主动关闭
)

// 此常量组定义了 Item 定义及实现中可能会发生的错误文本
const (
	// 服务端主动关闭，导致使用了已经关闭了的连接
	ErrNetClosingText = "use of closed network connection"
	// 客户端主动断开，表现为 EOF
	EOFText = "EOF"
	// gRPC 客户端主动断开，表现为 上下文取消
	ContextCancelText = "context canceled"
)

// 新建端
func NewItem(conn Conn, packager Packager, encrypter Encrypter) Item {
	id := ID(conn.UniqueSymbol()) // 对 Conn 的唯一标识进行包装

	return &item{
		id:        EncryptID(id), // 加密
		conn:      conn,
		state:     ItemStateCreate,
		packager:  packager,
		encrypter: encrypter,
		FbFunc:    nil,
	}
}

// 端状态可读化
func (is ItemState) String() string {
	switch is {
	case ItemStateCreate:
		return "Create"
	case ItemStateRunning:
		return "Running"
	case ItemStateInitiativeClose:
		return "InitiativeClose"
	case ItemStatePassiveClose:
		return "PassiveClose"
	case ItemStateUnKnowClose:
		return "UnKnowClose"
	default:
		return "UnKnowState"
	}
}

// 获取唯一标识
func (i *item) ID() ID {
	return i.id
}

// 接收消息
func (i *item) Receive(channel chan<- Message) {
	// 运行状态
	i.state = ItemStateRunning
	//log.Printf("[%s] start receive-%s", i.ID().String(), i.state.String())

	var e error
	// 监听 conn 的关闭状态
	for !i.conn.IsClosed() {
		b, err := i.conn.Read()
		if err != nil {
			e = err
			break
		}

		// 解包数据，反序列化到 BaseRequest 结构中，附带上内部唯一标识，发送到 Service 的请求消息流 channel 中
		for _, data := range i.packager.Unpack(b) {
			request := Message{}
			if err := request.Unmarshal(i.encrypter.Decrypt(data)); err != nil {
				log.ErrorF("unmarshal data to BaseRequest error : %s", err.Error())
				continue
			}

			request.ID = i.ID().String()

			go func() {
				if channel != nil {
					channel <- request
				}
			}()
		}
	}

	// 如果跳出了该读取循环，通过 e(error) 的值可以判断当前 Item 的关闭状态，根据不同的状态进行关闭处理
	if e != nil {
		if strings.Contains(e.Error(), EOFText) || strings.Contains(e.Error(), ContextCancelText) { // 客户端主动断开
			i.state = ItemStatePassiveClose
		} else if strings.Contains(e.Error(), ErrNetClosingText) { // 服务端主动断开
			i.state = ItemStateInitiativeClose
		} else { // 未知错误
			i.state = ItemStateUnKnowClose
			log.ErrorF("[%s] receive error : %s", i.ID().String(), e.Error())
		}
	}

	// 客户端主动断开，服务端再调用一次 Close()
	// 服务端主动关闭，会直接调用 Close()
	// 捕获到未知错误，主动断开
	if i.state == ItemStatePassiveClose || i.state == ItemStateUnKnowClose {
		i.Close()
	}
}

// 发送消息
func (i *item) Send(response Message) {
	// 将 BaseResponse 序列化
	data, err := response.Marshal()
	if err != nil {
		log.ErrorF("[%s] unmarshal response error : %s", i.ID().String(), err.Error())
		return
	}

	// 通过装包者打包，往 Item 持有的 Conn 中写入
	if err := i.conn.Write(i.packager.Pack(i.encrypter.Encrypt(data))); err != nil {
		i.state = ItemStateUnKnowClose
		i.Close()
		log.ErrorF("[%s] send error : %s", i.ID().String(), err.Error())
	}
}

// 关闭
func (i *item) Close() {
	if err := i.conn.Close(); err != nil {
		log.ErrorF("[%s] close error : %s", i.ID().String(), err.Error())
	}
	// 断开反馈
	// 向上反馈给 Service 持有的 Manager ， Manager 会通过钩子函数向业务反馈
	if i.FbFunc != nil {
		if err := i.FbFunc(i.ID().String()); err != nil {
			log.ErrorF("[%s] - [%s] close feedback error -%s", i.ID().String(), i.state.String(), err.Error())
		}
	}
	//log.Printf("[%s] closed-%s", i.ID().String(), i.state.String())
}

// 客户端断开反馈
func (i *item) PassiveCloseFeedback(function PassiveCloseFeedbackFunc) {
	// 持有断开反馈钩子函数
	i.FbFunc = function
}
