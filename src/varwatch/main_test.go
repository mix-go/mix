package varwatch

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var Config struct {
	Logger struct {
		Level string `json:"level"`
	} `json:"logger" varwatch:"logger"`
	Database struct {
		User    string `json:"user"`
		Pwd     string `json:"pwd"`
		Db      string `json:"db"`
		MaxOpen int    `json:"max_open"`
		MaxIdle int    `json:"max_idle"`
	} `json:"database" varwatch:"database"`
}

func TestNewWatcher(t *testing.T) {
	a := assert.New(t)

	w, err := NewWatcher(&Config)
	if err != nil {
		panic(err)
	}
	err = w.Watch("logger", func() {
		a.Equal(Config.Logger.Level, "debug")
	})
	err = w.Watch("database", func() {
		a.Equal(Config.Database.MaxOpen, 100)
		a.Equal(Config.Database.MaxIdle, 50)
	})
	if err != nil {
		panic(err)
	}
	w.Run(100 * time.Millisecond)

	go func() {
		time.Sleep(1 * time.Second)
		Config.Logger.Level = "debug"
		Config.Database.MaxOpen = 100
		Config.Database.MaxIdle = 50
	}()

	time.Sleep(2 * time.Second)
}
