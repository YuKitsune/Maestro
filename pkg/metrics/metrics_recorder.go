package metrics

type Recorder interface {
	CountRequest()
	ReportRequestDuration(func())

	CountServerError()
	CountClientError()

	CountDatabaseCall()

	CountAppleMusicRequest()
	CountSpotifyRequest()
	CountDeezerRequest()
}
