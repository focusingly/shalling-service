package util

import (
	"encoding/json"
	"fmt"
	"space-api/constants"
	"space-api/util/ptr"
	"strings"

	"github.com/coocood/freecache"
)

type bizCache struct {
	namespace string
	instance  *freecache.Cache
}

type JsonCache bizCache
type Second = int

var DefaultJsonCache *JsonCache

// Group 返回新的命名空间存储
func (jc *JsonCache) Group(namespace string) *JsonCache {
	return &JsonCache{
		namespace: namespace + ":",
		instance:  jc.instance,
	}
}

func (jc *JsonCache) SetWith(key string, val any, ttl ...Second) error {
	key = jc.namespace + key

	le := len(ttl)
	bf, err := json.Marshal(val)
	if err != nil {
		return err
	}
	switch {
	case le == 0:
		return jc.instance.Set(ptr.String2Bytes(key), bf, 0)
	case le == 1:
		return jc.instance.Set(ptr.String2Bytes(key), bf, ttl[0])
	default:
		return fmt.Errorf("want 0 or 1 ttl arg, but got: %d", le)
	}
}

func (jc *JsonCache) GetById(key string, receivePointer any) error {
	key = jc.namespace + key
	bf, err := jc.instance.Get(ptr.String2Bytes(key))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bf, receivePointer); err != nil {
		return nil
	}
	return nil
}

func (jc *JsonCache) Delete(key string) bool {
	key = jc.namespace + key
	return jc.instance.Del(ptr.String2Bytes(key))
}

// ClearAll 清空缓存, 如果是根命名空间, 那么会清空所有的缓存, 如果是具体的命名空间, 那么会通过遍历逐一删除
func (jc *JsonCache) ClearAll() {
	if jc.namespace == "" {
		jc.instance.Clear()
	} else {
		it := jc.instance.NewIterator()
		entry := it.Next()
		for {
			if entry == nil {
				break
			}
			if strings.HasPrefix(ptr.Bytes2String(entry.Key), jc.namespace) {
				jc.instance.Del(entry.Key)
			}
		}

	}
}

func (jc *JsonCache) GetTTL(key string) (uint32, error) {
	key = jc.namespace + key
	return jc.instance.TTL(ptr.String2Bytes(key))
}

func (jc *JsonCache) SetTTL(key string, ttl Second) (err error) {
	key = jc.namespace + key
	_, _, e := jc.instance.Update(
		ptr.String2Bytes(key),
		func(value []byte, found bool) (newValue []byte, replace bool, expireSeconds int) {
			if !found {
				err = fmt.Errorf("could not find match key: %s", key)
				return
			}

			newValue = value
			replace = true
			expireSeconds = ttl
			return
		})
	if e != nil {
		return e
	}

	return
}

func (jc *JsonCache) GetAndDel(key string, receiverPtr any) (err error) {
	key = jc.namespace + key
	bf, err := jc.instance.Get(ptr.String2Bytes(key))
	if err != nil {
		return
	}
	if err = json.Unmarshal(bf, receiverPtr); err != nil {
		return
	}
	jc.instance.Del(ptr.String2Bytes(key))

	return
}

// GetBeforeAdd 获取某个整数类型的值, 并更新, 如果不存在, 那么进行创建
func (jc *JsonCache) GetBeforeAdd(key string, inc int, ttl Second) (count int, err error) {
	key = jc.namespace + key
	_, _, err = jc.instance.Update(
		ptr.String2Bytes(key),
		func(value []byte, found bool) (newValue []byte, replace bool, expireSeconds int) {
			var intVal int
			if err := json.Unmarshal(value, &intVal); err != nil {
				panic(err)
			}
			intVal++
			count = intVal

			if bf, err := json.Marshal(intVal); err != nil {
				panic(err)
			} else {
				newValue = bf
			}
			replace = true
			expireSeconds = ttl

			return
		})

	return
}

func (jc *JsonCache) ExposeInstance() *freecache.Cache {
	return jc.instance
}

func init() {
	maxBfSize := constants.MB * 16
	cacheInstance := freecache.NewCache(int(maxBfSize))
	DefaultJsonCache = (*JsonCache)(&bizCache{
		instance:  cacheInstance,
		namespace: "",
	})
}
