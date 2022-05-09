package migrations

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionName = "migrations"
const versionKey = "version"

type Migrator struct {
}

func (m *Migrator) Execute(ctx context.Context, provider MigrationProvider, db *mongo.Database) error {

	// Run all migrations in a session (transaction)
	session, err := db.Client().StartSession()
	if err != nil {
		return err
	}

	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {

		migrations := provider.Migrations()
		for _, migration := range migrations {

			executed, err := m.hasExecuted(sessCtx, db, migration)
			if err != nil {
				return nil, err
			}

			if executed {
				continue
			}

			err = migration.Execute(sessCtx, db)
			if err != nil {
				return nil, err
			}

			err = m.recordExecution(sessCtx, db, migration)
			if err != nil {
				return nil, err
			}
		}

		return nil, nil
	})

	return err
}

func (m *Migrator) hasExecuted(ctx context.Context, db *mongo.Database, migration Migration) (bool, error) {
	coll := db.Collection(collectionName)
	count, err := coll.CountDocuments(ctx, bson.D{{versionKey, migration.Version()}})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (m *Migrator) recordExecution(ctx context.Context, db *mongo.Database, migration Migration) error {
	coll := db.Collection(collectionName)
	_, err := coll.InsertOne(ctx, bson.D{{versionKey, migration.Version()}})
	if err != nil {
		return err
	}

	return nil
}
