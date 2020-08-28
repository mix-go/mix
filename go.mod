module github.com/mix-go/mix

go 1.13

replace (
	github.com/mix-go/bean => ./src/bean
	github.com/mix-go/console => ./src/console
	github.com/mix-go/event => ./src/event
	github.com/mix-go/logrus => ./src/logrus
)

require (
	github.com/mix-go/bean v1.0.4
	github.com/mix-go/console v1.0.5
	github.com/mix-go/dotenv v1.0.1
	github.com/mix-go/event v1.0.1
	github.com/mix-go/logrus v1.0.1
)
