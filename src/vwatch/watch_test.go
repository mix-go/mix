package vwatch

import (
	"fmt"
	"testing"
	"time"
)

var Config struct {
	AppDebug bool
	Host     struct {
		Address string `json:"address"`
		Ports   []int  `json:"ports"`
	} `json:"host"`
	Datastore struct {
		Metric struct {
			Host string `json:"host"`
			Port int    `json:"port"`
		} `json:"metric"`
		Warehouse struct {
			Host string `json:"host"`
			Port int    `json:"port"`
		} `json:"warehouse"`
	} `json:"datastore"`
}

func init() {
	Config.AppDebug = true
	Config.Host.Address = "127.0.0.1"
	Config.Datastore.Warehouse.Host = "127.0.0.1"
}

func TestNewWatcher(t *testing.T) {
	w, err := NewWatcher(&Config)
	if err != nil {
		panic(err)
	}
	err = w.Watch(&Config.Host.Address, func() {
		addr := Config.Host.Address
		fmt.Println(addr)
	})
	if err != nil {
		panic(err)
	}
	w.Run(100 * time.Millisecond)

	go func() {
		time.Sleep(1 * time.Second)
		Config.AppDebug = false
		Config.Host.Address = "127.0.0.2"
		Config.Datastore.Warehouse.Host = "127.0.0.2"
	}()

	time.Sleep(5 * time.Second)
}
