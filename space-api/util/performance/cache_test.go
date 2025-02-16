package performance_test

import (
	"fmt"
	"space-api/constants"
	"space-api/util/performance"
	"testing"
	"time"
)

func TestCacheIncrOperation(*testing.T) {
	cache := performance.NewCache(constants.MB * 1).Group("pv")
	const key = "pv:123"

	val, _ := cache.IncAndGet(key, 1, 0)
	fmt.Println(val)
	val2, _ := cache.IncAndGet(key, 3, 0)
	fmt.Println(val2)

}

func TestCacheIncrOperation2(*testing.T) {
	cache := performance.NewCache(constants.MB * 1).Group("pv")
	const key = "pv:123"

	val, _ := cache.GetAndIncr(key, 1, 0)
	fmt.Println(val)
	val2, _ := cache.GetAndIncr(key, 2, 0)
	val3, _ := cache.GetAndIncr(key, 3, 0)
	fmt.Println(val2)
	fmt.Println(val3)

	r, _ := cache.GetInt64(key)
	fmt.Println(r)

}

func TestTTL(*testing.T) {
	cache := performance.NewCache(constants.MB * 1).Group("pv")
	const key = "k1"

	cache.Set(key, 12, performance.Second(3))
	time.Sleep(time.Second * 2)

	leave, err := cache.GetTTL(key)
	fmt.Println(leave, err == nil)
	time.Sleep(time.Second)
	l2, e2 := cache.GetTTL(key)
	fmt.Println(l2, e2 != nil)
}
