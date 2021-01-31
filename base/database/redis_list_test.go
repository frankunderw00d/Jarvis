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

func TestRedisList(t *testing.T) {
	//if v, err := LPush("bbc", 1, 2, 3, 4, 5); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}

	//if v, err := LPushX("UsersInfo", 7); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}

	//if v, err := RPush("UsersInfo", 0,1,2,3,4,5,"music"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}

	//if v, err := RPushX("UsersInfo", 7); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := RPushX("bbc", 7); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := RPushX("bbcd", 7); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}

	//if v, err := LPop("UsersInfo"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := LPop("bbc"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := LPop("bbcd"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}

	//if v, err := RPop("UsersInfo"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := RPop("bbc"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := RPop("bbcd"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}

	//if v, err := RPopLPush("UsersInfo","b"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := RPopLPush("a","b"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := RPopLPush("a","c"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}

	//if v, err := LRem("bbcd", 0, 22); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}

	//if v, err := LLen("a"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := LLen("ab"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := LLen("UsersInfo"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}

	//if v, err := LIndex("a", 0); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := LIndex("ab", 0); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := LIndex("UsersInfo", 0); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}

	//if v, err := LInsert("a", RedisListInsertBefore, "3","1675"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := LInsert("ab", RedisListInsertBefore, "3","1"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := LInsert("UsersInfo", RedisListInsertBefore, "3","1"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}

	//if err := LSet("a", 1, "1993"); err != nil {
	//	log.Println(err.Error())
	//}
	//
	//if err := LSet("ab", 1, "1"); err != nil {
	//	log.Println(err.Error())
	//}
	//
	//if err := LSet("UsersInfo", 1, "1"); err != nil {
	//	log.Println(err.Error())
	//}

	//if v, err := LRange("a", 0, 3); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := LRange("ab", 0, 3); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := LRange("UsersInfo", 0, 3); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}

	//if err := LTrim("a", 2, 7); err != nil {
	//	log.Println(err.Error())
	//}
	//
	//if err := LTrim("ab", 2, 7); err != nil {
	//	log.Println(err.Error())
	//}
	//
	//if err := LTrim("UsersInfo", 2, 7); err != nil {
	//	log.Println(err.Error())
	//}

	//if v, err := BLPop([]string{"a", "ab"}, 0); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Printf("%+v", v)
	//}

	//if v, err := BRPop([]string{"a", "ab"}, 0); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Printf("%+v", v)
	//}

	if v, err := BRPopLPush("a","b",0); err != nil {
		log.Println(err.Error())
	} else {
		log.Printf("%s", v)
	}
}
