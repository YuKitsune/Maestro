package provider

import (
	"fmt"
	"github.com/yukitsune/maestro/pkg/metrics"
	"github.com/yukitsune/maestro/pkg/model"
	"github.com/yukitsune/maestro/pkg/streamingservice"
	"github.com/yukitsune/maestro/pkg/streamingservice/applemusic"
	"github.com/yukitsune/maestro/pkg/streamingservice/deezer"
	"github.com/yukitsune/maestro/pkg/streamingservice/spotify"
)

type defaultServiceProvider struct {
	cfg  streamingservice.Config
	svcs []streamingservice.StreamingService
}

func NewDefaultProvider(cfg streamingservice.Config, rec metrics.Recorder) (streamingservice.ServiceProvider, error) {

	var services []streamingservice.StreamingService
	for key, config := range cfg {
		if !config.Enabled() {
			continue
		}

		switch key {
		case applemusic.Key:
			cfg := config.(*applemusic.Config)
			s := applemusic.NewAppleMusicStreamingService(cfg, rec)
			services = append(services, s)
			break

		case deezer.Key:
			s := deezer.NewDeezerStreamingService(rec)
			services = append(services, s)
			break

		case spotify.Key:
			cfg := config.(*spotify.Config)
			s, err := spotify.NewSpotifyStreamingService(cfg, rec)
			if err != nil {
				return nil, fmt.Errorf("failed to initialize spotify streaming service: %s", err.Error())
			}

			services = append(services, s)
			break
		}
	}

	return &defaultServiceProvider{
		cfg:  cfg,
		svcs: services,
	}, nil
}

func (p *defaultServiceProvider) GetService(key model.StreamingServiceKey) streamingservice.StreamingService {
	for _, svc := range p.svcs {
		if svc.Key() == key {
			return svc
		}
	}

	return nil
}

func (p *defaultServiceProvider) ListServices() []streamingservice.StreamingService {
	return p.svcs
}

func (p *defaultServiceProvider) GetConfig(key model.StreamingServiceKey) streamingservice.ServiceConfig {
	for k, cfg := range p.cfg {
		if k == key {
			return cfg
		}
	}

	return nil
}

func (p *defaultServiceProvider) ListConfigs() streamingservice.Config {
	return p.cfg
}
