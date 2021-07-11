package discovery

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/dollarkillerx/light/transport"
	"github.com/dollarkillerx/light/utils"
	"github.com/gomodule/redigo/redis"
)

type RedisDiscovery struct {
	hearBeat uint // sec
	pool     *redis.Pool
	close    chan struct{}
}

func NewRedisDiscovery(addr string, hearBeat uint, auth *string) (*RedisDiscovery, error) {
	dis := &RedisDiscovery{
		hearBeat: 3,
		close:    make(chan struct{}, 0),
	}
	if hearBeat >= 3 {
		dis.hearBeat = hearBeat
	}
	pool := &redis.Pool{
		MaxIdle:     10,                // 最大空闲连接数
		MaxActive:   10,                // 最大连接数
		IdleTimeout: 300 * time.Second, // 超时回收
		Dial: func() (conn redis.Conn, e error) {
			// 1. 打开连接
			dial, e := redis.Dial("tcp", addr)
			if e != nil {
				log.Fatalln("Redis Pool Err: ", e)
				return nil, e
			}
			// 2. 访问认证
			if auth != nil {
				dial.Do("AUTH", *auth)
			}
			return dial, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error { // 定时检查连接是否可用
			// time.Since(t) 获取离现在过了多少时间
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			log.Fatalln("Redis Pool Err: ", err)
			return err
		},
	}

	dis.pool = pool
	return dis, nil
}

func (r *RedisDiscovery) Discovery(serName string) ([]*Server, error) {
	path := fmt.Sprintf("/registry/%s/*", serName)
	red := r.pool.Get()
	defer red.Close()

	sers := make([]*Server, 0)

	values, err := redis.Strings(red.Do("keys", path))
	if err != nil {
		return nil, err
	}

	for _, v := range values {
		byt, err := redis.Bytes(red.Do("get", v))
		if err != nil {
			return nil, err
		}

		ser := Server{}
		err = json.Unmarshal(byt, &ser)
		if err != nil {
			log.Println(err)
			continue
		}

		sers = append(sers, &ser)
	}

	return sers, nil
}

func (r *RedisDiscovery) Registry(serName, addr string, weights float64, protocol transport.Protocol, serID *string) error {
	if serID == nil {
		id, err := utils.DistributedID()
		if err != nil {
			return err
		}
		serID = &id
	}

	serVer := Server{
		Addr:     addr,
		ID:       *serID,
		Weights:  weights,
		Protocol: protocol,
	}

	serNode, err := json.Marshal(serVer)
	if err != nil {
		return err
	}
	path := r.getRedisPath(serName, *serID)
	err = r.registry(path, serNode)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case <-r.close:
				return
			case <-time.After(time.Second * time.Duration(r.hearBeat)):
				err = r.registry(path, serNode)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}()

	return nil
}

func (r *RedisDiscovery) registry(path string, byte []byte) error {
	rds := r.pool.Get()
	defer rds.Close()

	_, err := rds.Do("setex", path, r.hearBeat, byte)
	return err
}

func (r *RedisDiscovery) UnRegistry(serName string, serID string) error {
	close(r.close)
	return nil
}

// 获取redis 存储 路径 [redis]中存储的格式 /registry/服务名称/服务id
func (r *RedisDiscovery) getRedisPath(serName string, id string) string {
	return fmt.Sprintf("/registry/%s/%s", serName, id)
}
