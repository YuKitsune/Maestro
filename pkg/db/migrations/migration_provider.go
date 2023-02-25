package migrations

type MigrationProvider interface {
	Migrations() []Migration
}

type mongoMigrationProvider struct {
}

func NewMongoMigrationProvider() MigrationProvider {
	return &mongoMigrationProvider{}
}

func (mp *mongoMigrationProvider) Migrations() []Migration {
	return []Migration{
		&Migration0001SplitThings{},
	}
}
