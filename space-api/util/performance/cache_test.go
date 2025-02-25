package performance_test

import (
	"fmt"
	"space-api/util/performance"
	"testing"
	"time"
)

func TestCacheIncrOperation(t *testing.T) {
	cache := performance.NewBuntDBCache(":memory:").Group("pv")
	defer cache.Close()

	const key = "k1"
	val, _ := cache.GetAndIncr(key, 2, time.Second*2)
	fmt.Println(val)
	time.Sleep(time.Second * 1)
	val2, _ := cache.IncrAndGet(key, 3)

	d, _ := cache.GetTTL(key)
	fmt.Println(val2, d)
}

func TestCacheIncrOperation2(*testing.T) {
	cache := performance.NewBuntDBCache(":memory:").Group("pv")
	defer cache.Close()
	const key = "k1"

	val, _ := cache.GetAndIncr(key, 1)
	fmt.Println(val)
	val2, _ := cache.GetAndIncr(key, 2)
	val3, _ := cache.GetAndIncr(key, 3)
	fmt.Println(val2)
	fmt.Println(val3)

	r, _ := cache.GetInt64(key)
	fmt.Println(r)

}

func TestTTL(*testing.T) {
	cache := performance.NewBuntDBCache(":memory:").Group("pv")
	defer cache.Close()
	const key = "k1"

	cache.Set(key, 12, time.Second*3)
	time.Sleep(time.Second * 2)

	leave, err := cache.GetTTL(key)
	fmt.Println(leave, err == nil)
	time.Sleep(time.Second)
	l2, e2 := cache.GetTTL(key)
	fmt.Println(l2, e2 != nil)
}
