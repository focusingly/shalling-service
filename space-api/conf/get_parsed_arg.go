package conf

import (
	"fmt"
	"log"
	"os"
	"path"
	"space-api/constants"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

type cmdConfig struct {
	configPath string
	useDebug   bool
}

var (
	cmdConfigIns = &cmdConfig{}
)
var fn = sync.OnceFunc(func() {

	rootCmd := &cobra.Command{
		Use:          "配置应用服务",
		Short:        "应用服务配置文件路径设置",
		SilenceUsage: false,
	}

	flagSet := rootCmd.PersistentFlags()
	flagSet.StringVarP(&cmdConfigIns.configPath, "config", "c", "", "the project option config")
	flagSet.BoolVarP(&cmdConfigIns.useDebug, "use-debug", "d", false, "set debug mode")

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}

	cmdConfigIns.configPath = strings.TrimSpace(cmdConfigIns.configPath)
	if cmdConfigIns.configPath == "" {
		t := fmt.Sprintf(
			"%snot config set, service use default configuration%s",
			constants.CYAN,
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

	return cmdConfigIns.configPath, cmdConfigIns.useDebug
}
