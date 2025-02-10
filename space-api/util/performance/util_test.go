package performance_test

import (
	"fmt"
	"log"
	"runtime"
	"space-api/constants"
	"space-api/util"
	"space-api/util/performance"
	"sync"
	"testing"
	"time"

	"golang.org/x/sync/singleflight"
)

func init() {
	runtime.GOMAXPROCS(2)
}

func TestSchedSingleflight(t *testing.T) {
	logger := log.Default()
	logger.SetFlags(0)

	t.Run("Test-Normal", func(t *testing.T) {
		var group performance.Group[int]
		var wg sync.WaitGroup
		logger.SetPrefix(constants.CYAN)

		const itCount = 100
		ints := [itCount]int{}

		for i := 0; i < itCount; i++ {
			wg.Add(1)
			go func(j int) {
				defer wg.Done()
				v, _, _ := group.Do("k1", func() (value int, err error) {
					return j, nil
				}, time.Millisecond*10)
				ints[j] = v
			}(i)
		}
		wg.Wait()
		m := map[int]int{}
		for i := 0; i < len(ints); i++ {
			m[ints[i]]++
		}
		for v, c := range m {
			logger.Printf("cached value: %d used: %d nums%s\n", v, c, constants.RESET)
		}
	})

	t.Run("Test-ErrorHandle", func(t *testing.T) {
		var group performance.Group[struct{}]
		var wg sync.WaitGroup

		const itCount = 10
		errs := [itCount]error{}
		logger.SetPrefix(constants.WARN)

		for i := 0; i < itCount; i++ {
			wg.Add(1)
			go func(j int) {
				defer wg.Done()
				_, _, err := group.Do("err", func() (value struct{}, err error) {
					err = &util.BizErr{
						Msg: fmt.Sprintf("test err: %d", j),
					}
					return
				}, time.Millisecond*10)
				errs[j] = err
			}(i)
		}

		wg.Wait()
		for i := 0; i < len(errs); i++ {
			logger.Println(errs[i], constants.RESET)
		}
	})
}

func TestSchedSingleflightPanicHandle(t *testing.T) {
	logger := log.Default()
	logger.SetFlags(0)

	t.Run("Handle-Panic", func(t *testing.T) {
		var g performance.Group[struct{}]
		for i := 0; i < 20; i++ {
			func(i int) {
				defer func() {
					if err := recover(); err != nil {
						logger.Println(constants.BG_RED, err, constants.RESET)
					} else {
						t.Fatal("expect recover a panic error, but got <nil>")
					}
				}()

				g.Do("fatal", func() (value struct{}, err error) {
					panic(fmt.Sprintf("oops...: %d", i))
				})
			}(i)
		}

	})
}

func TestXSingleflight(t *testing.T) {
	logger := log.Default()
	logger.SetFlags(0)

	t.Run("Test-Normal", func(t *testing.T) {
		var group singleflight.Group
		var wg sync.WaitGroup
		logger.SetPrefix(constants.CYAN)

		const itCount = 10
		ints := [itCount]int{}

		for i := 0; i < itCount; i++ {
			wg.Add(1)
			go func(j int) {
				defer wg.Done()
				v, _, _ := group.Do("err", func() (any, error) {
					return j, nil
				})
				ints[j] = v.(int)
			}(i)
		}

		wg.Wait()
		m := map[int]int{}
		for i := 0; i < len(ints); i++ {
			m[ints[i]]++
		}
		for v, c := range m {
			logger.Printf("cached value: %d used: %d nums%s\n", v, c, constants.RESET)
		}
	})

	t.Run("Test-ErrorHandle", func(t *testing.T) {
		var group singleflight.Group
		var wg sync.WaitGroup

		const itCount = 10
		errs := [itCount]error{}
		logger.SetPrefix(constants.RED)

		for i := 0; i < itCount; i++ {
			wg.Add(1)
			go func(j int) {
				defer wg.Done()
				_, err, _ := group.Do("err", func() (any, error) {
					return nil, &util.BizErr{Msg: fmt.Sprintf("%d", j)}
				})
				errs[j] = err
			}(i)
		}

		wg.Wait()
		for i := 0; i < len(errs); i++ {
			logger.Println(errs[i], constants.RESET)
		}
	})
}
