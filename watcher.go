package watchdog

type Watcher interface {
	Watch()
	Stop()
}
