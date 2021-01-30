package network

type (
	// 观察者定义
	Observer interface {
		// 观察连接
		ObserveConnect(string)

		// 观察断开连接
		ObserveDisconnect(string)

		// 观察者主动发送至端 ， 此函数会被放置于一个独立的线程执行，因此必须阻塞，否则只调用一次
		InitiativeSend(Context)
	}
)

const ()

var ()
