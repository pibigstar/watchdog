package watchdog

// 执行者，当监控指标到达指定的值后触发的操作
type Executors interface {
	Execute() error
}
