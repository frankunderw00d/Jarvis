package database

import (
	"context"
	"errors"
	"fmt"
	mongoGo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type (
	// Mongo 定义
	Mongo interface {
		// 初始化
		Initialize(host string, port int, maxConnIdleTime time.Duration, poolSize uint64) error

		// 获取指向指定 database 的 collection 的连接
		Get(database, collection string) (*mongoGo.Collection, error)

		// 关闭连接池
		Close() error
	}

	// Mongo 定义实现
	mongo struct {
		ctx        context.Context
		cancelFunc context.CancelFunc
		client     *mongoGo.Client
	}
)

const (
	// 默认 Mongo host
	DefaultMongoHost = "localhost"
	// 默认 Mongo port
	DefaultMongoPort = 27017
	// 默认 Mongo 最大连接空闲时间
	DefaultMongoMaxConnIdleTime = time.Duration(120) * time.Second
	// 默认 Mongo 池大小
	DefaultMongoPoolSize = 30
	// 默认 Mongo 驱动名
	DefaultMongoScheme = "mongodb"
)

const (
	ErrNilClientText           = "a_client is nil"
	ErrEmptyDatabaseNameText   = "database name is empty"
	ErrEmptyCollectionNameText = "collection name is empty"
)

var (
	// a_client 为 nil 错误
	ErrNilClient = errors.New(ErrNilClientText)
	// database name 为空 错误
	ErrEmptyDatabaseName = errors.New(ErrEmptyDatabaseNameText)
	// collection name 为空 错误
	ErrEmptyCollectionName = errors.New(ErrEmptyCollectionNameText)
)

// 新建 Mongo
func NewMongo() Mongo {
	return &mongo{}
}

// 初始化
func (m *mongo) Initialize(host string, port int, maxConnIdleTime time.Duration, poolSize uint64) error {
	if host == "" {
		host = DefaultMongoHost
	}
	if port == 0 {
		port = DefaultMongoPort
	}
	if maxConnIdleTime == 0 {
		maxConnIdleTime = DefaultMongoMaxConnIdleTime
	}
	if poolSize == 0 {
		poolSize = DefaultMongoPoolSize
	}

	ctx, cancel := context.WithCancel(context.Background())

	uri := fmt.Sprintf("%s://%s:%d", DefaultMongoScheme, host, port)

	clientOption := options.Client().ApplyURI(uri)
	clientOption.SetMaxConnIdleTime(maxConnIdleTime)
	clientOption.SetMaxPoolSize(poolSize)

	client, err := mongoGo.Connect(ctx, clientOption)
	if err != nil {
		return err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}

	m.client = client
	m.ctx = ctx
	m.cancelFunc = cancel

	return nil
}

// 获取指向指定 database 的 collection 的连接
func (m *mongo) Get(database, collection string) (*mongoGo.Collection, error) {
	if m.client == nil {
		return nil, ErrNilClient
	}
	if database == "" {
		return nil, ErrEmptyDatabaseName
	}
	if collection == "" {
		return nil, ErrEmptyCollectionName
	}

	return m.client.Database(database).Collection(collection), nil
}

// 关闭连接池
func (m *mongo) Close() error {
	if m.cancelFunc != nil {
		m.cancelFunc()
	}

	return m.client.Disconnect(context.TODO())
}
