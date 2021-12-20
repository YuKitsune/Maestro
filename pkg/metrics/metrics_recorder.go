package metrics

type Recorder interface {
	CountRequest()
	CountDatabaseCall()
	CountServerError()
	CountClientError()
	ReportRequestDuration(func())
}
