package beans

import (
	"github.com/mix-go/bean"
	"github.com/mix-go/logrus"
)

func Log() {
	Beans = append(Beans,
		bean.BeanDefinition{
			Name:    "logger",
			Reflect: bean.NewReflect(logrus.NewLogger),
			Scope:   bean.SINGLETON,
		},
	)
}
