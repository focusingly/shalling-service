package performance

import (
	"encoding/json"
	"fmt"
	"space-api/constants"
	"space-api/util/ptr"
	"strconv"
	"strings"

	"github.com/coocood/freecache"
)

type bizCache struct {
	namespace string
	instance  *freecache.Cache
}

type JsonCache bizCache
type Second = int

var DefaultJsonCache *JsonCache = NewCache(constants.MB * 16)

func NewCache(maxSize constants.MemoryByteSize) *JsonCache {
	return (*JsonCache)(&bizCache{
		instance:  freecache.NewCache(int(maxSize)),
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
	var expireLeaves int
	if expire, e := jc.instance.TTL([]byte(key)); e == nil {
		expireLeaves = int(expire)
	}
	switch len(ttl) {
	case 0:
	case 1:
		expireLeaves = ttl[0]
	default:
		panic(fmt.Errorf("ttl expect 0 or zero param, but got: %d", len(ttl)))
	}

	bf, err := json.Marshal(val)
	if err != nil {
		return err
	}

	return jc.instance.Set(ptr.String2Bytes(key), bf, expireLeaves)
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

// ClearAll 清空根命名空间缓存
func (jc *JsonCache) ClearAll() {
	jc.instance.Clear()
}

// ClearGroupCache 清空所属命名空间下的所有缓存, 会通过遍历逐一删除(不是原子操作)
func (jc *JsonCache) ClearGroupCache() {
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
func (jc *JsonCache) GetAndIncr(key string, inc int64, ttl ...Second) (count int64, err error) {
	key = jc.namespace + key
	var expireLeaves int

	if expire, e := jc.instance.TTL([]byte(key)); e == nil {
		expireLeaves = int(expire)
	}

	switch len(ttl) {
	case 0:
	case 1:
		expireLeaves = ttl[0]
	default:
		panic(fmt.Errorf("ttl expect 0 or zero param, but got: %d", len(ttl)))
	}

	_, _, err = jc.instance.Update(
		ptr.String2Bytes(key),
		func(value []byte, found bool) (newValue []byte, replace bool, expireSeconds int) {
			replace = true
			expireSeconds = expireLeaves

			// 未找到值
			if !found {
				count = 0
				newValue = ptr.String2Bytes(fmt.Sprintf("%d", inc))
				return
			} else {
				parsedInt, err := strconv.ParseInt(ptr.Bytes2String(value), 10, 64)
				if err != nil {
					panic(fmt.Errorf("already exists key in using: %s", err.Error()))
				}
				count = parsedInt
				newValue = ptr.String2Bytes(fmt.Sprintf("%d", inc+parsedInt))
			}

			return
		})

	return
}

// IncAndGet 先更新再获取一个整数记录, 如果不存在, 那么进行创建
func (jc *JsonCache) IncAndGet(key string, inc int64, ttl ...Second) (count int64, err error) {
	key = jc.namespace + key
	var expireLeaves int

	if expire, e := jc.instance.TTL([]byte(key)); e == nil {
		expireLeaves = int(expire)
	}

	switch len(ttl) {
	case 0:
	case 1:
		expireLeaves = ttl[0]
	default:
		panic(fmt.Errorf("ttl expect 0 or zero param, but got: %d", len(ttl)))
	}

	_, _, err = jc.instance.Update(
		ptr.String2Bytes(key),
		func(value []byte, found bool) (newValue []byte, replace bool, expireSeconds int) {
			// 进行替换
			replace = true
			expireSeconds = expireLeaves
			// 未找到值
			if !found {
				count = inc
				newValue = ptr.String2Bytes(fmt.Sprintf("%d", inc))
				return
			} else {
				parsedInt, err := strconv.ParseInt(ptr.Bytes2String(value), 10, 64)
				if err != nil {
					panic(fmt.Errorf("already exists key in using: %s", err.Error()))
				}
				count = parsedInt + inc
				newValue = ptr.String2Bytes(fmt.Sprintf("%d", count))
			}

			return
		})

	return
}

func (jc *JsonCache) GetInt64(key string) (count int64, err error) {
	key = jc.namespace + key
	val, err := jc.instance.Get(ptr.String2Bytes(key))
	if err != nil {
		return
	}
	count, err = strconv.ParseInt(ptr.Bytes2String(val), 10, 64)

	return
}

// ExposeInstance 直接暴露内部所使用的 freeCache 实例
func (jc *JsonCache) ExposeInstance() *freecache.Cache {
	return jc.instance
}
