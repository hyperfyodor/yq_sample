package consumer

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

type Metrics struct {
	TaskCountPerState         *prometheus.GaugeVec
	TaskTotalReceived         prometheus.Counter
	TaskTotalPerTaskType      *prometheus.CounterVec
	TaskTotalValuePerTaskType *prometheus.CounterVec
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
	m.TaskTotalReceived.Inc()
	m.TaskCountPerState.WithLabelValues("received").Inc()
}

func (m *Metrics) TaskIsProcessing() {
	m.TaskCountPerState.WithLabelValues("received").Dec()
	m.TaskCountPerState.WithLabelValues("processing").Inc()
}

func (m *Metrics) TaskIsDone(taskType int, taskValue int) {
	m.TaskCountPerState.WithLabelValues("processing").Dec()
	m.TaskCountPerState.WithLabelValues("done").Inc()
	m.TaskTotalPerTaskType.WithLabelValues(strconv.Itoa(taskType)).Inc()
	m.TaskTotalValuePerTaskType.WithLabelValues(strconv.Itoa(taskType)).Add(float64(taskValue))
}

func (m *Metrics) Unregister() {
	prometheus.Unregister(m.TaskTotalPerTaskType)
	prometheus.Unregister(m.TaskTotalValuePerTaskType)
	prometheus.Unregister(m.TaskCountPerState)
	prometheus.Unregister(m.TaskTotalReceived)
}
