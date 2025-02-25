package performance

import (
	"encoding/json"
	"fmt"
	"log"
	"space-api/util"
	"space-api/util/ptr"
	"strconv"
	"time"

	"github.com/tidwall/buntdb"
)

type (
	CacheGroupInf interface {
		Group(namespace string) CacheGroupInf
		Set(key string, val any, ttl ...time.Duration) error
		Get(key string, receivePtr any) error
		Delete(key string) error
		GetTTL(key string) (time.Duration, error)
		SetTTL(key string, ttl time.Duration) error
		GetAndDel(key string, receivePtr any) error
		GetAndIncr(key string, inc int64, ttl ...time.Duration) (int64, error)
		IncrAndGet(key string, inc int64, ttl ...time.Duration) (int64, error)
		GetInt64(key string) (int64, error)
		ClearAll()
		Close()
	}

	BuntDBCacheAdapter struct {
		instance  *buntdb.DB
		namespace string
	}
)

var _ CacheGroupInf = (*BuntDBCacheAdapter)(nil)

var DefaultJsonCache CacheGroupInf

func init() {
	DefaultJsonCache = NewBuntDBCache(":memory:")
}

func NewBuntDBCache(cachePath string) *BuntDBCacheAdapter {
	instance, err := buntdb.Open(cachePath)

	if err != nil {
		log.Fatal("initialize cache failure: ", err)
	}
	return &BuntDBCacheAdapter{
		instance: instance,
	}
}

func (c *BuntDBCacheAdapter) getKey(key string) string {
	return util.TernaryExpr(c.namespace != "", c.namespace+":"+key, key)
}

func (c *BuntDBCacheAdapter) Group(namespace string) CacheGroupInf {
	return &BuntDBCacheAdapter{
		instance: c.instance,
		namespace: util.TernaryExpr(
			c.namespace != "",
			c.namespace+":"+namespace,
			namespace,
		),
	}
}

func (c *BuntDBCacheAdapter) Set(key string, val any, ttl ...time.Duration) (err error) {
	key = c.getKey(key)

	ttlLen := len(ttl)
	var newTTL time.Duration
	switch ttlLen {
	case 0:
	case 1:
		newTTL = ttl[0]
	default:
		return fmt.Errorf("want got 0 or zero ttl parameter, but got: %d", ttlLen)
	}

	bf, marshalErr := json.Marshal(val)
	if marshalErr != nil {
		return marshalErr
	}

	return c.instance.Update(func(tx *buntdb.Tx) error {
		expLeave, ttlErr := tx.TTL(key)
		var op *buntdb.SetOptions
		switch {
		case ttlLen == 1:
			op = &buntdb.SetOptions{
				Expires: true,
				TTL:     newTTL,
			}
		case ttlErr == nil && expLeave > 0:
			op = &buntdb.SetOptions{
				Expires: true,
				TTL:     expLeave,
			}
		}
		_, _, setErr := tx.Set(key, ptr.Bytes2String(bf), op)

		return setErr
	})
}

func (c *BuntDBCacheAdapter) Get(key string, receivePtr any) error {
	key = c.getKey(key)

	return c.instance.View(func(tx *buntdb.Tx) error {
		val, getErr := tx.Get(key)
		if getErr != nil {
			return getErr
		}
		return json.Unmarshal(ptr.String2Bytes(val), receivePtr)
	})
}

func (c *BuntDBCacheAdapter) Delete(key string) error {
	key = c.getKey(key)

	return c.instance.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(key)

		return err
	})

}

func (c *BuntDBCacheAdapter) GetTTL(key string) (time.Duration, error) {
	key = c.getKey(key)

	var expLeave time.Duration
	opErr := c.instance.View(func(tx *buntdb.Tx) error {
		d, err := tx.TTL(key)
		expLeave = d
		return err
	})

	return expLeave, opErr
}

func (c *BuntDBCacheAdapter) SetTTL(key string, ttl time.Duration) error {
	key = c.getKey(key)

	return c.instance.Update(func(tx *buntdb.Tx) error {
		val, err := tx.Get(key)
		if err != nil {
			return err
		}

		_, _, setErr := tx.Set(key, val, &buntdb.SetOptions{
			Expires: true,
			TTL:     ttl,
		})
		return setErr
	})
}

