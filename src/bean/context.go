package bean

import (
	"fmt"
	"sync"
)

// 创建实例
func NewApplicationContext(definitions []BeanDefinition) *ApplicationContext {
	context := &ApplicationContext{Definitions: definitions}
	context.Init()
	return context
}

// 应用上下文
type ApplicationContext struct {
	Definitions     []BeanDefinition
	tidyDefinitions sync.Map
	instances       sync.Map
}

// 初始化
func (t *ApplicationContext) Init() {
	t.tidyDefinitions = sync.Map{}
	for _, d := range t.Definitions {
		d.Context = t
		t.tidyDefinitions.Store(d.Name, d)
	}
}

// 获取定义
func (t *ApplicationContext) GetDefinition(name string) BeanDefinition {
	var (
		inf interface{}
		ok  bool
	)
	if inf, ok = t.tidyDefinitions.Load(name); !ok {
		panic(fmt.Sprintf("Bean not found: %s", name))
	}
	return inf.(BeanDefinition)
}

// 获取实例
func (t *ApplicationContext) GetBean(name string, prop Fields, args ConstructorArgs) interface{} {
	def := merge(t.GetDefinition(name), prop, args)
	if def.Scope == SINGLETON {
		if ins, ok := t.instances.Load(name); ok {
			return ins
		}
		ins := def.Instance()
		t.instances.Store(name, ins)
		return ins
	}
	return def.Instance()
}

// 快速获取实例
func (c *ApplicationContext) Get(name string) interface{} {
	return c.GetBean(name, Fields{}, ConstructorArgs{})
}

// 合并
// args | fields 内的字段会替换之前定义的值
// args 内的 nil 值将会忽略，不会替换处理
func merge(def BeanDefinition, fields Fields, args ConstructorArgs) BeanDefinition {
	iFields := len(fields) > 0
	iArgs := len(args) > 0
	if iFields || iArgs {
		newDef := BeanDefinition{
			Name:            def.Name,
			Scope:           def.Scope,
			Reflect:         def.Reflect,
			InitMethod:      def.InitMethod,
			ConstructorArgs: nil,
			Fields:          nil,
			Context:         def.Context,
		}
		if iFields {
			// 合并替换字段
			tmp := Fields{}
			for k, v := range def.Fields {
				tmp[k] = v
			}
			for k, v := range fields {
				tmp[k] = v
			}
			newDef.Fields = tmp
		}
		if iArgs {
			// 合并替换参数，nil 忽略
			tmp := ConstructorArgs{}
			tmp = append(tmp, def.ConstructorArgs...)
			for k, v := range args {
				if v == nil {
					continue
				}
				ok := false
				for sk, _ := range def.ConstructorArgs {
					if sk == k {
						ok = true
					}
				}
				if ok {
					tmp[k] = v
				} else {
					tmp = append(tmp, v)
				}
			}
			newDef.ConstructorArgs = tmp
		}
		return newDef
	}
	return def
}
