package di

import (
	"github.com/go-redis/redis/v8"
	"github.com/mix-go/dotenv"
	"github.com/mix-go/xdi"
	"time"
)

func init() {
	obj := xdi.Object{
		Name: "goredis",
		New: func() (i interface{}, e error) {
			opt := redis.Options{
				Addr:        dotenv.Getenv("REDIS_ADDR").String(),
				Password:    dotenv.Getenv("REDIS_PASSWORD").String(),
				DB:          int(dotenv.Getenv("REDIS_DATABASE").Int64()),
				DialTimeout: time.Duration(dotenv.Getenv("REDIS_DIAL_TIMEOUT").Int64(10)) * time.Second,
			}
			return redis.NewClient(&opt), nil
		},
	}
	if err := xdi.Provide(&obj); err != nil {
		panic(err)
	}
}

func GoRedis() (client *redis.Client) {
	if err := xdi.Populate("goredis", &client); err != nil {
		panic(err)
	}
	return
}
