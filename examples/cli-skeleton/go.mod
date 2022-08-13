module github.com/mix-go/cli-skeleton

go 1.14

replace (
	github.com/mix-go/dotenv => ../../src/dotenv
	github.com/mix-go/xcli => ../../src/xcli
	github.com/mix-go/xdi => ../../src/xdi
	github.com/mix-go/xsql => ../../src/xsql
)

require (
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/go-redis/redis/v8 v8.7.1
	github.com/go-sql-driver/mysql v1.6.0
	github.com/jinzhu/configor v1.2.1
	github.com/mix-go/dotenv v1.1.15
	github.com/mix-go/xcli v1.1.20
	github.com/mix-go/xdi v1.1.16
	github.com/mix-go/xsql v1.1.8
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cast v1.4.1 // indirect
	github.com/spf13/viper v1.8.1
	go.uber.org/zap v1.17.0
	golang.org/x/sys v0.0.0-20210903071746-97244b99971b // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gorm.io/driver/mysql v1.0.5
	gorm.io/gorm v1.21.3
	xorm.io/xorm v1.0.7
)
