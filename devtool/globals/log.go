package globals

import (
    "github.com/mix-go/console"
    "github.com/mix-go/logrus"
)

func Logger() *logrus.Logger {
    return console.Context().Get("logger").(*logrus.Logger)
}
