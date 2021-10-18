package watchdog

import "net/http"

func runFileServer(path string) {
	err := http.ListenAndServe(":9999", http.FileServer(http.Dir(path)))
	if err != nil {
		panic(err)
	}
}
