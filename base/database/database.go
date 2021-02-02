package database

import (
	"database/sql"
	mongoGo "go.mongodb.org/mongo-driver/mongo"
	"time"
)

type ()

const ()

var (
	// 默认的 MySQL 实例
	defaultMySQL = NewMySQL()

	//// 默认的 Redis 实例
	//defaultRedis = redis.NewRedis()

	// 默认的 Mongo 实例
	defaultMongo = NewMongo()
)

// --------------------------------------------------- MySQL Public Methods --------------------------------------------
// 1.初始化 MySQL
func InitializeMySQL(username, password, host string, port int, database string) error {
	return defaultMySQL.Initialize(username, password, host, port, database)
}

// 2.设置 MySQL
func SetUpMySQL(connMaxLife time.Duration, maxIdle, maxOpen int) {
	defaultMySQL.SetUp(connMaxLife, maxIdle, maxOpen)
}

// 3.MySQL 获取连接
func GetMySQLConn() (*sql.Conn, error) {
	return defaultMySQL.Get()
}

// 4.关闭 MySQL
func CloseMySQL() error {
	return defaultMySQL.Close()
}

//// --------------------------------------------------- Redis Public Methods --------------------------------------------
//// 1.初始化 Redis
//func InitializeRedis(idleTimeout time.Duration, maxIdle, maxActive int, host string, port int, password string) {
//	defaultRedis.Initialize(idleTimeout, maxIdle, maxActive, host, port, password)
//}
//
//// 2.Redis 获取连接
//func GetRedisConn() (redisGo.Conn, error) {
//	return defaultRedis.Get()
//}
//
//// 3.关闭 Redis
//func CloseRedis() error {
//	return defaultRedis.Close()
//}

// --------------------------------------------------- Mongo Public Methods --------------------------------------------
// 1.初始化 Mongo
func InitializeMongo(username, password, database, host string, port int, maxConnIdleTime time.Duration, poolSize uint64) error {
	return defaultMongo.Initialize(username, password, database, host, port, maxConnIdleTime, poolSize)
}

// 2.Mongo 获取连接
func GetMongoConn(collection string) (*mongoGo.Collection, error) {
	return defaultMongo.Get(collection)
}

// 3.关闭 Mongo
func CloseMongo() error {
	return defaultMongo.Close()
}

/*
MySQL 使用:

	if err := InitializeMySQL(
		"frank",
		"frank123",
		"localhost",
		7000,
		"jarvis",
	); err != nil {
		log.Printf("Initialize MySQL error : %s", err)
		return
	}

	SetUpMySQL(time.Minute*time.Duration(5), 10, 30)

	conn, err := GetMySQLConn()
	if err != nil {
		log.Printf("Get MySQL conn error : %s", err)
		return
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("db Close : %s", err)
			return
		}
	}()

	if err := conn.PingContext(context.Background()); err != nil {
		log.Printf("db ping error : %s", err)
		return
	}

	rows, err := conn.QueryContext(context.Background(), "select * from user;")
	if err != nil {
		log.Printf("conn query error : %s", err)
		return
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("rows Close : %s", err)
			return
		}
	}()

	type User struct {
		ID       int
		Platform int
		UserID   int
		UserName string
	}

	users := make([]User, 0)
	for rows.Next() {
		user := User{}
		if err := rows.Scan(&user.ID, &user.Platform, &user.UserID, &user.UserName); err != nil {
			log.Printf("rows scan : %s", err)
			continue
		}
		users = append(users, user)
	}

	if len(users) > 0 {
		for _, user := range users {
			log.Printf("%+v", user)
		}
	}

*/

/*
Redis 使用:

	InitializeRedis(0, 0, 0, "", 0, "frank123")

	conn, err := GetRedisConn()
	if err != nil {
		log.Printf("Redis GetRedisConn error : %s", err.Error())
		return
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("redis.Conn.Close error : %s", err.Error())
			return
		}
	}()

	r, err := conn.Do("get", "name")
	if err != nil {
		log.Printf("redis.Conn.Do error : %s", err.Error())
		return
	}

	data, ok := r.([]byte)
	if !ok {
		log.Printf("redis.Conn.Do.Reply : %v", r)
		return
	}

	log.Printf("Reply : %s", string(data))
*/

/*
Mongo 使用:

	if err := InitializeMongo("", 9000, time.Minute*time.Duration(5), 30); err != nil {
		log.Printf("Initialize mongo error : %s", err.Error())
		return
	}
	defer func() {
		if err := CloseMongo(); err != nil {
			log.Printf("Close mongo error : %s", err.Error())
			return
		}
	}()

	collection, err := GetMongoConn("log", "log_2020_11_24")
	if err != nil {
		log.Printf("get mongo connection error : %s", err.Error())
		return
	}

	singleResult := collection.FindOne(context.Background(), bson.D{{"Level", "debug"}})
	if singleResult.Err() != nil {
		log.Printf("mongoGo.FindOne error : %s", singleResult.Err().Error())
		return
	}
	raw, err := singleResult.DecodeBytes()
	if err != nil {
		log.Printf("mongoGo.SingleResult decodeBytes error : %s", err.Error())
		return
	}

	log.Printf("%s", raw.String())
*/