func (c *BuntDBCacheAdapter) GetAndDel(key string, receivePtr any) (err error) {
	key = c.getKey(key)

	return c.instance.Update(func(tx *buntdb.Tx) error {
		val, getErr := tx.Delete(key)
		if getErr != nil {
			return getErr
		}

		return json.Unmarshal(ptr.String2Bytes(val), receivePtr)
	})
}

func (c *BuntDBCacheAdapter) GetAndIncr(key string, inc int64, ttl ...time.Duration) (count int64, err error) {
	key = c.getKey(key)

	ttlLen := len(ttl)
	var newTTL time.Duration
	switch ttlLen {
	case 0:
	case 1:
		newTTL = ttl[0]
	default:
		err = fmt.Errorf("want got 0 or zero ttl parameter, but got: %d", ttlLen)
		return
	}

	err = c.instance.Update(func(tx *buntdb.Tx) error {
		ttlLeave, ttlErr := tx.TTL(key)
		existsVal, getErr := tx.Get(key)
		var op *buntdb.SetOptions

		switch {
		case ttlLen == 1:
			op = &buntdb.SetOptions{
				Expires: true,
				TTL:     newTTL,
			}
		case ttlErr == nil && ttlLeave > 0:
			op = &buntdb.SetOptions{
				Expires: true,
				TTL:     ttlLeave,
			}
		}

		if getErr != nil {
			count = 0
			_, _, setErr := tx.Set(key, fmt.Sprintf("%d", inc), op)

			return setErr
		} else {
			c, parsedErr := strconv.ParseInt(existsVal, 10, 64)
			if parsedErr != nil {
				return parsedErr
			}
			count = c
			_, _, setErr := tx.Set(key, fmt.Sprintf("%d", inc+c), op)

			return setErr
		}
	})

	return
}

func (c *BuntDBCacheAdapter) IncrAndGet(key string, inc int64, ttl ...time.Duration) (count int64, err error) {
	key = c.getKey(key)

	ttlLen := len(ttl)
	var newTTL time.Duration
	switch ttlLen {
	case 0:
	case 1:
		newTTL = ttl[0]
	default:
		err = fmt.Errorf("want got 0 or zero ttl parameter, but got: %d", ttlLen)
		return
	}

	err = c.instance.Update(func(tx *buntdb.Tx) error {
		ttlLeave, ttlErr := tx.TTL(key)
		existsVal, getErr := tx.Get(key)
		var op *buntdb.SetOptions

		switch {
		case ttlLen == 1:
			op = &buntdb.SetOptions{
				Expires: true,
				TTL:     newTTL,
			}
		case ttlErr == nil && ttlLeave > 0:
			op = &buntdb.SetOptions{
				Expires: true,
				TTL:     ttlLeave,
			}
		}

		if getErr != nil {
			count = inc
			_, _, setErr := tx.Set(key, fmt.Sprintf("%d", inc), op)

			return setErr
		} else {
			c, parsedErr := strconv.ParseInt(existsVal, 10, 64)
			if parsedErr != nil {
				return parsedErr
			}
			count = c + inc
			_, _, setErr := tx.Set(key, fmt.Sprintf("%d", count), op)

			return setErr
		}
	})

	return
}

func (c *BuntDBCacheAdapter) GetInt64(key string) (count int64, err error) {
	key = c.getKey(key)

	err = c.instance.View(func(tx *buntdb.Tx) error {
		val, getErr := tx.Get(key)
		if getErr != nil {
			return getErr
		}

		return json.Unmarshal(ptr.String2Bytes(val), &count)
	})

	return
}

func (c *BuntDBCacheAdapter) ExposeInstance() *buntdb.DB {
	return c.instance
}

func (c *BuntDBCacheAdapter) ClearAll() {
	c.instance.Update(func(tx *buntdb.Tx) error {
		return tx.DeleteAll()
	})
}

func (c *BuntDBCacheAdapter) Close() {
	c.instance.Close()
}
