package consumer

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

type Metrics struct {
	taskCountPerState         *prometheus.GaugeVec
	taskTotalReceived         prometheus.Counter
	taskTotalPerTaskType      *prometheus.CounterVec
	taskTotalValuePerTaskType *prometheus.CounterVec
}

func MustLoad() *Metrics {
	gauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "task_count_per_state",
		Help: "counts tasks in each state",
	}, []string{"task_state"})

	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "task_total_received",
		Help: "total tasks received",
	})

	counterVec := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "task_total_per_task_type",
		Help: "total tasks per task type",
	}, []string{"task_type"})

	counterVec2 := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "task_total_value_per_task_type",
		Help: "total value sum of processed tasks per task type",
	}, []string{"task_type"})

	prometheus.MustRegister(gauge, counter, counterVec, counterVec2)

	return &Metrics{gauge, counter, counterVec, counterVec2}
}

func (m *Metrics) TaskJustReceived() {
	m.taskTotalReceived.Inc()
	m.taskCountPerState.WithLabelValues("received").Inc()
}

func (m *Metrics) TaskIsProcessing() {
	m.taskCountPerState.WithLabelValues("received").Dec()
	m.taskCountPerState.WithLabelValues("processing").Inc()
}

func (m *Metrics) TaskIsDone(taskType int, taskValue int) {
	m.taskCountPerState.WithLabelValues("processing").Dec()
	m.taskCountPerState.WithLabelValues("done").Inc()
	m.taskTotalPerTaskType.WithLabelValues(strconv.Itoa(taskType)).Inc()
	m.taskTotalValuePerTaskType.WithLabelValues(strconv.Itoa(taskType)).Add(float64(taskValue))
}
