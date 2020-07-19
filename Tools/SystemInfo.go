package Tools

import (
	"github.com/shirou/gopsutil/cpu"
	"runtime"
	"time"
)

func CpuUsagePercent(trackTime time.Duration) float64 {
	percent, _ := cpu.Percent(trackTime,false)
	return percent[0]
}

func CurrentMemUsageMB() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	return bToMb(m.Alloc)
}

func TotalMemUsageMB() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	return bToMb(m.TotalAlloc)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}