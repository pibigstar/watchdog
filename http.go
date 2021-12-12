package watchdog

import (
	"fmt"
	"net/http"
	"os"

	auth "github.com/abbot/go-http-auth"
)

const (
	pprofPasswordEnv = "PPROF_PASSWORD"
	defaultUser      = "sns"
	defaultPassword  = "$2a$10$h/B69yeHyt9BwjeGO7CJSu7YMiN7VpjAw/8UcLNkLlcNXhHJwpOw."
)

func Secret(user, realm string) string {
	password := defaultPassword
	if user == defaultUser {
		if s := os.Getenv(pprofPasswordEnv); s != "" {
			password = s
		}
		return password
	}
	return ""
}

func runFileServer() {
	authenticator := auth.NewBasicAuthenticator("", Secret)
	handFunc := auth.JustCheck(authenticator, handleFileServer(defaultCollectPath, "/"))
	err := http.ListenAndServe(fmt.Sprintf(":%d", defaultPprofPort), handFunc)
	if err != nil {
		log.Println("ListenAndServe", err.Error())
	}
}

func handleFileServer(dir, prefix string) http.HandlerFunc {
	fs := http.FileServer(http.Dir(dir))
	realHandler := http.StripPrefix(prefix, fs).ServeHTTP
	return func(w http.ResponseWriter, req *http.Request) {
		log.Println(req.URL)
		realHandler(w, req)
	}
}
