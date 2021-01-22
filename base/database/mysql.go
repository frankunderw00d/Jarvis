package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type (
	MySQL interface {
		// 初始化
		Initialize(username, password, host string, port int, database string) error

		// 设置
		SetUp(connMaxLife time.Duration, maxIdle, maxOpen int)

		// 获取连接
		Get() (*sql.Conn, error)

		// 关闭连接池
		Close() error
	}

	mysql struct {
		db *sql.DB
	}
)

const (
	// MySQL 默认 host
	DefaultMySQLHost = "127.0.0.1"
	// MySQL 默认 port
	DefaultMySQLPort = 3306
	// MySQL 驱动名
	MySQLDriverName = "mysql"
)

const (
	ErrNilDBText = "db is nil"
)

var (
	// db 为 nil 错误
	ErrNilDB = errors.New(ErrNilDBText)
)

// 新建 MySQL
func NewMySQL() MySQL {
	return &mysql{}
}

// 初始化
func (m *mysql) Initialize(username, password, host string, port int, database string) error {
	if host == "" {
		host = DefaultMySQLHost
	}
	if port == 0 {
		port = DefaultMySQLPort
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true",
		username,
		password,
		host,
		port,
		database,
	)

	db, err := sql.Open(MySQLDriverName, dsn)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	m.db = db

	return nil
}

// 设置
func (m *mysql) SetUp(connMaxLife time.Duration, maxIdle, maxOpen int) {
	if m.db == nil {
		return
	}

	m.db.SetConnMaxLifetime(connMaxLife)
	m.db.SetMaxIdleConns(maxIdle)
	m.db.SetMaxOpenConns(maxOpen)
}

// 获取连接
func (m *mysql) Get() (*sql.Conn, error) {
	if m.db == nil {
		return nil, ErrNilDB
	}

	return m.db.Conn(context.Background())
}

// 关闭连接
func (m *mysql) Close() error {
	return m.db.Close()
}
