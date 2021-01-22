// Manager 是 Item 管理者的角色，对外提供函数供外部调用来管理由 Gate 下放的 Conn 构建成的 Item
package network

import (
	"errors"
	"sync"
)

type (
	// 端管理者定义
	Manager interface {
		// 添加一个端到当前管理
		ManageItem(Item) error

		// 从当前管理移除指定端
		RemoveItem(string) error

		// 从当前管理寻找指定 id 的端
		FindItem(id string) (Item, error)

		// 目前管理端数目
		CurrentManageCount() int64

		// 注册观察者
		RegisterObserver(Observer)
	}

	// 端管理者定义实现
	manager struct {
		max       int64           // 最大管理连接
		mutex     sync.Mutex      // 端管理竞态锁
		items     map[string]Item // 端管理
		observers []Observer      // 观察者列表
	}
)

const (
	DefaultMaxConnection = 5000 // 默认最大端管理数
)

// 此常量组定义了 Manager 定义及实现中可能会发生的错误文本
const (
	ErrNilItemText     = "item is nil"
	ErrNilIdText       = "id is nil"
	ErrMaxConnectText  = "already reach max connect"
	ErrItemExistText   = "the id of item already exist"
	ErrItemUnExistText = "the id of item doesn't exist"
)

// 此常量组定义了 Manager 定义及实现中可能会发生的错误
var (
	// item 为空 错误
	ErrNilItem = errors.New(ErrNilItemText)
	// id 为空 错误
	ErrNilId = errors.New(ErrNilIdText)
	// 已经到达最大管理数 错误
	ErrMaxConnect = errors.New(ErrMaxConnectText)
	// item 的 id 已存在 错误
	ErrItemExist = errors.New(ErrItemExistText)
	// item 的 id 不存在 错误
	ErrItemUnExist = errors.New(ErrItemUnExistText)
)

// 新建端管理者
func NewManage(max int64) Manager {
	// 不填入最大持有链接，则设为默认最大持有链接数 5000
	if max <= 0 {
		max = DefaultMaxConnection
	}
	return &manager{
		max:       max,
		mutex:     sync.Mutex{},
		items:     make(map[string]Item),
		observers: make([]Observer, 0),
	}
}

// 添加一个端到当前管理
func (m *manager) ManageItem(i Item) error {
	if i == nil {
		return ErrNilItem
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 已存在
	if _, exist := m.items[i.ID().String()]; exist {
		return ErrItemExist
	}

	// 最大管理数校验
	if int64(len(m.items)) > m.max {
		return ErrMaxConnect
	}

	// 持有
	m.items[i.ID().String()] = i

	for _, observer := range m.observers {
		go observer.ObserveConnect(i.ID().String())
	}

	//log.Printf("Service now having %d connect", len(m.items))

	return nil
}

// 从当前管理移除指定端
func (m *manager) RemoveItem(id string) error {
	if id == "" {
		return ErrNilId
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 不存在
	if _, exist := m.items[id]; !exist {
		return ErrItemUnExist
	}

	// 删除
	delete(m.items, id)

	for _, observer := range m.observers {
		go observer.ObserveDisconnect(id)
	}

	//log.Printf("Service now having %d connect", len(m.items))

	return nil
}

// 从当前管理寻找指定 id 的端
func (m *manager) FindItem(id string) (Item, error) {
	if id == "" {
		return nil, ErrNilId
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	i, exist := m.items[id]
	// 不存在
	if !exist {
		return nil, ErrItemUnExist
	}

	return i, nil
}

// 目前管理端数目
func (m *manager) CurrentManageCount() int64 {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return int64(len(m.items))
}

// 注册观察者
func (m *manager) RegisterObserver(observer Observer) {
	m.observers = append(m.observers, observer)
}
