package redis

import (
	"fmt"
	redisGo "github.com/gomodule/redigo/redis"
	"log"
	"testing"
	"time"
)

func init() {
	// 初始化 Redis
	InitializeRedis(time.Minute*time.Duration(5), 10, 5000, "localhost", 8000, "frank123")
}

func TestRedisDistributedSync(t *testing.T) {
	ds := NewDistributedSync("ANNOUNCE")
	if err := ds.Join(); err != nil {
		log.Printf("Join error : %s", err.Error())
		return
	}
	defer ds.Clean()

	go func() {
		i := 0
		for i < 10 {
			v := fmt.Sprintf("This is [%d] announce", i+1)
			if err := ds.Publish(v); err != nil {
				log.Println(err.Error())
				break
			}

			time.Sleep(time.Second * time.Duration(10))
			i++
		}
	}()

	for {
		v, err := ds.Subscribe()
		if err != nil {
			if err == redisGo.ErrNil {
				time.Sleep(time.Second * time.Duration(50))
				continue
			}
			log.Println(err.Error())
			break
		}

		log.Printf("Subscribe : %s", v)
	}
}
