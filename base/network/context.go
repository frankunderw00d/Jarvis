// context 被设计用于在 BaseRequest 的路由分发上，装载了 BaseRequest 和 *BaseResponse ，所有的 RouteHandleFunc 函数只接收一个
// context 的接口(即为实现的指针)，context 可以读取请求和发送多次响应，调用 Done() 函数即在调用链中截止
package network

import (
	"encoding/json"
	"errors"
)

type (
	// 上下文定义
	Context interface {
		// 获取请求
		Request() Message

		// 设置额外信息
		SetExtra(string, interface{})

		// 获取额外信息
		Extra(string, interface{}) interface{}

		// 根据 id 查询和调用 Item 的 Send() 函数
		FindAndSendReply(string, string, []byte) error

		// 原始数据返回，不经过 Reply 结构包含，常用于跨服务调用时，中间曾转发
		BinaryReply(d []byte) error

		// 请求成功回复
		Success([]byte) error

		// 请求错误回复
		BadRequest(string) error

		// 服务器错误回复
		ServerError(error) error

		// 钩住查询 Item 的函数
		HookFind(func(string) (Item, error))

		// 结束调用链传递
		Done()

		// 是否已经结束
		IsDone() bool
	}

	// 上下文定义实现
	context struct {
		request  Message                    // 请求
		done     bool                       // 是否中断调用链
		findFunc func(string) (Item, error) // 寻找响应调用
		extra    map[string]interface{}     // 额外附带信息
	}
)

// 此常量组定义了 Context 定义及实现中可能会发生的错误文本
const (
	ErrNilFindFuncText = "find function is nil"
)

// 此变量组定义了 Context 定义及实现中可能会发生的错误
var (
	// 寻找 Item 函数为 nil 错误
	ErrNilFindFunc = errors.New(ErrNilFindFuncText)
)

func NewContext(request Message) Context {
	return &context{
		request:  request,
		done:     false, // 默认未结束调用链
		findFunc: nil,   // 默认不持有任何查找 Item 函数
		extra:    map[string]interface{}{},
	}
}

// 获取请求
func (c *context) Request() Message {
	// 复制持有的请求
	request := c.request
	return request
}

// 设置额外信息
func (c *context) SetExtra(key string, value interface{}) {
	c.extra[key] = value
}

// 获取额外信息
func (c *context) Extra(key string, defaultValue interface{}) interface{} {
	v, exist := c.extra[key]
	if !exist {
		return defaultValue
	}
	return v
}

// 钩住查询 Item 的函数
func (c *context) HookFind(findFunc func(string) (Item, error)) {
	if findFunc != nil {
		c.findFunc = findFunc
	}
}

// 结束调用链传递
func (c *context) Done() {
	c.done = true
}

// 是否已经结束
func (c *context) IsDone() bool {
	return c.done
}

// 根据 id 查询和调用 Item 的 Send() 函数
func (c *context) FindAndSendReply(id, reply string, data []byte) error {
	if id == "" {
		return ErrNilId
	}
	if c.findFunc == nil {
		return ErrNilFindFunc
	}

	// 查找 Item
	i, err := c.findFunc(id)
	if err != nil {
		return err
	}

	// 发送消息
	i.Send(Message{
		ID:    id,
		Data:  data,
		Reply: reply,
	})

	return nil
}

// 原始数据返回，不经过 Reply 结构包含，常用于跨服务调用时，中间曾转发
func (c *context) BinaryReply(d []byte) error {
	return c.FindAndSendReply(c.request.ID, c.request.Reply, d)
}

// 请求成功回复
func (c *context) Success(d []byte) error {
	reply := ReplySuccess(d)
	data, err := json.Marshal(&reply)
	if err != nil {
		return err
	}

	return c.FindAndSendReply(c.request.ID, c.request.Reply, data)
}

// 请求错误回复
func (c *context) BadRequest(str string) error {
	reply := ReplyBadRequest([]byte(str))
	data, err := json.Marshal(&reply)
	if err != nil {
		return err
	}
	return c.FindAndSendReply(c.request.ID, c.request.Reply, data)
}

// 服务器错误回复
func (c *context) ServerError(e error) error {
	reply := ReplyServerError([]byte(e.Error()))
	data, err := json.Marshal(&reply)
	if err != nil {
		return err
	}
	return c.FindAndSendReply(c.request.ID, c.request.Reply, data)
}
