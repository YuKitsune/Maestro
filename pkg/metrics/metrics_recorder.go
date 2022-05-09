package metrics

type Recorder interface {
	CountRequest()
	ReportRequestDuration(func())

	CountServerError()

	CountDatabaseCall()

	CountAppleMusicRequest()
	CountSpotifyRequest()
	CountDeezerRequest()
}
