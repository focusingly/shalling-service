package performance_test

import (
	"fmt"
	"space-api/constants"
	"space-api/util/performance"
	"testing"
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
