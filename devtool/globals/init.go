package globals

import (
	"os"
)

func Init() {
	// 日志配置
	logger := Logger()
	logger.SetOutput(os.Stdout)
}
