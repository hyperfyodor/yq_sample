package profiling

import (
	"net/http"
	_ "net/http/pprof"
)

func Listen(port string) error {
	return http.ListenAndServe(":"+port, nil)
}
