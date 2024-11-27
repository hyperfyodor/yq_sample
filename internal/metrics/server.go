package metrics

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func Listen(port string) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(":"+port, mux)
}
