package network

type (
	// 调用链表定义
	CallLinkedList interface {
		// 添加一个调用节点
		AddNode(RouteHandleFunc)

		// 执行调用链
		Run(Context)

		// 复制当前调用链
		Copy() CallLinkedList
	}

	// 调用链表节点定义
	node struct {
		handleFunc RouteHandleFunc // 当前节点的调用函数
		next       *node           // 下一个节点
	}

	// 调用链表定义实现
	callLinkedList struct {
		rootNode *node // 根节点
	}
)

const ()

var ()

// 新建 节点
func newNode(handleFunc RouteHandleFunc) *node {
	return &node{
		handleFunc: handleFunc, // 持有传入调用函数
		next:       nil,        // 新建节点默认下一个节点为 nil
	}
}

// 新建调用链
func NewCallLinkedList() CallLinkedList {
	return &callLinkedList{}
}

// 添加一个调用节点
func (cll *callLinkedList) AddNode(handleFunc RouteHandleFunc) {
	// 如果根节点为空，新节点为根节点
	if cll.rootNode == nil {
		cll.rootNode = newNode(handleFunc)
		return
	}

	// 否则获取当前调用链的最后一个节点，将新节点置为其下一个节点
	cll.lastNode().next = newNode(handleFunc)
}

// 执行调用链
func (cll *callLinkedList) Run(ctx Context) {
	currentNode := cll.rootNode

	// 从根节点开始，逐步向下调用
	// 无下一个节点 或 中断 则停止调用
	for currentNode != nil && !ctx.IsDone() {
		currentNode.handleFunc(ctx)
		currentNode = currentNode.next
	}
}

// 复制当前调用链
func (cll *callLinkedList) Copy() CallLinkedList {
	another := &callLinkedList{}

	// 从根节点开始，复制当前调用链的所有节点
	currentNode := cll.rootNode
	for currentNode != nil {
		another.AddNode(currentNode.handleFunc)
		currentNode = currentNode.next
	}

	return another
}

// 获取到最后一个节点
func (cll *callLinkedList) lastNode() *node {
	lastNode := cll.rootNode

	// 从根节点开始，遍历调用链到最后一个节点
	for lastNode.next != nil {
		lastNode = lastNode.next
	}

	return lastNode
}
