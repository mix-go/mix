module github.com/mix-go/cli

go 1.13

replace (
	github.com/mix-go/bean => ../bean
	github.com/mix-go/event => ../event
	github.com/mix-go/logrus => ../logrus
)

require (
	github.com/mix-go/bean v0.0.0-00010101000000-000000000000
	github.com/mix-go/event v0.0.0-00010101000000-000000000000
	github.com/mix-go/logrus v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.6.1
)
