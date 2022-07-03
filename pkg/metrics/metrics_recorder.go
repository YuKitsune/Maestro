package metrics

type Recorder interface {
	ReportRequestDuration(traceId string, path string, fn func())

	CountServerError()

	CountDatabaseCall()

	CountAppleMusicRequest()
	CountSpotifyRequest()
	CountDeezerRequest()
}
