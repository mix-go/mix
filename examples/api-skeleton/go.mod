module github.com/mix-go/api-skeleton

go 1.13

replace (
	github.com/mix-go/dotenv => ../../src/dotenv
	github.com/mix-go/xcli => ../../src/xcli
	github.com/mix-go/xdi => ../../src/xdi
	github.com/mix-go/xsql => ../../src/xsql
)

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/bytedance/sonic v1.8.8 // indirect
	github.com/gin-gonic/gin v1.9.0
	github.com/go-playground/validator/v10 v10.13.0 // indirect
	github.com/go-redis/redis v6.15.9+incompatible // indirect
	github.com/go-session/redis v3.0.1+incompatible
	github.com/go-session/session v3.1.2+incompatible
	github.com/go-sql-driver/mysql v1.7.1
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0
	github.com/jinzhu/configor v1.2.1
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/klauspost/cpuid/v2 v2.2.4 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/mattn/go-isatty v0.0.18 // indirect
	github.com/mix-go/dotenv v1.1.15
	github.com/mix-go/xcli v1.1.21
	github.com/mix-go/xdi v1.1.17
	github.com/mix-go/xsql v1.1.11
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.18.1 // indirect
	github.com/pelletier/go-toml/v2 v2.0.7 // indirect
	github.com/redis/go-redis/v9 v9.0.4
	github.com/sijms/go-ora/v2 v2.7.3 // indirect
	github.com/sirupsen/logrus v1.9.0
	github.com/spf13/afero v1.9.5 // indirect
	github.com/spf13/viper v1.15.0
	github.com/ugorji/go/codec v1.2.11 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.24.0
	golang.org/x/arch v0.3.0 // indirect
	golang.org/x/crypto v0.8.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
	gorm.io/driver/mysql v1.5.0
	gorm.io/gorm v1.25.0
	xorm.io/builder v0.3.12 // indirect
	xorm.io/xorm v1.3.2
)
