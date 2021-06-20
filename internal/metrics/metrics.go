package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

const (
	label  = "action"
	create = "create"
	update = "update"
	remove = "remove"
)

type Metrics interface {
	Init()
	IncrementSuccessfulCreates(uint64)
	IncrementSuccessfulUpdates(uint64)
	IncrementSuccessfulRemoves(uint64)
}

func New() Metrics {
	return &metrics{}
}

type metrics struct {
	counters *prometheus.CounterVec
}

func (m *metrics) Init() {
	m.counters = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ocp_video_api",
			Help: "Video API metrics",
		},
		[]string{label},
	)

	prometheus.MustRegister(m.counters)

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe("0.0.0.0:9091", nil)
		if err != nil {
			panic(err)
		}
	}()
}

func (m *metrics) increment(action string, times uint64) {
	m.counters.With(prometheus.Labels{label: action}).Add(float64(times))
}

func (m *metrics) IncrementSuccessfulCreates(times uint64) {
	m.increment(create, times)
}

func (m *metrics) IncrementSuccessfulUpdates(times uint64) {
	m.increment(update, times)
}

func (m *metrics) IncrementSuccessfulRemoves(times uint64) {
	m.increment(remove, times)
}
