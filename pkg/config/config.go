package config

import (
	"github.com/yukitsune/maestro/pkg/model"
)

type Config struct {
	API      *API      `mapstructure:"api"`
	Log      *Log      `mapstructure:"logging"`
	Database *Database `mapstructure:"database"`
	Services *Services `mapstructure:"services"`
}

type API struct {
	Port            int    `mapstructure:"port"`
	AssetsDirectory string `mapstructure:"assets_dir"`
}

type Log struct {
	Level string `mapstructure:"level"`
	Loki  *struct {
		Host   string            `mapstructure:"host"`
		Labels map[string]string `mapstructure:"labels"`
	} `mapstructure:"loki"`
}

type Database struct {
	URI      string `mapstructure:"uri"`
	Database string `mapstructure:"name"`
}

type Services struct {
	AppleMusic *AppleMusic `mapstructure:"apple_music"`
	Spotify    *Spotify    `mapstructure:"spotify"`
	Deezer     *Deezer     `mapstructure:"deezer"`
}

func (s *Services) AsMap() map[model.StreamingServiceType]Service {
	configs := make(map[model.StreamingServiceType]Service)
	configs[model.AppleMusicStreamingService] = s.AppleMusic
	configs[model.SpotifyStreamingService] = s.Spotify
	configs[model.DeezerStreamingService] = s.Deezer
	return configs
}

type AppleMusic struct {
	ServiceLogoFileName string `mapstructure:"logo_file_name"`
	IsEnabled           bool   `mapstructure:"enabled"`
	Token               string `mapstructure:"token"`
}

func (c *AppleMusic) Type() model.StreamingServiceType {
	return model.AppleMusicStreamingService
}

func (c *AppleMusic) Name() string {
	return "Apple Music"
}

func (c *AppleMusic) LogoFileName() string {
	return c.ServiceLogoFileName
}

func (c *AppleMusic) Enabled() bool {
	return c.IsEnabled
}

type Spotify struct {
	ServiceLogoFileName string `mapstructure:"logo_file_name"`
	IsEnabled           bool   `mapstructure:"enabled"`
	ClientID            string `mapstructure:"client_id"`
	ClientSecret        string `mapstructure:"client_secret"`
}

func (c *Spotify) Type() model.StreamingServiceType {
	return model.SpotifyStreamingService
}

func (c *Spotify) Name() string {
	return "Spotify"
}

func (c *Spotify) LogoFileName() string {
	return c.ServiceLogoFileName
}

func (c *Spotify) Enabled() bool {
	return c.IsEnabled
}

type Deezer struct {
	ServiceLogoFileName string `mapstructure:"logo_file_name"`
	IsEnabled           bool   `mapstructure:"enabled"`
}

func (c *Deezer) Type() model.StreamingServiceType {
	return model.DeezerStreamingService
}

func (c *Deezer) Name() string {
	return "Deezer"
}

func (c *Deezer) LogoFileName() string {
	return c.ServiceLogoFileName
}

func (c *Deezer) Enabled() bool {
	return c.IsEnabled
}
