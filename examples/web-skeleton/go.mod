module github.com/mix-go/web-skeleton

go 1.13

replace (
	github.com/mix-go/dotenv => ../../src/dotenv
	github.com/mix-go/xcli => ../../src/xcli
	github.com/mix-go/xdi => ../../src/xdi
)

require (
	github.com/gin-gonic/gin v1.7.4
	github.com/go-playground/validator/v10 v10.9.0 // indirect
	github.com/go-redis/redis v6.15.9+incompatible // indirect
	github.com/go-redis/redis/v8 v8.8.0
	github.com/go-session/redis v3.0.1+incompatible
	github.com/go-session/session v3.1.2+incompatible
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/snappy v0.0.3 // indirect
	github.com/gorilla/websocket v1.4.2
	github.com/jinzhu/configor v1.2.0
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/mix-go/dotenv v1.1.15
	github.com/mix-go/xcli v1.1.15
	github.com/mix-go/xdi v1.1.11
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/viper v1.8.1
	github.com/ugorji/go v1.2.6 // indirect
	go.uber.org/zap v1.17.0
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5 // indirect
	golang.org/x/sys v0.0.0-20210903071746-97244b99971b // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gorm.io/driver/mysql v1.0.5
	gorm.io/gorm v1.21.6
	xorm.io/builder v0.3.9 // indirect
	xorm.io/xorm v1.0.7
)
