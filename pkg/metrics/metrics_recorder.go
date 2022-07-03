package metrics

type Recorder interface {
	ReportRequestDuration(path string, fn func())

	CountServerError()

	CountDatabaseCall()

	CountAppleMusicRequest()
	CountSpotifyRequest()
	CountDeezerRequest()
}
