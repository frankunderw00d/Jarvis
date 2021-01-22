package network

type (
	// 观察者定义
	Observer interface {
		// 观察连接
		ObserveConnect(string)

		// 观察断开连接
		ObserveDisconnect(string)
	}
)

const ()

var ()
