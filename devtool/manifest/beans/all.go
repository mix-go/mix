package beans

import "github.com/mix-go/bean"

var (
	Beans []bean.BeanDefinition
)

func Init() {
	// 因为 beans 中会使用到 .env conf 中的配置信息
	// 因此只能 Init 手动加载，不能使用 init
	Error()
	Event()
	Log()
}
