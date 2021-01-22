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
)

// 反序列化
func (m *Message) Unmarshal(b []byte) error {
	return json.Unmarshal(b, m)
}

// 序列化
func (m *Message) Marshal() ([]byte, error) {
	return json.Marshal(m)
}
