package performance

import (
	"github.com/bytedance/gopkg/util/gopool"
)

// DefaultTaskRunner 默认的全局任务执行器, 默认并行度: 2
var DefaultTaskRunner = gopool.NewPool("task-runner", 2, gopool.NewConfig())
