## Mix Dotenv

基于 [GoDotEnv](https://github.com/joho/godotenv) 开发的具有**类型转换功能**的环境配置库

Based on GoDotEnv library, with type conversion function

## Usage

- 安装

```
go get -u github.com/mix-go/dotenv
```

- 使用

~~~
_ = dotenv.Load(".env")
i := dotenv.Getenv("key").String()
i := dotenv.Getenv("key").Bool()
i := dotenv.Getenv("key").Int64()
i := dotenv.Getenv("key").Float64()
~~~

## License

Apache License Version 2.0, http://www.apache.org/licenses/
