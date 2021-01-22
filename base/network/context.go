// context 被设计用于在 BaseRequest 的路由分发上，装载了 BaseRequest 和 *BaseResponse ，所有的 RouteHandleFunc 函数只接收一个
// context 的接口(即为实现的指针)，context 可以读取请求和发送多次响应，调用 Done() 函数即在调用链中截止
package network

import "errors"

type (
	// 上下文定义
	Context interface {
		// 获取请求
		Request() Message

		// 响应
		Reply(string, string, []byte) error

		// 错误响应
		Error(error) error

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
	}
}

// 获取请求
func (c *context) Request() Message {
	// 复制持有的请求
	request := c.request
	return request
}

// 响应
func (c *context) Reply(id string, reply string, data []byte) error {
	return c.findAndSendReply(id, reply, data)
}

// 错误响应
func (c *context) Error(err error) error {
	return c.findAndSendReply(c.request.ID, c.request.Reply, []byte(err.Error()))
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
func (c *context) findAndSendReply(id, reply string, data []byte) error {
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
