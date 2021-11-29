package db

import (
	"fmt"
	"github.com/yukitsune/camogo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/url"
)

type DatabaseModule struct {
	Config *Config
}

func (m *DatabaseModule) Register(cb camogo.ContainerBuilder) error {

	// Config
	if err := cb.RegisterInstance(m.Config); err != nil {
		return err
	}

	// Database
	err := cb.RegisterFactory(func(cfg *Config) (*mongo.Client, error) {

		uri := fmt.Sprintf(
			"mongodb://%s:%d/%s",
			url.QueryEscape(cfg.Host),
			cfg.Port,
			url.QueryEscape(cfg.Database))

		creds := options.Credential{
			Username: cfg.User,
			Password: cfg.Password,
		}
		opts := options.Client().ApplyURI(uri).SetAuth(creds)

		client, err := mongo.NewClient(opts)
		if err != nil {
			return nil, err
		}

		err = client.Connect(ctx)
		if err != nil {
			return nil, err
		}

		return client, nil
	},
		camogo.TransientLifetime)
	if err != nil {
		return err
	}

	return nil
}
