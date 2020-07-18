package Tools

import (
	"github.com/shirou/gopsutil/cpu"
	"runtime"
	"time"
)

func CpuUsagePercent() float64 {
	percent, _ := cpu.Percent(time.Second,true)
	return percent[0]
	//clockSeconds := float64(C.Clock()-startTicks) / float64(C.CLOCKS_PER_SEC)
	//realSeconds := time.Since(startTime).Seconds()
	//return clockSeconds / realSeconds * 100
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