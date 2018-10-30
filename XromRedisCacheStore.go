package jgin

import (
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/gomodule/redigo/redis"
	"github.com/vmihailenco/msgpack"
	"reflect"
)



//var _ core.CacheStore = XromRedisCacheStore()

// XromRedisStore represents in-memory store
type XromRedisStore struct {
	server redis.Conn
	timeOutInSeconds int
}

var typeRegistry = make(map[string]interface{})

// NewXromRedisStore creates a new store in memory
func XromRedisCacheStore(conn redis.Conn,exptime int) *XromRedisStore {
	if exptime<0 {
		exptime = 86400*20
	}
	ins := &XromRedisStore{server: conn,timeOutInSeconds:exptime}
	var _ core.CacheStore = ins
	return ins
}

// Put puts object into store
func (s *XromRedisStore) Put(key string, value interface{}) error {

	storeSet := make(map[string]interface{})
	tname := reflect.TypeOf(value).String()

	switch value.(type) {
		case bool, string, rune, byte, int, int8, int16, int64, float32, float64, complex64, complex128:
			storeSet["t"] = "_origin"
			//storeSet["tp"] = tname
			storeSet["d"] = value
		default:
			data,_ := msgpack.Marshal(value)
			storeSet["d"] = string(data)
			storeSet["t"] = tname
			if _,ok := typeRegistry[tname];!ok{
				typeRegistry[tname] = value //reflect.TypeOf(value).Elem()?
			}
	}

	store, _ := msgpack.Marshal(storeSet)
	if s.timeOutInSeconds>0{
		s.server.Do("SET", key, string(store), "EX", s.timeOutInSeconds)
	} else {
		s.server.Do("SET", key, string(store))
	}
	return nil
}

// Get gets object from store
func (s *XromRedisStore) Get(key string) (interface{}, error) {
	v, err := redis.String(s.server.Do("GET", key))

	if err == nil {
		b := []byte(v)
		storeSet := make(map[string]interface{})
		err = msgpack.Unmarshal(b, &storeSet)

		if packedData,ok := storeSet["d"];ok{
			if t,oked := storeSet["t"];oked{
				tname := t.(string)
				if tname == "_origin" {
						return packedData,nil
						//out, _ := base64.StdEncoding.DecodeString(outstring)
						//return string(out),nil
					return nil, xorm.ErrNotExist
				} else {
					if newStruct, geted := typeRegistry[tname]; geted{
						err := msgpack.Unmarshal([]byte(packedData.(string)),newStruct)
						if err==nil {
							return newStruct,nil
						}
					}
				}
			}
		}

	}
	return nil, xorm.ErrNotExist
}

// Del deletes object
func (s *XromRedisStore) Del(key string) error {
	s.server.Do("DELETE", key)
	return nil
}

