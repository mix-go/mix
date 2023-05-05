> OpenMix 出品：[https://openmix.org](https://openmix.org/mix-go)

## Mix XEnv

基于 [GoDotEnv](https://github.com/joho/godotenv) 开发的具有**类型转换功能**的环境配置库

Based on GoDotEnv library, with type conversion function

## Usage

载入 `.env` 到环境变量

~~~go
_ = xenv.Load(".env")
~~~

获取环境变量

~~~go
i := xenv.Getenv("key").String()
i := xenv.Getenv("key").Bool()
i := xenv.Getenv("key").Int64()
i := xenv.Getenv("key").Float64()
~~~

设置默认值

~~~go
i := xenv.Getenv("key").String("default")
i := xenv.Getenv("key").Bool(false)
i := xenv.Getenv("key").Int64(123)
i := xenv.Getenv("key").Float64(123.4)
~~~

## License

Apache License Version 2.0, http://www.apache.org/licenses/
