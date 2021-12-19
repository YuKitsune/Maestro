package metrics

type Recorder interface {
	CountRequest()
	CountDatabaseCall()
	ReportRequestDuration(func())
}
