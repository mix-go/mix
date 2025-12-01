package xutil

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
)

var statsSingleton *MethodStats
var statsOnce sync.Once

func getMethodStatsSingleton(logger Logger) *MethodStats {
	statsOnce.Do(func() {
		statsSingleton = NewMethodStats(logger)
	})
	return statsSingleton
}

// MethodStats 存储每个 gRPC 方法的调用计数，使用 sync.Map 和原子操作实现高性能并发访问。
type MethodStats struct {
	counts *sync.Map // 键: string, 值: *int64
	logger Logger
	// 用于记录上一次重置的时间，方便日志输出时间跨度
	lastResetTime atomic.Value
}

// NewMethodStats 创建并返回一个新的 MethodStats 实例。
func NewMethodStats(logger Logger) *MethodStats {
	if logger == nil {
		logger = &defaultLogger{}
	}
	ms := &MethodStats{
		counts: &sync.Map{},
		logger: logger,
	}
	ms.lastResetTime.Store(time.Now()) // 首次创建时初始化上次重置时间
	return ms
}

// Inc 方法：原子性地增加指定方法的调用计数。
func (ms *MethodStats) Inc(fullMethodName string) {
	actual, _ := ms.counts.LoadOrStore(fullMethodName, new(int64))
	counterPtr := actual.(*int64)
	atomic.AddInt64(counterPtr, 1) // 高性能原子增量操作
}

// GetSnapshot 方法：获取当前的统计数据，但**不重置**计数器。
func (ms *MethodStats) GetSnapshot() map[string]uint64 {
	snapshot := make(map[string]uint64)

	ms.counts.Range(func(key, value interface{}) bool {
		methodName := key.(string)
		counterPtr := value.(*int64)

		// 原子地加载当前值
		currentValue := atomic.LoadInt64(counterPtr)

		if currentValue > 0 {
			snapshot[methodName] = uint64(currentValue)
		}
		return true
	})

	return snapshot
}

// Reset 方法：原子地清零所有计数器，并更新重置时间。
func (ms *MethodStats) Reset() {
	duration := ms.GetIntervalDuration()
	ms.logger.Debugw(fmt.Sprintf("Grpc Monitor: gRPC Method Call Daily Reset at %v. No calls recorded in the last %v.", time.Now().Format("15:04:05"), duration.Truncate(time.Second)))

	now := time.Now()

	ms.counts.Range(func(key, value interface{}) bool {
		counterPtr := value.(*int64)
		// 原子地交换旧值并将其设置为 0
		atomic.SwapInt64(counterPtr, 0)
		return true
	})

	// 更新上次重置时间
	ms.lastResetTime.Store(now)
}

// GetIntervalDuration 获取从上次重置到当前的时间跨度
func (ms *MethodStats) GetIntervalDuration() time.Duration {
	lastTime := ms.lastResetTime.Load().(time.Time)
	return time.Since(lastTime)
}

func StatsUnaryServerInterceptor(logger Logger) grpc.UnaryServerInterceptor {
	stats := getMethodStatsSingleton(logger)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		stats.Inc(info.FullMethod)
		resp, err := handler(ctx, req)
		return resp, err
	}
}

func StatsStreamServerInterceptor(logger Logger) grpc.StreamServerInterceptor {
	stats := getMethodStatsSingleton(logger)
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		stats.Inc(info.FullMethod)
		err := handler(srv, ss)
		return err
	}
}

func StatsServerOptions(logger Logger) []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.UnaryInterceptor(StatsUnaryServerInterceptor(logger)),
		grpc.StreamInterceptor(StatsStreamServerInterceptor(logger)),
	}
}

// calculateNextMidnight 计算从当前时间到下一个午夜的时间间隔。
func calculateNextMidnight(now time.Time) time.Duration {
	// 获取今天的午夜
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// 下一个午夜是今天的午夜 + 24小时
	nextMidnight := midnight.Add(24 * time.Hour)

	// 返回时间差
	return nextMidnight.Sub(now)
}

// StartGrpcMonitoring 启动两个独立的定时器：
// 1. Snapshot 定时器：按 interval 周期性地打印统计快照 (不重置)。
// 2. Midnight 重置器：每天零点准时清零计数器。
func StartGrpcMonitoring(interval time.Duration, logger Logger) {
	stats := getMethodStatsSingleton(logger)
	stats.logger = logger // 设置 logger

	if interval <= 0 {
		stats.logger.Errorw("Grpc Monitor: Snapshot interval is zero or negative. Only daily reset will be performed.")
	} else {
		// --- 1. 启动 Snapshot 定时器 (仅打印，不重置) ---
		go func() {
			ticker := time.NewTicker(interval)
			stats.logger.Debugw(fmt.Sprintf("Grpc Monitor: Snapshot timer started with interval %v.", interval))
			for range ticker.C {
				// 打印统计快照，但不清零计数器
				stats.PrintStatsSnapshot()
			}
		}()
	}

	// --- 2. 启动 Midnight 重置器 (每天零点清零) ---

	// 第一次触发：计算到下一个零点的时长
	initialDelay := calculateNextMidnight(time.Now())
	stats.logger.Debugw(fmt.Sprintf("Grpc Monitor: Next midnight reset scheduled for %v (delay: %v).", time.Now().Add(initialDelay).Format("2006-01-02 15:04:05"), initialDelay))

	// 使用 time.AfterFunc 在 initialDelay 后执行第一次重置
	time.AfterFunc(initialDelay, func() {
		// 第一次重置操作
		stats.Reset()

		// 后续触发：启动一个 24 小时的 Ticker
		ticker := time.NewTicker(24 * time.Hour)
		stats.logger.Debugw("Grpc Monitor: Daily reset Ticker started (24h).")
		for range ticker.C {
			stats.Reset()
		}
	})
}

// PrintStatsSnapshot 仅打印当前快照，不重置计数。
func (ms *MethodStats) PrintStatsSnapshot() {
	snapshot := ms.GetSnapshot() // 获取快照，不重置

	if len(snapshot) == 0 {
		return
	}

	ms.logger.Debugw("Grpc Monitor", "snapshot", snapshot)
}

type defaultLogger struct{}

func (*defaultLogger) Debugw(msg string, keysAndValues ...interface{}) {
}

func (*defaultLogger) Errorw(msg string, keysAndValues ...interface{}) {
}
