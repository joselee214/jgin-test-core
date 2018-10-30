package jgin

import (
	"github.com/gomodule/redigo/redis"
	"time"
)



//一个map，别名->连接池子
type ConnCachePools struct {
	conns map[string]interface{} //保存实际的对象..
	connsmap map[string]string  //别名=>conns的key的转译
}

var ConnPools *ConnCachePools

func init()  {
	ConnPools = &ConnCachePools{conns: make(map[string]interface{}),connsmap:make(map[string]string)}
}

func (cpool *ConnCachePools)GetRedisA(aliasName string) (redis.Conn) {
	return cpool.GetRedis(aliasName,"","")
}

func (cpool *ConnCachePools)GetRedis(aliasName,hostAndPort,passwd string) (redis.Conn) {

	var connsKey = ""
	if aliasName!=""{
		if v,ok := cpool.connsmap[aliasName]; ok {
			connsKey = v
		}
	}
	if connsKey=="" {
		connsKey = getMd5(hostAndPort+passwd)
	}

	if value, ok := cpool.conns[connsKey]; ok {
		if pool,ok := value.(*redis.Pool); ok {
			return pool.Get()
		}
	}
	return nil
}

func (cpool *ConnCachePools)AddRedis(aliasName,hostAndPort,passwd string,arg ...int) (redis.Conn) {
	connsKey := getMd5(hostAndPort+passwd)
	if aliasName=="" {
		aliasName = connsKey
	}
	if value, ok := cpool.conns[connsKey]; ok {
		if pool,ok := value.(*redis.Pool); ok {
			cpool.connsmap[aliasName] = connsKey
			return  pool.Get()
		}
	}
	pool := newPool(hostAndPort,passwd,arg...)
	cpool.conns[connsKey] = pool
	cpool.connsmap[aliasName] = connsKey
	return  pool.Get()
}


//初始化一个Redis Pool
func newPool(hostAndPort,passwd string,arg ... int) *redis.Pool {
	//db,maxIdle,maxActive,idleTimeoutSecond
	params := make([]int, 4)
	for i, v := range arg {
		params[i] = v
	}
	db := params[0]
	if db<0 {
		db = 0
	}
	maxIdle := params[1]
	if maxIdle<=0 {
		maxIdle = 10
	}
	maxActive := params[2]
	if maxActive<=0 {
		maxActive = 2000
	}
	idleTimeoutSecond := params[3]
	if idleTimeoutSecond<=0 {
		idleTimeoutSecond = 240
	}
	idleTimeout := time.Duration(idleTimeoutSecond) * time.Second
	return &redis.Pool{
		MaxIdle:     maxIdle, //最大空闲连接
		MaxActive:   maxActive, //最大连接数
		IdleTimeout: idleTimeout, //超时时间
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", hostAndPort)
			if err != nil {
				return nil, err
			}
			if passwd!="" {
				if _, err := c.Do("AUTH", passwd); err != nil {
					c.Close()
					return nil, err
				}
			}
			// 选择db
			c.Do("SELECT", db)
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}