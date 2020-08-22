module github.com/mix-go/mix

go 1.13

replace (
	github.com/mix-go/bean => ./src/bean
	github.com/mix-go/console => ./src/console
	github.com/mix-go/dotenv => ./src/dotenv
	github.com/mix-go/event => ./src/event
	github.com/mix-go/logrus => ./src/logrus
)

require (
	github.com/mix-go/bean v1.0.0-beta.10
	github.com/mix-go/console v1.0.0-beta.12
	github.com/mix-go/dotenv v1.0.0-beta.10
	github.com/mix-go/event v1.0.0-beta.10
	github.com/mix-go/logrus v1.0.0-beta.10
	golang.org/x/sys v0.0.0-20200821140526-fda516888d29 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
	gorm.io/gorm v0.2.34 // indirect
)
