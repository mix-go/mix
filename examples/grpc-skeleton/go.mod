module github.com/mix-go/grpc-skeleton

go 1.13

replace (
	github.com/mix-go/dotenv => ../../src/dotenv
	github.com/mix-go/xcli => ../../src/xcli
	github.com/mix-go/xdi => ../../src/xdi
)

require (
	github.com/go-redis/redis/v8 v8.8.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/jinzhu/configor v1.2.0
	github.com/mix-go/dotenv v1.1.15
	github.com/mix-go/xcli v1.1.15
	github.com/mix-go/xdi v1.1.11
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sirupsen/logrus v1.6.0
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0
	google.golang.org/genproto v0.0.0-20191009194640-548a555dbc03 // indirect
	google.golang.org/grpc v1.30.0
	google.golang.org/protobuf v1.23.0
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gorm.io/driver/mysql v1.0.5
	gorm.io/gorm v1.21.3
	xorm.io/xorm v1.0.7
)
