package watchdog

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"sync/atomic"
	"time"
)

type GoroutineCollector struct {
	BaseCollector
}

type GoroutineWatcher struct {
	BaseWatcher
}

func NewGoroutineWatcher() *GoroutineWatcher {
	w := &GoroutineWatcher{}
	w.MaxThreshold = 3000 // 超过 3000 进行告警
	w.IncrThreshold = 500 // 2s内goroutine上涨超过500
	w.Interval = 2 * time.Second
	w.close = make(chan struct{})

	executor := &GoroutineCollector{}
	executor.CollectFilePath = defaultCollectPath
	executor.CollectSec = 10
	executor.MaxFileBackup = 3
	w.Executors = executor
	return w
}

// 设置goroutine当前值
func (c *GoroutineWatcher) set() {
	v := c.get()
	if c.lastPercent == 0 {
		c.lastPercent = v
	} else {
		c.lastPercent = c.currentPercent
	}
	c.currentPercent = v
}

func (c *GoroutineWatcher) get() float64 {
	i := runtime.NumGoroutine()
	return float64(i)
}

// 触发判断
func (c *GoroutineWatcher) trigger() bool {
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

func (c *GoroutineWatcher) check() error {
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

func (c *GoroutineWatcher) Watch() {
	if err := c.check(); err != nil {
		panic(err)
	}
	go func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Println("cpu watch recover:", e)
				panic(e)
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
					if err := c.Executors.Execute(); err != nil {
						fmt.Println(err.Error())
					}
				}
			case <-c.close:
				t.Stop()
				return
			}
		}
	}()
}

func (c *GoroutineWatcher) Stop() {
	close(c.close)
}

// 采集
func (c *GoroutineCollector) Execute() error {
	if !atomic.CompareAndSwapInt32(&c.flag, 0, 1) {
		return nil
	}
	// 删除一个最老的
	if c.fileIndex >= c.MaxFileBackup {
		atomic.StoreInt32(&c.fileIndex, 0)
	}

	// 删除之前的
	if err := removeFileByPrefix(c.CollectFilePath, fmt.Sprintf("%s-%d", goPrefix, c.fileIndex)); err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s-%d-%s.pprof", goPrefix, c.fileIndex, time.Now().Format("01-02-15:04:05"))
	w, err := os.Create(filepath.Join(c.CollectFilePath, fileName))
	if err != nil {
		return err
	}
	if err := pprof.Lookup(goPrefix).WriteTo(w, 0); err != nil {
		return err
	}

	atomic.AddInt32(&c.fileIndex, 1)
	atomic.StoreInt32(&c.flag, 0)

	return nil
}
