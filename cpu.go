package watchdog

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/pprof"
	"sync/atomic"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

// CPU采集器
type CpuCollector struct {
	BaseCollector
}

// CPU监测者
type CpuWatcher struct {
	BaseWatcher
}

func NewCpuWatcher() *CpuWatcher {
	w := &CpuWatcher{}
	w.MaxThreshold = 80
	w.IncrThreshold = 30
	w.Interval = 5 * time.Second
	w.close = make(chan struct{})

	executor := &CpuCollector{}
	executor.CollectSec = 10
	executor.MaxFileBackup = 3
	w.Executors = append(w.Executors, executor)
	return w
}

// 设置cpu当前值
func (c *CpuWatcher) set() {
	v := c.get()
	if c.lastPercent == 0 {
		c.lastPercent = v
	} else {
		c.lastPercent = c.currentPercent
	}
	c.currentPercent = v
}

// 获取cpu使用率
func (c *CpuWatcher) get() float64 {
	cpus, _ := cpu.Percent(c.Interval, false)
	return cpus[0]
}

// 触发判断
func (c *CpuWatcher) trigger() bool {
	// 指标在下降
	if c.currentPercent <= c.lastPercent {
		return false
	}
	// 当前指标超过最大阈值
	if c.currentPercent >= c.MaxThreshold {
		return true
	}
	// 上涨量超过上涨阈值
	if (c.currentPercent - c.lastPercent) >= c.IncrThreshold {
		return true
	}
	return false
}

func (c *CpuCollector) Execute() error {
	if !atomic.CompareAndSwapInt32(&c.flag, 0, 1) {
		return nil
	}

	if c.fileIndex >= c.MaxFileBackup {
		atomic.StoreInt32(&c.fileIndex, 0)
	}

	// 环形删除
	if err := removeFileByPrefix(defaultCollectPath, fmt.Sprintf("%s-%d", cpuPrefix, c.fileIndex)); err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s-%d-%s.pprof", cpuPrefix, c.fileIndex, time.Now().Format("01-02-15:04:05"))
	w, err := os.Create(filepath.Join(defaultCollectPath, fileName))
	if err != nil {
		return err
	}
	if err := pprof.StartCPUProfile(w); err != nil {
		return fmt.Errorf("could not enable CPU profiling: %w", err)
	}
	time.Sleep(time.Duration(c.CollectSec) * time.Second)
	pprof.StopCPUProfile()

	atomic.AddInt32(&c.fileIndex, 1)
	atomic.StoreInt32(&c.flag, 0)

	return nil
}

func (c *CpuWatcher) check() error {
	if c.Interval <= 0 {
		return fmt.Errorf("interval less 0")
	}
	if c.MaxThreshold <= 0 {
		return fmt.Errorf("MaxThreshold less 0")
	}
	if c.IncrThreshold <= 0 {
		return fmt.Errorf("IncrThreshold less 0")
	}
	return nil
}

func (c *CpuWatcher) Watch() {
	if err := c.check(); err != nil {
		log.Println("check cpu", err.Error())
		return
	}
	go func() {
		defer func() {
			if e := recover(); e != nil {
				log.Println("cpu watch recover:", e)
			}
		}()

		t := time.NewTicker(c.Interval)
		for {
			select {
			case <-t.C:
				// 设置当前cpu信息
				c.set()
				// 判断是否超过阈值
				if c.trigger() {
					// 采集
					for _, exec := range c.Executors {
						if err := exec.Execute(); err != nil {
							log.Println(err.Error())
						}
					}
				}
			case <-c.close:
				t.Stop()
				return
			}
		}
	}()
}

func (c *CpuWatcher) Stop() {
	close(c.close)
}
