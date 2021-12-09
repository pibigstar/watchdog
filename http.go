package watchdog

import (
	"fmt"
	"net/http"
)

func runFileServer() {
	err := http.ListenAndServe(fmt.Sprintf(":%d", defaultPprofPort), http.FileServer(http.Dir(defaultCollectPath)))
	if err != nil {
		log.Println("ListenAndServe", err.Error())
	}
}
