# 编程须知

## 编译与执行

- 不可直接 `go run main.go`

由于骨架代码有读取 `.env`、`conf/config.yml` 的配置文件，而配置文件的路径采用的是编译的二进制文件所在的 `bin` 目录的相对路径，由于 `go run main.go` 生成的二进制文件路径是由go编译器随机指定在 `/var/folders/` 目录底下一个随机临时目录，此时依据相对路径是无法成功读取配置文件的，因此会抛出 `panic` 异常提示找不到配置文件。

```
// 切换到项目根目录
cd project

// 只能这样编译到 project/bin 目录
go build -o bin/go_build_main_go main.go

// 执行编译好的程序
bin/go_build_main_go api
```

上面的命令可以合并执行

```
go build -o bin/go_build_main_go main.go && bin/go_build_main_go api
```

## `Goland` 如何使用

使用的 `Goland` 请阅读以下文章，学习更加快速的编译方法

- [MixGo 在 IDE Goland 中的如何使用](https://zhuanlan.zhihu.com/p/391857663)

