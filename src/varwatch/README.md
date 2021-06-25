> OpenMix 出品：[https://openmix.org](https://openmix.org/mix-go)

## Mix VarWatch

监视配置结构体变量的数据变化并执行一些任务

Monitor the data changes of configuration structure variables and perform some tasks

## Installation

```
go get github.com/mix-go/varwatch
```
## Usage

当采用 [spf13/viper](https://github.com/spf13/viper) [jinzhu/configor](https://github.com/jinzhu/configor) 这种绑定变量的配置库来动态更新配置信息

> 任何采用 &Config 指针绑定数据的配置库都可以

~~~go
var Config struct {
	Logger struct {
		Level int `json:"level"`
	} `json:"logger" varwatch:"logger"`
	Database struct {
		User    string `json:"user"`
		Pwd     string `json:"pwd"`
		Db      string `json:"db"`
		MaxOpen int    `json:"max_open"`
		MaxIdle int    `json:"max_idle"`
	} `json:"database" varwatch:"database"`
}

err := viper.Unmarshal(&Config)
~~~

以动态修改日志级别举例：当 `Config.Logger.Level` 发生变化时我们需要执行一些代码修改日志的级别

 - 首先将 Logger 节点配置 `varwatch:"logger"` 标签信息
 - 然后采用以下代码执行监听逻辑

```go
w, err := varwatch.NewWatcher(&Config, 10 * time.Second)
if err != nil {
    panic(err)
}
w.Watch("logger", func() {
    // 获取变化后的值
    lv := Config.Logger.Level
    // 修改 logrus 的日志级别
    logrus.SetLevel(logrus.Level(uint32(lv)))
})
```

需要动态修改连接池信息，或者数据库账号密码都可以通过上面的范例实现。

## License

Apache License Version 2.0, http://www.apache.org/licenses/
