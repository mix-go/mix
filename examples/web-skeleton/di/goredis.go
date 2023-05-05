package di

import (
	"github.com/mix-go/xdi"
	"github.com/mix-go/xutil/xenv"
	"github.com/redis/go-redis/v9"
	"time"
)

func init() {
	obj := xdi.Object{
		Name: "goredis",
		New: func() (i interface{}, e error) {
			opt := redis.Options{
				Addr:        xenv.Getenv("REDIS_ADDR").String(),
				Password:    xenv.Getenv("REDIS_PASSWORD").String(),
				DB:          int(xenv.Getenv("REDIS_DATABASE").Int64()),
				DialTimeout: time.Duration(xenv.Getenv("REDIS_DIAL_TIMEOUT").Int64(10)) * time.Second,
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
