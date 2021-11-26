package streamingService

const ShareLinkPatternKey = "shareLinkPattern"

type Config struct {
	ServiceName string `mapstructure:"serviceName"`
	Enabled     bool   `mapstructure:"enabled"`
	Properties  map[string]string
}
