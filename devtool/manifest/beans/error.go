package beans

import (
	"github.com/mix-go/bean"
	"github.com/mix-go/console"
)

func Error() {
	Beans = append(Beans,
		bean.BeanDefinition{
			Name:            "error",
			Reflect:         bean.NewReflect(console.NewError),
			Scope:           bean.SINGLETON,
			ConstructorArgs: bean.ConstructorArgs{bean.NewReference("logger")},
		},
	)
}
