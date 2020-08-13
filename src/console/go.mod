module github.com/mix-go/console

go 1.13

replace (
	github.com/mix-go/bean => ../bean
	github.com/mix-go/event => ../event
	github.com/mix-go/logrus => ../logrus
)

require (
	github.com/astaxie/beego v1.12.2
	github.com/mix-go/bean v1.0.0-beta4
	github.com/mix-go/event v1.0.0-beta4
	github.com/mix-go/logrus v1.0.0-beta4
	github.com/stretchr/testify v1.6.1
)
