package database

import (
	"log"
	"testing"
	"time"
)

func init() {
	// 初始化 Redis
	InitializeRedis(time.Minute*time.Duration(5), 10, 5000, "localhost", 8000, "frank123")
}

func TestRedisHash(t *testing.T) {
	//if v, err := HSet("a", map[string]interface{}{
	//	"sex":   true,
	//	"music": "igo",
	//}); err != nil {
	//	log.Println(err)
	//} else {
	//	log.Println(v)
	//}

	//if v, err := HSetNX("a","name","frank"); err != nil {
	//	log.Println(err)
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := HSetNX("a","name","bbbb"); err != nil {
	//	log.Println(err)
	//} else {
	//	log.Println(v)
	//}

	//if v, err := HGet("UsersInfo", "name"); err != nil {
	//	log.Println(err)
	//} else {
	//	log.Println(v)
	//}

	//if v, err := HExists("a", "name"); err != nil {
	//	log.Println(err)
	//} else {
	//	log.Println(v)
	//}

	//if v, err := HDel("a", "name", "age", "sex"); err != nil {
	//	log.Println(err)
	//} else {
	//	log.Println(v)
	//}

	//if v, err := HLen("ab"); err != nil {
	//	log.Println(err)
	//} else {
	//	log.Println(v)
	//}

	//if v, err := HStrLen("a", "music"); err != nil {
	//	log.Println(err)
	//} else {
	//	log.Println(v)
	//}

	//if v, err := HIncrBy("a", "age", 5); err != nil {
	//	log.Println(err)
	//} else {
	//	log.Println(v)
	//}

	//if v, err := HIncrByFloat("ab", "age", 5.56); err != nil {
	//	log.Println(err)
	//} else {
	//	log.Println(v)
	//}

	//if  err := HMSet("a", map[string]interface{}{
	//	"sex":   true,
	//	"music": "igo",
	//}); err != nil {
	//	log.Println(err)
	//}

	//if v, err := HMGet("a", "naem", "sex", "age"); err != nil {
	//	log.Println(err)
	//} else {
	//	log.Printf("%+v", v)
	//}
	//
	//if v, err := HMGet("bbc", "naem", "sex", "age"); err != nil {
	//	log.Println(err)
	//} else {
	//	log.Printf("%+v", v)
	//}
	//
	//if v, err := HMGet("b", "naem", "sex", "age"); err != nil {
	//	log.Println(err)
	//} else {
	//	log.Printf("%+v", v)
	//}

	//if v, err := HKeys("a"); err != nil {
	//	log.Println(err)
	//} else {
	//	log.Printf("%+v", v)
	//}
	//
	//if v, err := HKeys("bbc"); err != nil {
	//	log.Println(err)
	//} else {
	//	log.Printf("%+v", v)
	//}
	//
	//if v, err := HKeys("b"); err != nil {
	//	log.Println(err)
	//} else {
	//	log.Printf("%+v", v)
	//}

	//if v, err := HVals("a"); err != nil {
	//	log.Println(err)
	//} else {
	//	log.Printf("%+v", v)
	//}
	//
	//if v, err := HVals("bbc"); err != nil {
	//	log.Println(err)
	//} else {
	//	log.Printf("%+v", v)
	//}
	//
	//if v, err := HVals("b"); err != nil {
	//	log.Println(err)
	//} else {
	//	log.Printf("%+v", v)
	//}

	if v, err := HGetAll("a"); err != nil {
		log.Println(err)
	} else {
		log.Printf("%+v", v)
	}

	if v, err := HGetAll("bbc"); err != nil {
		log.Println(err)
	} else {
		log.Printf("%+v", v)
	}

	if v, err := HGetAll("b"); err != nil {
		log.Println(err)
	} else {
		log.Printf("%+v", v)
	}
}
