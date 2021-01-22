package network

import "encoding/json"

type (
	// 框架要求客户端和服务端发送的消息在装包者打包之前必须为 Message 结构
	Message struct {
		Module string `json:"module"` // 模块名
		Route  string `json:"route"`  // 路由名
		ID     string `json:"-"`      // 内部唯一标识，仅由 Item 带出，最终回到 Item ，用于标识消息的来源和去处
		Data   []byte `json:"data"`   // 具体数据，业务所需要的数据序列化后装入此项发送，由另一端按照约定结构反序列化
		Reply  string `json:"reply"`  // 回覆，当客户端不处于双工通信状态时，可以进行一次性请求，此项用于客户端收到回复后的消息转发
	}

	// 框架[建议]在返回时封装响应，便于客户端统一响应
	Reply struct {
		Code    int    `json:"code"`    // 状态码
		Message string `json:"message"` // 状态语句
		Data    []byte `json:"data"`    // 数据
	}
)

const (
	// 成功响应
	ReplySuccessCode = 200
	// 请求不正确
	ReplyBadRequestCode = 400
	// 服务器错误
	ReplyServerErrorCode = 500
)

var (
	// 成功响应
	ReplySuccessMessage = "Success"
	// 请求不正确
	ReplyBadRequestMessage = "Bad request"
	// 服务器错误
	ReplyServerErrorMessage = "Server error"
)

// 反序列化
func (m *Message) Unmarshal(b []byte) error {
	return json.Unmarshal(b, m)
}

// 序列化
func (m *Message) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

func ReplySuccess(data []byte) Reply {
	return Reply{
		Code:    ReplySuccessCode,
		Message: ReplySuccessMessage,
		Data:    data,
	}
}

func ReplyBadRequest(data []byte) Reply {
	return Reply{
		Code:    ReplyBadRequestCode,
		Message: ReplyBadRequestMessage,
		Data:    data,
	}
}

func ReplyServerError(data []byte) Reply {
	return Reply{
		Code:    ReplyServerErrorCode,
		Message: ReplyServerErrorMessage,
		Data:    data,
	}
}
