package Tools

import (
	"github.com/shirou/gopsutil/cpu"
	"os"
	"path/filepath"
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

func WalkMatch(root, pattern string) ([]string, error) {
	// Copied from
	// https://stackoverflow.com/questions/55300117/how-do-i-find-all-files-that-have-a-certain-extension-in-go-regardless-of-depth
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}