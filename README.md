# 看门狗> 当服务器波动时，自动采集pprof- 默认提供 CPU、内存、goroutine 波动 监控    - 检查间隔可配置    - 环比上涨可配置    - 可配置最大值，当到达该值时也需要自动采集- 采集pprof    - 采集文件数上限可配置，不能一直采集    - 采集目录可配置    - 提供http 接口查看采集的所有文件    - 采集的文件需要支持下载- 扩展    - 用户可自定义自己的监控    - 用户可自定义触发时间（不能仅仅是采集）## Quick Start```go// 监控cpuwatchdog.Run(watchdog.NewCpuWatcher())```