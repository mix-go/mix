module github.com/mix-go/console

go 1.13

replace (
	github.com/mix-go/bean => ../bean
	github.com/mix-go/event => ../event
	github.com/mix-go/logrus => ../logrus
)

require (
	github.com/mix-go/bean v1.0.16
	github.com/mix-go/event v1.0.16
	github.com/mix-go/logrus v1.0.19
	github.com/stretchr/testify v1.6.1
)
