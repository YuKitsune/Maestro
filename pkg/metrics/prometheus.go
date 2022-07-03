package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type prometheusMetricsRecorder struct {
	requestDurationHistogram *prometheus.HistogramVec

	databaseCallCounter prometheus.Counter

	serverErrorCounter prometheus.Counter
	clientErrorCounter prometheus.Counter

	appleMusicRequestCounter prometheus.Counter
	spotifyRequestCounter    prometheus.Counter
	deezerRequestCounter     prometheus.Counter
}

func NewPrometheusMetricsRecorder() (Recorder, error) {

	reqDur := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "maestro_request_duration_seconds",
		Help:    "The total duration of a HTTP request",
		Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
	}, []string{"trace_id", "method", "route"})

	if err := prometheus.Register(reqDur); err != nil {
		return nil, err
	}

	serverErrorCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "maestro_server_error_count",
		Help: "The total number of server-side errors",
	})

	if err := prometheus.Register(serverErrorCounter); err != nil {
		return nil, err
	}

	clientErrorCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "maestro_client_error_count",
		Help: "The total number of client-side errors",
	})

	if err := prometheus.Register(clientErrorCounter); err != nil {
		return nil, err
	}

	dbCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "maestro_db_call_count",
		Help: "The total number of database calls",
	})

	if err := prometheus.Register(dbCounter); err != nil {
		return nil, err
	}

	amCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "maestro_apple_music_request_count",
		Help: "The total number of requests sent to the Apple Music API",
	})

	if err := prometheus.Register(amCounter); err != nil {
		return nil, err
	}

	spCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "maestro_spotify_request_count",
		Help: "The total number of requests sent to the Spotify API",
	})

	if err := prometheus.Register(spCounter); err != nil {
		return nil, err
	}

	dzCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "maestro_deezer_request_count",
		Help: "The total number of requests sent to the Deezer API",
	})

	if err := prometheus.Register(dzCounter); err != nil {
		return nil, err
	}

	rec := &prometheusMetricsRecorder{
		reqDur,
		dbCounter,
		serverErrorCounter,
		clientErrorCounter,
		amCounter,
		spCounter,
		dzCounter,
	}

	return rec, nil
}

func (p prometheusMetricsRecorder) CountDatabaseCall() {
	p.databaseCallCounter.Inc()
}

func (p prometheusMetricsRecorder) CountServerError() {
	p.serverErrorCounter.Inc()
}

func (p prometheusMetricsRecorder) ReportRequestDuration(traceId string, path string, fn func()) {
	begin := time.Now()
	fn()
	dur := time.Since(begin)
	p.requestDurationHistogram.WithLabelValues(traceId, path).Observe(dur.Seconds())
}

func (p prometheusMetricsRecorder) CountAppleMusicRequest() {
	p.appleMusicRequestCounter.Inc()
}

func (p prometheusMetricsRecorder) CountSpotifyRequest() {
	p.spotifyRequestCounter.Inc()
}

func (p prometheusMetricsRecorder) CountDeezerRequest() {
	p.deezerRequestCounter.Inc()
}
