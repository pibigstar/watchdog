package watchdog

import (
	wlog "log"
	"os"
	"strconv"
)

var (
	defaultCollectPath     = "/var/tmp/pprof"
	defaultWatchDogLogFile = "watchdog.log"
	defaultPprofPort       = 9999

	log  = wlog.New(os.Stdout, "", wlog.Ldate|wlog.Ltime)
	list []Watcher
)

// env
const (
	collectPathEnv = "COLLECT_PATH"
	pprofPort      = "PPROF_PORT"
	logFile        = "WATCH_DOG_LOG_FILE"
)

// executor file prefix
const (
	cpuPrefix    = "cpu"
	goPrefix     = "goroutine"
	memoryPrefix = "memory"
)

func init() {
	if s := os.Getenv(collectPathEnv); s != "" {
		defaultCollectPath = s
	}
	if s := os.Getenv(pprofPort); s != "" {
		if port, err := strconv.Atoi(s); err == nil {
			defaultPprofPort = port
		}
	}
	if s := os.Getenv(logFile); s != "" {
		defaultWatchDogLogFile = s
	}

	if f, err := os.Create(defaultWatchDogLogFile); err == nil {
		log = wlog.New(f, "", wlog.Ldate|wlog.Ltime)
	}
}

func Run(watchers ...Watcher) {
	// 检查目录是否存在与可写权限
	if err := checkPath(defaultCollectPath); err != nil {
		log.Println("checkPath", err.Error())
		return
	}
	list = append(list, watchers...)
	for _, l := range list {
		l.Watch()
	}
	// 启动http文件server
	go runFileServer()
}

func SetCollectFilePath(path string) {
	defaultCollectPath = path
}
