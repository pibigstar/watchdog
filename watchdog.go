package watchdog

import "os"

var (
	defaultCollectPath = "/var/tmp/pprof"
	list               []Watcher
)

const (
	CollectPathEnv = "COLLECT_PATH"
	cpuPrefix    = "cpu"
	goPrefix     = "goroutine"
	memoryPrefix = "memory"
)

func init() {
	if s := os.Getenv(CollectPathEnv); s != "" {
		defaultCollectPath = s
	}
}

func Run(watchers ...Watcher) {
	list = append(list, watchers...)
	for _, l := range list {
		l.Watch()
	}
	// 检查目录是否存在与可写权限
	if err := checkPath(defaultCollectPath); err != nil {
		panic(err)
	}

	// 启动http文件server
	go runFileServer(defaultCollectPath)
}

func SetCollectFilePath(path string) {
	defaultCollectPath = path
}

func main() {
	watchers := []Watcher{
		NewMemoryWatcher(),
		NewCpuWatcher(),
		NewGoroutineWatcher(),
	}
	Run(watchers...)

	select {}
}
