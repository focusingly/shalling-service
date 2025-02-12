package conf

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"space-api/constants"
	"strings"
	"sync"
)

var (
	_confPath    string
	_isDebugMode bool
)
var fn = sync.OnceFunc(func() {
	flag.StringVar(&_confPath, "c", "", "the project option config")
	flag.BoolVar(&_isDebugMode, "use-debug", false, "set debug mode")
	flag.Parse()

	if strings.TrimSpace(_confPath) == "" {
		t := fmt.Sprintf(
			"%snot config set, service use default configuration%s",
			constants.BG_CYAN,
			constants.RESET,
		)
		fmt.Println(t)

		if err := os.MkdirAll(_defaultStore, os.ModePerm); err != nil {
			log.Fatal("create default store error: ", err)
		}
		if err := os.MkdirAll(path.Join(_defaultStore, "db"), os.ModePerm); err != nil {
			log.Fatal("create store error: ", err)
		}
		if err := os.MkdirAll(path.Join(_defaultStore, "files"), os.ModePerm); err != nil {
			log.Fatal("create store error: ", err)
		}

		return
	}
})

func GetParsedArgs() (confPath string, isDebugMode bool) {
	fn()

	return _confPath, _isDebugMode
}
