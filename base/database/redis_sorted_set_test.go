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

func TestRedisSortedSet(t *testing.T) {
	//if v, err := ZAdd("a", map[string]float64{
	//	"frank": 100.15,
	//	"tom":   74.213,
	//}); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := ZAdd("bbc", map[string]float64{
	//	"frank": 100.15,
	//	"tom":   74.213,
	//}); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := ZAdd("UsersInfo", map[string]float64{
	//	"frank": 100.15,
	//	"tom":   74.213,
	//}); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}

	//if v, err := ZScore("a", "1"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := ZScore("b", "1"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := ZScore("UsersInfo", "1"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}

	//if v, err := ZIncrBy("a", 1.1, "frank"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := ZIncrBy("b", 1.1, "frank"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := ZIncrBy("UsersInfo", 1.1, "frank"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}

	//if v, err := ZCard("a"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := ZCard("b"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := ZCard("UsersInfo"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}

	//if v, err := ZCount("a",0.1,10000); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := ZCount("b",0.1,10000); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := ZCount("UsersInfo",0.1,10000); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}

	//if v, err := ZRange("a",0,-1,RedisSortedSetWithScores); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := ZRange("a",0,-1,RedisSortedSetWithOutScores); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := ZRange("b",0,-1,RedisSortedSetWithScores); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := ZRange("UsersInfo",0,-1,RedisSortedSetWithScores); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}

	//if v, err := ZRevRange("a", 0, -1, RedisSortedSetWithScores); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := ZRevRange("a", 0, -1, RedisSortedSetWithOutScores); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := ZRevRange("b", 0, -1, RedisSortedSetWithScores); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := ZRevRange("UsersInfo", 0, -1, RedisSortedSetWithScores); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}

	//if v, err := ZRangeByScore("a", RedisSortedSetValueNegativeInf, RedisSortedSetValuePositiveInf, RedisSortedSetWithOutScores, RedisSortedSetLimitNone); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := ZRangeByScore("b", RedisSortedSetValueNegativeInf, RedisSortedSetValuePositiveInf, RedisSortedSetWithOutScores, Limit("limit 0 1")); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := ZRangeByScore("UsersInfo", RedisSortedSetValueNegativeInf, RedisSortedSetValuePositiveInf, RedisSortedSetWithOutScores, Limit("limit 0 1")); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}

	//if v, err := ZRevRangeByScore("a", RedisSortedSetValuePositiveInf, RedisSortedSetValueNegativeInf, RedisSortedSetWithScores, Limit("limit 0 5")); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := ZRevRangeByScore("b", RedisSortedSetValuePositiveInf, RedisSortedSetValueNegativeInf, RedisSortedSetWithScores, Limit("limit 0 5")); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := ZRevRangeByScore("UsersInfo", RedisSortedSetValuePositiveInf, RedisSortedSetValueNegativeInf, RedisSortedSetWithScores, Limit("limit 0 5")); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}

	//if v, err := ZRevRank("a", "frank"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := ZRevRank("a", "frank123"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := ZRevRank("b", "frank"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := ZRank("UsersInfo", "frank"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}

	//if v, err := ZRem("a", "frank", "tom", "chaison"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := ZRem("b", "frank", "tom", "chaison"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := ZRem("UsersInfo", "frank", "tom", "chaison"); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}

	//if v, err := ZRemRangeByRank("a", 1,3); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := ZRemRangeByRank("b", 1,3); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := ZRemRangeByRank("UsersInfo", 1,3); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}

	//if v, err := ZRemRangeByScore("a", RedisSortedSetValueNegativeInf, RedisSortedSetValuePositiveInf); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := ZRemRangeByScore("b", Value("1"), Value("3")); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := ZRemRangeByScore("UsersInfo", Value("1"), Value("3")); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}

	//if v, err := ZRangeByLex("a", LEXValue("(a"), LEXValue("[d"), RedisSortedSetLimitNone); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := ZRangeByLex("a", LEXValue("(a"), LEXValue("[d"), Limit("limit 0 2")); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := ZRangeByLex("a", RedisSortedSetValueMin, RedisSortedSetValueMax, Limit("limit 0 3")); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := ZRangeByLex("b", LEXValue("(a"), LEXValue("[d"), RedisSortedSetLimitNone); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := ZRangeByLex("UsersInfo", LEXValue("(a"), LEXValue("[d"), RedisSortedSetLimitNone); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}

	//if v, err := ZRemRangeByLex("a", LEXValue("(a"), LEXValue("[d")); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := ZRemRangeByLex("b", LEXValue("(a"), LEXValue("[d")); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}
	//
	//if v, err := ZRemRangeByLex("UsersInfo", LEXValue("(a"), LEXValue("[d")); err != nil {
	//	log.Println(err.Error())
	//} else {
	//	log.Println(v)
	//}

	if v, err := ZInterStore("ab", []string{"a", "b"}, []float64{2, 1}, "sum"); err != nil {
		log.Println(err.Error())
	} else {
		log.Println(v)
	}

	if v, err := ZInterStore("ac", []string{"a", "c"}, []float64{2, 1}, "sum"); err != nil {
		log.Println(err.Error())
	} else {
		log.Println(v)
	}

	if v, err := ZInterStore("au", []string{"a", "UsersInfo"}, []float64{2, 1}, "sum"); err != nil {
		log.Println(err.Error())
	} else {
		log.Println(v)
	}
}
