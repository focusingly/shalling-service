package performance

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

func init() {
	maxBfSize := constants.MB * 16
	cacheInstance := freecache.NewCache(int(maxBfSize))
	DefaultJsonCache = (*JsonCache)(&bizCache{
		instance:  cacheInstance,
		namespace: "",
	})
}

// Group 返回新的命名空间存储
func (jc *JsonCache) Group(namespace string) *JsonCache {
	return &JsonCache{
		namespace: namespace + ":",
		instance:  jc.instance,
	}
}

// Set 设置值, 如果 ttl <=0, 那么表示永不过期
func (jc *JsonCache) Set(key string, val any, ttl ...Second) error {
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

func (jc *JsonCache) Get(key string, receivePointer any) error {
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

// ClearAll 清空缓存, 如果是根命名空间, 那么会清空所有的缓存, 如果是具体的命名空间, 那么会通过遍历逐一删除(不是原子操作)
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

// SetTTL 重新设置过期时间, 如果 ttl <=0, 那么表示永远不过期
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

// GetAndDel 获取并删除一个值, 这个过程不是原子操作(受限于非重入锁)
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

// GetAndIncr 先获取再更新一个整数记录, 如果不存在, 那么进行创建
func (jc *JsonCache) GetAndIncr(key string, inc int, ttl Second) (count int, err error) {
	key = jc.namespace + key
	_, _, err = jc.instance.Update(
		ptr.String2Bytes(key),
		func(value []byte, found bool) (newValue []byte, replace bool, expireSeconds int) {
			var intVal int
			if err := json.Unmarshal(value, &intVal); err != nil {
				panic(err)
			}
			count = intVal
			intVal++

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

// IncAndGet 先更新再获取一个整数记录, 如果不存在, 那么进行创建
func (jc *JsonCache) IncAndGet(key string, inc int, ttl Second) (count int, err error) {
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

// ExposeInstance 直接暴露内部所使用的 freeCache 实例
func (jc *JsonCache) ExposeInstance() *freecache.Cache {
	return jc.instance
}
