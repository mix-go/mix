package xhttp

import (
	"sync"
	"sync/atomic"
)

var shutdownController = NewShutdownController()

type ShutdownController struct {
	shutdownFlag int32          // 原子标记，表示是否进入停机状态
	waitGroup    sync.WaitGroup // 用于等待所有处理中的请求完成
}

func NewShutdownController() *ShutdownController {
	return &ShutdownController{}
}

func (sc *ShutdownController) BeginRequest() bool {
	if atomic.LoadInt32(&sc.shutdownFlag) == 1 {
		return false
	}
	sc.waitGroup.Add(1)
	return true
}

func (sc *ShutdownController) EndRequest() {
	sc.waitGroup.Done()
}

func (sc *ShutdownController) InitiateShutdown() {
	atomic.StoreInt32(&sc.shutdownFlag, 1)
	sc.waitGroup.Wait()
}

func Shutdown() {
	shutdownController.InitiateShutdown()
}
