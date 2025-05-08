package db

import (
	"context"

	_ "github.com/patricksferraz/time-record-service/infrastructure/db/migrations"
	migrate "github.com/xakep666/mongo-migrate"
	"go.elastic.co/apm/module/apmmongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Mongo struct {
	Database *mongo.Database
}

func (m *Mongo) Connect(ctx context.Context, uri string, database string) error {
	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(uri),
		options.Client().SetMonitor(apmmongo.CommandMonitor()),
	)
	if err != nil {
		return err
	}

	m.Database = client.Database(database)
	return nil
}

func (m *Mongo) Close(ctx context.Context) error {
	err := m.Database.Client().Disconnect(ctx)
	return err
}

func (m *Mongo) Test(ctx context.Context) error {
	err := m.Database.Client().Ping(ctx, readpref.Primary())
	return err
}

func (m *Mongo) Migrate() error {
	migrate.SetDatabase(m.Database)
	if err := migrate.Up(migrate.AllAvailable); err != nil {
		return err
	}
	return nil
}

func NewMongo(ctx context.Context, uri, dbName string) (*Mongo, error) {
	mongo := &Mongo{}

	err := mongo.Connect(ctx, uri, dbName)
	if err != nil {
		return nil, err
	}

	return mongo, nil
}
