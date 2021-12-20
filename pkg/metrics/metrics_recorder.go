package metrics

type Recorder interface {
	CountRequest()
	CountDatabaseCall()
	CountError()
	ReportRequestDuration(func())
}
