package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type prometheusMetricsRecorder struct {
	requestCounter           prometheus.Counter
	databaseCallCounter      prometheus.Counter
	requestDurationHistogram prometheus.Histogram
	linkSourceSummary        prometheus.Summary
}

func NewPrometheusMetricsRecorder() (Recorder, error) {

	reqCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "maestro_request_count",
		Help: "The total number of HTTP requests served",
	})

	if err := prometheus.Register(reqCounter); err != nil {
		return nil, err
	}

	dbCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "maestro_db_call_count",
		Help: "The total number of database calls",
	})

	if err := prometheus.Register(dbCounter); err != nil {
		return nil, err
	}

	reqDur := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "maestro_request_duration",
		Help: "The total duration of a HTTP request",
	})

	if err := prometheus.Register(reqDur); err != nil {
		return nil, err
	}

	// Todo: Better description
	linkSummary := prometheus.NewSummary(prometheus.SummaryOpts{
		Name: "maestro_link_sources",
		Help: "The number of links per streaming service used",
	})

	if err := prometheus.Register(linkSummary); err != nil {
		return nil, err
	}

	rec := &prometheusMetricsRecorder{
		reqCounter,
		dbCounter,
		reqDur,
		linkSummary,
	}

	return rec, nil
}

func (p prometheusMetricsRecorder) CountRequest() {
	p.requestCounter.Inc()
}

func (p prometheusMetricsRecorder) CountDatabaseCall() {
	p.databaseCallCounter.Inc()
}

func (p prometheusMetricsRecorder) ReportRequestDuration(fn func()) {
	timer := prometheus.NewTimer(p.requestDurationHistogram)
	fn()
	timer.ObserveDuration()
}
