package db

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/yukitsune/camogo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
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
	err := cb.RegisterFactory(func(ctx context.Context, cfg *Config) (*mongo.Client, error) {

		creds := options.Credential{
			Username:   cfg.User,
			Password:   cfg.Password,
			AuthSource: cfg.Database,
		}

		opts := options.Client().ApplyURI(cfg.Uri).SetAuth(creds)
		if cfg.CACertFile != nil && len(*cfg.CACertFile) > 0 {

			pemData, err := ioutil.ReadFile(*cfg.CACertFile)
			if err != nil {
				return nil, err
			}

			pool := x509.NewCertPool()
			pool.AppendCertsFromPEM(pemData)

			tlsConfig := &tls.Config{
				ClientCAs: pool,
			}

			opts.SetTLSConfig(tlsConfig)
		}

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

	err = cb.RegisterFactory(func(ctx context.Context, cfg *Config, c *mongo.Client) (*mongo.Database, error) {
		db := c.Database(cfg.Database)
		return db, nil
	},
		camogo.TransientLifetime)
	if err != nil {
		return err
	}

	// Repository
	err = cb.RegisterFactory(NewMongoRepository, camogo.TransientLifetime)
	if err != nil {
		return err
	}

	return nil
}
