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
	cfg      streamingservice.Config
	svcFuncs map[model.StreamingServiceKey]func(streamingservice.ServiceConfig) (streamingservice.StreamingService, error)
}

func NewDefaultProvider(cfg streamingservice.Config, rec metrics.Recorder) (streamingservice.ServiceProvider, error) {

	svcFuncs := make(map[model.StreamingServiceKey]func(streamingservice.ServiceConfig) (streamingservice.StreamingService, error))
	for key, config := range cfg {
		if !config.Enabled() {
			continue
		}

		switch key {
		case applemusic.Key:
			fn := func(cfg streamingservice.ServiceConfig) (streamingservice.StreamingService, error) {
				appleCfg := cfg.(*applemusic.Config)
				svc := applemusic.NewAppleMusicStreamingService(appleCfg, rec)
				return svc, nil
			}

			svcFuncs[key] = fn
			break

		case deezer.Key:
			fn := func(_ streamingservice.ServiceConfig) (streamingservice.StreamingService, error) {
				svc := deezer.NewDeezerStreamingService(rec)
				return svc, nil
			}

			svcFuncs[key] = fn
			break

		case spotify.Key:
			fn := func(cfg streamingservice.ServiceConfig) (streamingservice.StreamingService, error) {
				spotifyCfg := cfg.(*spotify.Config)
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
		cfg:      cfg,
		svcFuncs: svcFuncs,
	}, nil
}

func (p *defaultServiceProvider) GetService(key model.StreamingServiceKey) (streamingservice.StreamingService, error) {
	cfg, err := p.GetConfig(key)
	if err != nil {
		return nil, err
	}

	if !cfg.Enabled() {
		return nil, fmt.Errorf("service %s is disabled", key)
	}

	svcFunc, ok := p.svcFuncs[key]
	if !ok {
		return nil, fmt.Errorf("couldn't find service with key %s", key)
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

func (p *defaultServiceProvider) GetConfig(key model.StreamingServiceKey) (streamingservice.ServiceConfig, error) {
	cfg, ok := p.cfg[key]
	if !ok {
		return nil, fmt.Errorf("couldn't find service config with key %s", key)
	}

	return cfg, nil
}

func (p *defaultServiceProvider) ListConfigs() streamingservice.Config {
	return p.cfg
}
