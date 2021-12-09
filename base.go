package watchdog

import "time"

type BaseWatcher struct {
	lastPercent    float64       // 上一次的指标
	currentPercent float64       // 当前指标
	close          chan struct{} // 关闭
	IncrThreshold  float64       // 增长量阈值
	MaxThreshold   float64       // 最大指标阈值
	Interval       time.Duration // 监控间隔
	Executors      []Executors   // 执行者
}

type BaseCollector struct {
	flag            int32  // 标识当前是否正在采集
	fileIndex       int32  // 文件下标
	CollectSec      int    // 收集时长,单位秒
	CollectFilePath string // 收集的文件存储位置
	MaxFileBackup   int32  // 最多存储几个
}
