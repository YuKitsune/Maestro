package config

import (
	"github.com/spf13/viper"
	"github.com/yukitsune/maestro/pkg/model"
)

type Spotify interface {
	Type() model.StreamingServiceType
	Name() string
	Enabled() bool
	LogoFileName() string
	ClientId() string
	ClientSecret() string
}

type spotifyViperConfig struct {
	v *viper.Viper
}

func NewSpotifyViperConfig(v *viper.Viper) Spotify {
	v.SetDefault("services.spotify.enabled", true)
	v.SetDefault("services.spotify.logo_file_name", "deezer.png")

	return &spotifyViperConfig{v}
}

func (c *spotifyViperConfig) Type() model.StreamingServiceType {
	return model.SpotifyStreamingService
}

func (c *spotifyViperConfig) Name() string {
	return "Spotify"
}

func (c *spotifyViperConfig) Enabled() bool {
	return c.v.GetBool("services.spotify.enabled")
}

func (c *spotifyViperConfig) LogoFileName() string {
	return c.v.GetString("services.spotify.logo_file_name")
}

func (c *spotifyViperConfig) ClientId() string {
	if !c.v.IsSet("services.spotify.client_id") {
		panic("spotify client_id not set")
	}

	return c.v.GetString("services.spotify.client_id")
}

func (c *spotifyViperConfig) ClientSecret() string {
	if !c.v.IsSet("services.spotify.client_secret") {
		panic("spotify client_secret not set")
	}

	return c.v.GetString("services.spotify.client_secret")
}
