package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	counters *prometheus.CounterVec
)

const (
	label    = "action"
	create   = "create"
	update   = "update"
	remove   = "remove"
)

func Register() {
	counters = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ocp_video_api",
			Help: "Video API metrics",
		},
		[]string{label},
	)

	prometheus.MustRegister(counters)
}

func increment(action string, times uint64) {
	counters.With(prometheus.Labels{label: action}).Add(float64(times))
}


func IncrementSuccessfulCreates(times uint64) {
	increment(create, times)
}

func IncrementSuccessfulUpdates(times uint64) {
	increment(update, times)
}

func IncrementSuccessfulRemoves(times uint64) {
	increment(remove, times)
}
