package logrus

import (
    "errors"
    "fmt"
    "github.com/sirupsen/logrus"
    "io"
    "os"
    "runtime/debug"
    "testing"
)

func TestLog(t *testing.T) {
    logger := NewLogger()
    logger.Infof("test")
    logger.Infof("test\ntest")
    logger.WithField(logrus.FieldKeyLogrusError, "dfsdfsdf").Infof("test")
}

func TestFile(t *testing.T) {
    logger := NewLogger()

    pwd, _ := os.Getwd()
    file := NewFileWriter(fmt.Sprintf("%s/test.log", pwd), 7)
    writer := io.MultiWriter(os.Stdout, file)
    logger.SetOutput(writer)

    logger.Infof("test")
}

func TestErrorStack(t *testing.T) {
    logger := NewLogger()
    stack := debug.Stack()
    logger.ErrorStack(errors.New("panic test"), &stack)
}
