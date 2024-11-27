package profiling

import (
	"net/http"
	_ "net/http/pprof"
)

func Listen(port string) error {
	mux := http.NewServeMux()
	return http.ListenAndServe(":"+port, mux)
}
