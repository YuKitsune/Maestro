package config

import (
	"github.com/spf13/viper"
	"github.com/yukitsune/maestro/pkg/model"
)

type Service interface {
	Type() model.StreamingServiceType
	Name() string
	LogoFileName() string
	Enabled() bool
}

type Services interface {
	AppleMusic() AppleMusic
	Spotify() Spotify
	Deezer() Deezer
	AsMap() map[model.StreamingServiceType]Service
}

type servicesViperConfig struct {
	appleMusic AppleMusic
	spotify    Spotify
	deezer     Deezer
}

func NewServicesViperConfig(v *viper.Viper) Services {
	return &servicesViperConfig{
		appleMusic: NewAppleMusicViperConfig(v.Sub("apple_music")),
		spotify:    NewSpotifyViperConfig(v.Sub("spotify")),
		deezer:     NewDeezerViperConfig(v.Sub("deezer")),
	}
}

func (s *servicesViperConfig) AppleMusic() AppleMusic {
	return s.appleMusic
}

func (s *servicesViperConfig) Spotify() Spotify {
	return s.spotify
}

func (s *servicesViperConfig) Deezer() Deezer {
	return s.deezer
}

func (s *servicesViperConfig) AsMap() map[model.StreamingServiceType]Service {
	configs := make(map[model.StreamingServiceType]Service)
	configs[model.AppleMusicStreamingService] = s.AppleMusic()
	configs[model.SpotifyStreamingService] = s.Spotify()
	configs[model.DeezerStreamingService] = s.Deezer()
	return configs
}
