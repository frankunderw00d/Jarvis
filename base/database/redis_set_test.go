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

func TestRedisSet(t *testing.T) {
	//if v, err := SAdd("UsersInfo", 1, 2, 3, 4, 5); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}

	//if v, err := SIsMember("a", 1); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := SIsMember("bbc", 1); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := SIsMember("UsersInfo", 1); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}

	//if v, err := SPop("a", 10); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := SPop("bbc", 21); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := SPop("UsersInfo", 2); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}

	//if v, err := SRandMember("a", -4); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := SRandMember("bbc", 4); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := SRandMember("UsersInfo", 4); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}

	//if v, err := SRem("a", 3,4,5); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := SRem("bbc", 3,4,5); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := SRem("UsersInfo", 3,4,5); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}

	//if v, err := SMove("a", "b",8); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := SMove("bbc", "b",8); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := SMove("UsersInfo", "b",8); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := SMove("a", "UsersInfo",8); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}

	//if v, err := SCard("a"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := SCard("bbc"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := SCard("UsersInfo"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}

	//if v, err := SMembers("a"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := SMembers("bbc"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := SMembers("UsersInfo"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}

	//if v, err := SInter("a","b","c"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := SInter("a","b","bbc"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := SInter("a","b","UsersInfo"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}

	//if v, err := SDiff("a", "b", "c"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := SDiff("a", "b", "bbc"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := SDiff("a", "b", "UsersInfo"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}

	if v, err := SDiffStore("abc","a", "b", "c"); err != nil {
		log.Println(err.Error())
	} else {
		log.Println(v)
	}

	if v, err := SDiffStore("abbbc","a", "b", "bbc"); err != nil {
		log.Println(err.Error())
	} else {
		log.Println(v)
	}

	if v, err := SDiffStore("abUsersInfo","a", "b", "UsersInfo"); err != nil {
		log.Println(err.Error())
	} else {
		log.Println(v)
	}
}
