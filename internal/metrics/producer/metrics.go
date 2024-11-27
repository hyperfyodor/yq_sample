package producer

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	TotalProduced prometheus.Counter
}

func MustLoad() *Metrics {
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "total_produced",
		Help: "Total number of produced tasks",
	})

	prometheus.MustRegister(counter)

	return &Metrics{counter}
}

func (p *Metrics) TotalProducedInc() {
	p.TotalProduced.Inc()
}

func (p *Metrics) Unregister() {
	prometheus.Unregister(p.TotalProduced)
}
