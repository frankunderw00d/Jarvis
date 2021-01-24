package main
//
//import (
//	"context"
//	redisGo "github.com/gomodule/redigo/redis"
//	"go.mongodb.org/mongo-driver/bson"
//	"jarvis/base/database"
//	"log"
//	"time"
//)
//
//func init() {
//	// 初始化 MySQL
//	if err := database.InitializeMySQL(
//		"frank",
//		"frank123",
//		"localhost",
//		7000,
//		"jarvis",
//	); err != nil {
//		log.Panicf("Initialize MySQL error : %s", err.Error())
//		return
//	}
//
//	// 设置 MySQL
//	database.SetUpMySQL(time.Minute*time.Duration(5), 10, 5000)
//
//	// 初始化 Redis
//	database.InitializeRedis(time.Minute*time.Duration(5), 10, 5000, "localhost", 8000, "frank123")
//
//	// 初始化 Mongo
//	if err := database.InitializeMongo(
//		"frank",
//		"frank123",
//		"jarvis",
//		"localhost",
//		9000,
//		time.Minute*time.Duration(5),
//		5000); err != nil {
//		log.Panicf("Initialize Mongo error : %s", err.Error())
//		return
//	}
//}
//
//func main() {
//	//mysqlTest()
//	//redisTest()
//	mongoTest()
//
//}
//
//func mysqlTest() {
//	conn, err := database.GetMySQLConn()
//	if err != nil {
//		log.Printf("database.GetMySQLConn error : %s", err.Error())
//		return
//	}
//	defer func() {
//		if err := conn.Close(); err != nil {
//			log.Printf("conn error : %s", err.Error())
//			return
//		}
//	}()
//
//	result, err := conn.ExecContext(context.Background(), "insert into `jarvis`.`static_platform`(name, link, owner) values (?,?,?)",
//		"test",
//		"test-link",
//		"test-owner",
//	)
//	if err != nil {
//		log.Printf("conn exec error : %s", err.Error())
//		return
//	}
//
//	log.Println(result.LastInsertId())
//	log.Println(result.RowsAffected())
//}
//
//func redisTest() {
//	conn, err := database.GetRedisConn()
//	if err != nil {
//		log.Printf("database.GetRedisConn error : %s", err.Error())
//		return
//	}
//	defer func() {
//		if err := conn.Close(); err != nil {
//			log.Printf("conn error : %s", err.Error())
//			return
//		}
//	}()
//
//	a, err := redisGo.String(conn.Do("set", "name", "frank"))
//	if err != nil {
//		log.Printf("conn do error : %s", err.Error())
//		return
//	}
//	log.Println(a)
//}
//
//func mongoTest() {
//	conn, err := database.GetMongoConn("log")
//	if err != nil {
//		log.Printf("database.GetRedisConn error : %s", err.Error())
//		return
//	}
//
//	result, err := conn.InsertOne(context.Background(), bson.M{"title": "panic", "author": "frankunderw00d", "print": "test2"})
//	if err != nil {
//		log.Printf("conn insert error : %s", err.Error())
//		return
//	}
//
//	log.Println(result.InsertedID)
//}
