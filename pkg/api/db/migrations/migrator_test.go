package migrations_test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/yukitsune/maestro/pkg/api/db/migrations"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

type mockMigration1 struct {
	DidExecute bool
}

func (m *mockMigration1) Execute(_ context.Context, _ *mongo.Database) error {
	m.DidExecute = true
	return nil
}

func (m *mockMigration1) Version() int {
	return 1
}

type mockMigration2 struct {
	DidExecute bool
}

func (m *mockMigration2) Execute(_ context.Context, _ *mongo.Database) error {
	m.DidExecute = true
	return nil
}

func (m *mockMigration2) Version() int {
	return 2
}

type mockMigration3 struct{}

func (m *mockMigration3) Execute(_ context.Context, _ *mongo.Database) error {
	return fmt.Errorf("oops")
}

func (m *mockMigration3) Version() int {
	return 3
}

type mockMigrationProvider struct {
	migrations []migrations.Migration
}

func (mp *mockMigrationProvider) Migrations() []migrations.Migration {
	return mp.migrations
}

func Test_MigrationsAreRecorded(t *testing.T) {
	withTestDb(t, func(db *mongo.Database) {

		// Arrange
		m1 := &mockMigration1{}
		m2 := &mockMigration2{}
		mp := &mockMigrationProvider{
			[]migrations.Migration{
				m1,
				m2,
			},
		}
		m := &migrations.Migrator{}

		// Act
		err := m.Execute(context.Background(), mp, db)
		assert.NoError(t, err)

		// Assert
		assertMigrationsExecuted(t, db, m1, m2)
	})
}

func Test_ExecutedMigrationsAreSkipped(t *testing.T) {
	withTestDb(t, func(db *mongo.Database) {

		// Arrange
		// Execute the first migration
		mp1 := &mockMigrationProvider{
			[]migrations.Migration{
				&mockMigration1{},
			},
		}
		mr1 := &migrations.Migrator{}

		err := mr1.Execute(context.Background(), mp1, db)
		assert.NoError(t, err)

		m1 := &mockMigration1{}
		m2 := &mockMigration2{}
		mp2 := &mockMigrationProvider{
			[]migrations.Migration{
				m1,
				m2,
			},
		}
		mr2 := &migrations.Migrator{}

		// Act
		// Try to execute the first migration again, along with the second one
		err = mr2.Execute(context.Background(), mp2, db)
		assert.NoError(t, err)

		// Assert
		// First one should be skipped (already executed)
		// Second one should be executed (not yet executed)
		assert.False(t, m1.DidExecute)
		assert.True(t, m2.DidExecute)
		assertMigrationsExecuted(t, db, m1, m2)
	})
}

func Test_BadMigrationsAbortAllChanges(t *testing.T) {
	withTestDb(t, func(db *mongo.Database) {

		// Arrange
		mp := &mockMigrationProvider{
			[]migrations.Migration{
				&mockMigration1{},
				&mockMigration1{},
				&mockMigration3{}, // Bad migration
			},
		}
		mr := &migrations.Migrator{}

		// Act
		err := mr.Execute(context.Background(), mp, db)
		assert.NotNil(t, err)

		// Assert
		assertNoMigrationsExecuted(t, db)
	})
}

func assertMigrationsExecuted(t *testing.T, db *mongo.Database, migrations ...migrations.Migration) {
	migColl := db.Collection("migrations")
	for _, m := range migrations {
		count, err := migColl.CountDocuments(context.Background(), bson.D{{"version", m.Version()}})
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)
	}
}

func assertNoMigrationsExecuted(t *testing.T, db *mongo.Database) {
	migColl := db.Collection("migrations")
	count, err := migColl.CountDocuments(context.Background(), bson.D{})
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)
}
