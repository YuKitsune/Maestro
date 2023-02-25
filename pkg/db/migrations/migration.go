package migrations

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type Migration interface {
	Execute(context.Context, *mongo.Database) error
	Version() int
}
