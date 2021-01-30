package database

import (
	"testing"
	"time"
)

func init() {
	// 初始化 Redis
	InitializeRedis(time.Minute*time.Duration(5), 10, 5000, "localhost", 8000, "frank123")
}

func TestRedisString(t *testing.T) {
	//if err := Set("name", "frank", RedisStringSetTimeoutEX, 1, RedisStringSetNotExist); err != nil {
	//	log.Println(err.Error())
	//}
	//
	//time.Sleep(time.Second*time.Duration(2))
	//
	//if v, err := Get("name"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}

	//if err := MSet(map[string]interface{}{
	//	"name": "frank",
	//	"age":  1,
	//}); err != nil {
	//	log.Println(err.Error())
	//}

	//if v, err := MGet([]string{
	//	"name",
	//	"sex",
	//}); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Printf("%+v", v)
	//}

	//if v, err := Incr("sex"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Printf("%+v", v)
	//}

	//if v, err := Append("sex","fff"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Printf("%+v", v)
	//}

	//if v, err := StrLen("biu"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Printf("%+v", v)
	//}

	//if v, err := GetSet("biu","cooper"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Printf("%s", v)
	//}

	//if v, err := GetRange("aksfsn", 0, 3); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Printf("%s", v)
	//}
}
