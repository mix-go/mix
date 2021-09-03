module github.com/mix-go/cli-skeleton

go 1.14

replace (
	github.com/mix-go/dotenv => ../../src/dotenv
	github.com/mix-go/xcli => ../../src/xcli
	github.com/mix-go/xdi => ../../src/xdi
)

require (
	github.com/go-redis/redis/v8 v8.7.1
	github.com/go-sql-driver/mysql v1.5.0
	github.com/jinzhu/configor v1.2.1
	github.com/mix-go/dotenv v1.1.15
	github.com/mix-go/xcli v1.1.15
	github.com/mix-go/xdi v1.1.11
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sirupsen/logrus v1.8.1
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gorm.io/driver/mysql v1.0.5
	gorm.io/gorm v1.21.3
	xorm.io/xorm v1.0.7
)
