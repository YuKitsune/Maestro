package migrations_test

import (
	"context"
	"fmt"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"testing"
	"time"
)

// Note: These tests might be somewhat flakey

var client *mongo.Client

func TestMain(m *testing.M) {

	host := "localhost"
	port := "27017"
	username := "root"
	password := "root"
	replicaSet := "rs0"

	// Set up the mongo container
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pull mongodb docker image
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "bitnami/mongodb",
		Tag:        "5.0.8",
		Env: []string{
			// Username and password
			fmt.Sprintf("MONGODB_ROOT_USER=%s", username),
			fmt.Sprintf("MONGODB_ROOT_PASSWORD=%s", password),

			// Replica Set setup
			// Need to use localhost so the client doesn't try to resolve the containers host name
			fmt.Sprintf("MONGODB_ADVERTISED_HOSTNAME=%s", host),
			fmt.Sprintf("MONGODB_ADVERTISED_PORT_NUMBER=%s", port),
			fmt.Sprintf("MONGODB_REPLICA_SET_NAME=%s", replicaSet),
			"MONGODB_REPLICA_SET_MODE=primary",
			"MONGODB_REPLICA_SET_KEY=rs0key",
		},

		// Expose the MongoDB port
		PortBindings: map[docker.Port][]docker.PortBinding{
			docker.Port(fmt.Sprintf("%s/tcp", port)): {{HostIP: "", HostPort: port}},
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// Kill the container after a while
	err = resource.Expire(uint(5 * time.Minute))
	if err != nil {
		log.Fatalf("Could not set resource expiry: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	err = pool.Retry(func() error {

		uri := fmt.Sprintf(
			"mongodb://%s:%s@%s:%s/?replicaSet=%s",
			username,
			password,
			host,
			port,
			replicaSet)
		opts := options.Client().ApplyURI(uri)
		client, err = mongo.NewClient(opts)
		if err != nil {
			return nil
		}

		err = client.Connect(context.TODO())
		if err != nil {
			return err
		}

		return client.Ping(context.TODO(), nil)
	})

	code := 0
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	} else {
		// run tests
		code = m.Run()
	}

	// When we're done, kill and remove the container
	if err = pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	// disconnect mongodb client
	if err = client.Disconnect(context.TODO()); err != nil {
		log.Fatalf("Could not disconnect: %s", err)
	}

	os.Exit(code)
}

func withTestDb(t *testing.T, fn func(db *mongo.Database)) {

	// Get/Create a database with the same name as the test
	db := client.Database(t.Name())

	// Drop the database once the test has completed
	defer func() {
		err := db.Drop(context.Background())
		if err != nil {
			log.Fatalf("Failed to drop database %s: %s", db.Name(), err)
		}
	}()

	fn(db)
}
