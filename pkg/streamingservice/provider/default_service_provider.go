package provider

import (
	"fmt"

	"github.com/yukitsune/maestro/pkg/config"
	"github.com/yukitsune/maestro/pkg/metrics"
	"github.com/yukitsune/maestro/pkg/model"
	"github.com/yukitsune/maestro/pkg/streamingservice"
	"github.com/yukitsune/maestro/pkg/streamingservice/applemusic"
	"github.com/yukitsune/maestro/pkg/streamingservice/deezer"
	"github.com/yukitsune/maestro/pkg/streamingservice/spotify"
)

type defaultServiceProvider struct {
	cfgMap   map[model.StreamingServiceType]config.Service
	svcFuncs map[model.StreamingServiceType]func(config.Service) (streamingservice.StreamingService, error)
}

func NewDefaultProvider(cfg config.Services, rec metrics.Recorder) (streamingservice.ServiceProvider, error) {

	cfgMap := cfg.AsMap()
	svcFuncs := make(map[model.StreamingServiceType]func(config.Service) (streamingservice.StreamingService, error))

	for key, _ := range cfgMap {
		switch key {
		case model.AppleMusicStreamingService:
			fn := func(cfg config.Service) (streamingservice.StreamingService, error) {
				appleCfg := cfg.(config.AppleMusic)
				svc := applemusic.NewAppleMusicStreamingService(appleCfg, rec)
				return svc, nil
			}

			svcFuncs[key] = fn
			break

		case model.DeezerStreamingService:
			fn := func(cfg config.Service) (streamingservice.StreamingService, error) {
				deezerCfg := cfg.(config.Deezer)
				svc := deezer.NewDeezerStreamingService(deezerCfg, rec)
				return svc, nil
			}

			svcFuncs[key] = fn
			break

		case model.SpotifyStreamingService:
			fn := func(cfg config.Service) (streamingservice.StreamingService, error) {
				spotifyCfg := cfg.(config.Spotify)
				s, err := spotify.NewSpotifyStreamingService(spotifyCfg, rec)
				if err != nil {
					return nil, fmt.Errorf("failed to initialize spotify streaming service: %s", err.Error())
				}

				return s, nil
			}

			svcFuncs[key] = fn
			break
		}
	}

	return &defaultServiceProvider{
		cfgMap,
		svcFuncs,
	}, nil
}

func (p *defaultServiceProvider) GetService(serviceType model.StreamingServiceType) (streamingservice.StreamingService, error) {
	cfg, err := p.GetConfig(serviceType)
	if err != nil {
		return nil, err
	}

	svcFunc, ok := p.svcFuncs[serviceType]
	if !ok {
		return nil, fmt.Errorf("couldn't find service type %s", serviceType)
	}

	return svcFunc(cfg)
}

func (p *defaultServiceProvider) ListServices() (streamingservice.StreamingServices, error) {
	svcs := make(streamingservice.StreamingServices)
	for key, fn := range p.svcFuncs {
		cfg, err := p.GetConfig(key)
		if err != nil {
			return nil, err
		}

		if !cfg.Enabled() {
			return nil, fmt.Errorf("service %s is disabled", key)
		}

		svc, err := fn(cfg)
		if err != nil {
			return nil, err
		}

		svcs[key] = svc
	}

	return svcs, nil
}

func (p *defaultServiceProvider) GetConfig(key model.StreamingServiceType) (config.Service, error) {
	cfg, ok := p.cfgMap[key]
	if !ok {
		return nil, fmt.Errorf("couldn't find service config with key %s", key)
	}

	return cfg, nil
}

func (p *defaultServiceProvider) ListConfigs() map[model.StreamingServiceType]config.Service {
	return p.cfgMap
}
