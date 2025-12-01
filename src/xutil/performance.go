package xutil

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/process"
	"os"
	"runtime"
	"runtime/debug"
	"time"
)

// PerformanceStats 性能监控数据结构
type PerformanceStats struct {
	Timestamp         string  `json:"timestamp"`           // 时间戳
	AllocMB           float64 `json:"alloc_mb"`            // 当前堆对象占用的内存
	TotalAllocMB      float64 `json:"total_alloc_mb"`      // 累计分配的内存
	SysMB             float64 `json:"sys_mb"`              // 从系统获取的总内存
	HeapAllocMB       float64 `json:"heap_alloc_mb"`       // 堆对象占用内存
	HeapSysMB         float64 `json:"heap_sys_mb"`         // 堆从系统获取的内存
	HeapIdleMB        float64 `json:"heap_idle_mb"`        // 堆中空闲内存
	HeapInuseMB       float64 `json:"heap_inuse_mb"`       // 堆中使用中的内存
	NumGC             uint32  `json:"num_gc"`              // GC执行次数
	LastGC            uint64  `json:"last_gc"`             // 上次GC时间
	PauseTotalMs      uint64  `json:"pause_total_ms"`      // GC总暂停时间(毫秒)
	GCCPUFraction     float64 `json:"gc_cpu_fraction"`     // GC占用CPU时间比例
	Goroutines        int     `json:"goroutines"`          // 当前协程总数
	NumCPU            int     `json:"num_cpu"`             // CPU核心数
	GOMAXPROCS        int     `json:"gomaxprocs"`          // 最大并行执行的CPU数
	SystemCPUPercent  float64 `json:"system_cpu_percent"`  // 系统总体CPU使用率
	ProcessCPUPercent float64 `json:"process_cpu_percent"` // 当前进程CPU使用率
	ProcessMemoryMB   float64 `json:"process_memory_mb"`   // 进程物理内存使用量
	ProcessThreads    int32   `json:"process_threads"`     // 进程线程数
}

func logPerformanceStats(logger Logger, handler func(*PerformanceStats)) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	goroutineCount := runtime.NumGoroutine()

	// 获取CPU和进程信息
	cpuCount := runtime.NumCPU()
	maxProcs := runtime.GOMAXPROCS(0)
	systemCPUPercent := getSystemCPUPercent()
	processCPUPercent, processMemoryMB, processThreads := getProcessInfo()

	stats := PerformanceStats{
		Timestamp:         time.Now().Format("2006-01-02 15:04:05"),
		AllocMB:           bToMb(m.Alloc),
		TotalAllocMB:      bToMb(m.TotalAlloc),
		SysMB:             bToMb(m.Sys),
		HeapAllocMB:       bToMb(m.HeapAlloc),
		HeapSysMB:         bToMb(m.HeapSys),
		HeapIdleMB:        bToMb(m.HeapIdle),
		HeapInuseMB:       bToMb(m.HeapInuse),
		NumGC:             m.NumGC,
		LastGC:            m.LastGC,
		PauseTotalMs:      m.PauseTotalNs / 1000000,
		GCCPUFraction:     m.GCCPUFraction,
		Goroutines:        goroutineCount,
		NumCPU:            cpuCount,
		GOMAXPROCS:        maxProcs,
		SystemCPUPercent:  systemCPUPercent,
		ProcessCPUPercent: processCPUPercent,
		ProcessMemoryMB:   processMemoryMB,
		ProcessThreads:    processThreads,
	}

	logger.Debugw("Performance Monitor", "stats", stats)

	if handler != nil {
		handler(&stats)
	}
}

// 获取系统CPU使用率
func getSystemCPUPercent() float64 {
	percentages, err := cpu.Percent(time.Second, false)
	if err != nil || len(percentages) == 0 {
		return 0.0
	}
	return percentages[0]
}

// 获取进程信息（CPU使用率、内存、线程数）
func getProcessInfo() (float64, float64, int32) {
	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return 0.0, 0.0, 0
	}

	// 获取CPU使用率
	cpuPercent, err := p.CPUPercent()
	if err != nil {
		cpuPercent = 0.0
	}

	// 获取内存信息
	memInfo, err := p.MemoryInfo()
	var memoryMB float64
	if err == nil && memInfo != nil {
		memoryMB = float64(memInfo.RSS) / 1024 / 1024
	}

	// 获取线程数
	threads, err := p.NumThreads()
	if err != nil {
		threads = 0
	}

	return cpuPercent, memoryMB, threads
}

func bToMb(b uint64) float64 {
	return float64(b) / 1024 / 1024
}

// ForceGC 强制 GC 并释放内存
func ForceGC() {
	runtime.GC()
	debug.FreeOSMemory()
}

// StartPerformanceMonitoring 启动性能监控
func StartPerformanceMonitoring(interval time.Duration, logger Logger, handler func(*PerformanceStats)) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			logPerformanceStats(logger, handler)
		}
	}()
}
