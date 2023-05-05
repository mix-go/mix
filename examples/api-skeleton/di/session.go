package di

import (
	"github.com/go-session/redis"
	"github.com/go-session/session"
	"github.com/mix-go/xdi"
	"github.com/mix-go/xutil/xenv"
	"time"
)

func init() {
	obj := xdi.Object{
		Name: "session",
		New: func() (i interface{}, e error) {
			opts := redis.Options{
				Addr:        xenv.Getenv("REDIS_ADDR").String(),
				Password:    xenv.Getenv("REDIS_PASSWORD").String(),
				DB:          int(xenv.Getenv("REDIS_DATABASE").Int64()),
				DialTimeout: time.Duration(xenv.Getenv("REDIS_DIAL_TIMEOUT").Int64(10)) * time.Second,
			}
			opt := redis.NewRedisStore(&opts)
			return session.NewManager(session.SetStore(opt)), nil
		},
	}
	if err := xdi.Provide(&obj); err != nil {
		panic(err)
	}
}
