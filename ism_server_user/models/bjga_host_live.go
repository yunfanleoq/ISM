package models

import (
	"math"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

var (
	bjgaHostStart = time.Now()
	bjgaCPUBits   atomic.Uint64
	bjgaHostOnce  sync.Once
)

func bjgaStartHostSample() {
	bjgaHostOnce.Do(func() {
		go func() {
			_, _ = cpu.Percent(0, false)
			tk := time.NewTicker(2 * time.Second)
			defer tk.Stop()
			for range tk.C {
				p, err := cpu.Percent(0, false)
				if err != nil || len(p) == 0 {
					continue
				}
				bjgaCPUBits.Store(math.Float64bits(p[0]))
			}
		}()
	})
}

func bjgaRound1(v float64) float64 {
	return math.Round(v*10) / 10
}

// BjgaHostLive 只读正在跑 ism_server 的这一台。没有其它三台的采集探针。
func BjgaHostLive() map[string]interface{} {
	bjgaStartHostSample()
	if bjgaCPUBits.Load() == 0 {
		if p, err := cpu.Percent(200*time.Millisecond, false); err == nil && len(p) > 0 {
			bjgaCPUBits.Store(math.Float64bits(p[0]))
		}
	}
	cpuPct := math.Float64frombits(bjgaCPUBits.Load())
	memPct := 0.0
	if m, err := mem.VirtualMemory(); err == nil && m != nil {
		memPct = m.UsedPercent
	}
	diskPct := 0.0
	if d, err := disk.Usage("/"); err == nil && d != nil {
		diskPct = d.UsedPercent
	}
	rss := uint64(0)
	if proc, err := process.NewProcess(int32(os.Getpid())); err == nil {
		if mi, err := proc.MemoryInfo(); err == nil && mi != nil {
			rss = mi.RSS
		}
	}
	name, _ := os.Hostname()
	up := int(time.Since(bjgaHostStart).Seconds())
	return map[string]interface{}{
		"host":      name,
		"cpu":       bjgaRound1(cpuPct),
		"mem":       bjgaRound1(memPct),
		"disk":      bjgaRound1(diskPct),
		"rssMb":     rss / 1024 / 1024,
		"goroutine": runtime.NumGoroutine(),
		"uptimeSec": up,
		"at":        time.Now().Format("15:04:05"),
		"source":    "本机",
	}
}

func bjgaRecentAudit(n int) []BjgaAuditLog {
	var rows []BjgaAuditLog
	Db.Order("id desc").Limit(n).Find(&rows)
	if rows == nil {
		rows = []BjgaAuditLog{}
	}
	return rows
}
