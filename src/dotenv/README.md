> OpenMix 出品：[https://openmix.org](https://openmix.org/mix-go)

## Mix Dotenv

基于 [GoDotEnv](https://github.com/joho/godotenv) 开发的具有**类型转换功能**的环境配置库

Based on GoDotEnv library, with type conversion function

## Installation

```
go get github.com/mix-go/dotenv
```

## Usage

载入 `.env` 到环境变量

~~~go
_ = dotenv.Load(".env")
~~~

获取环境变量

~~~go
i := dotenv.Getenv("key").String()
i := dotenv.Getenv("key").Bool()
i := dotenv.Getenv("key").Int64()
i := dotenv.Getenv("key").Float64()
~~~

设置默认值

~~~go
i := dotenv.Getenv("key").String("default")
i := dotenv.Getenv("key").Bool(false)
i := dotenv.Getenv("key").Int64(123)
i := dotenv.Getenv("key").Float64(123.4)
~~~

## License

Apache License Version 2.0, http://www.apache.org/licenses/
