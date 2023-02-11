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
		// Todo: Update this to use sub once viper bug is fixed
		appleMusic: NewAppleMusicViperConfig(v),
		spotify:    NewSpotifyViperConfig(v),
		deezer:     NewDeezerViperConfig(v),
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
