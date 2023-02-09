package model

type StreamingServiceType string

const (
	AppleMusicStreamingService StreamingServiceType = "apple_music"
	SpotifyStreamingService    StreamingServiceType = "spotify"
	DeezerStreamingService     StreamingServiceType = "deezer"
)

func (s StreamingServiceType) String() string {
	return string(s)
}

type StreamingService struct {
	Key     StreamingServiceType
	Name    string
	Enabled bool
}
