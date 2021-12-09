package watchdog

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/pprof"
	"sync/atomic"
	"time"

	"github.com/shirou/gopsutil/mem"
)

type MemoryCollector struct {
	BaseCollector
}

type MemoryWatcher struct {
	BaseWatcher
}

func NewMemoryWatcher() *MemoryWatcher {
	w := &MemoryWatcher{}
	w.MaxThreshold = 80
	w.IncrThreshold = 30
	w.Interval = 5 * time.Second
	w.close = make(chan struct{})

	executor := &MemoryCollector{}
	executor.CollectFilePath = defaultCollectPath
	executor.CollectSec = 10
	executor.MaxFileBackup = 3
	w.Executors = append(w.Executors, executor)
	return w
}

// 设置memory当前值
func (c *MemoryWatcher) set() {
	v := c.get()
	if c.lastPercent == 0 {
		c.lastPercent = v
	} else {
		c.lastPercent = c.currentPercent
	}
	c.currentPercent = v
}

func (c *MemoryWatcher) get() float64 {
	memInfo, _ := mem.VirtualMemory()
	return memInfo.UsedPercent
}

// 触发判断
func (c *MemoryWatcher) trigger() bool {
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

func (c *MemoryWatcher) check() error {
	if c.Interval <= 0 {
		return fmt.Errorf("interval less 0")
	}
	if c.MaxThreshold <= 0 {
		return fmt.Errorf("maxThreshold less 0")
	}
	if c.IncrThreshold <= 0 {
		return fmt.Errorf("incrThreshold less 0")
	}
	return nil
}

func (c *MemoryWatcher) Watch() {
	if err := c.check(); err != nil {
		log.Println("check memory", err.Error())
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
				// 设置当前信息
				c.set()
				// 判断是否达到阈值
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

func (c *MemoryWatcher) Stop() {
	close(c.close)
}

// Execute 采集内存
func (c *MemoryCollector) Execute() error {
	if !atomic.CompareAndSwapInt32(&c.flag, 0, 1) {
		return nil
	}

	if c.fileIndex >= c.MaxFileBackup {
		atomic.StoreInt32(&c.fileIndex, 0)
	}

	// 环形删除
	if err := removeFileByPrefix(c.CollectFilePath, fmt.Sprintf("%s-%d", memoryPrefix, c.fileIndex)); err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s-%d-%s.pprof", memoryPrefix, c.fileIndex, time.Now().Format("01-02-15:04:05"))
	w, err := os.Create(filepath.Join(c.CollectFilePath, fileName))
	if err != nil {
		return err
	}
	if err := pprof.Lookup("heap").WriteTo(w, 0); err != nil {
		return err
	}

	atomic.AddInt32(&c.fileIndex, 1)
	atomic.StoreInt32(&c.flag, 0)

	return nil
}
