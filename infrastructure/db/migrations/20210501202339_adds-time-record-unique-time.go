package migrations

import (
	"context"

	"github.com/patricksferraz/time-record-service/infrastructure/db/collection"
	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	opt := options.Index().SetName("time_record_id_time_unique")
	opt.SetUnique(true)
	keys := bson.D{primitive.E{Key: "employee_id", Value: 1}, primitive.E{Key: "time", Value: 1}}
	index := mongo.IndexModel{Keys: keys, Options: opt}

	migrate.Register(func(db *mongo.Database) error {
		_, err := db.Collection(collection.TimeRecordCollection).Indexes().CreateOne(context.TODO(), index)
		if err != nil {
			return err
		}
		return nil
	}, func(db *mongo.Database) error {
		_, err := db.Collection(collection.TimeRecordCollection).Indexes().DropOne(context.TODO(), *opt.Name)
		if err != nil {
			return err
		}
		return nil
	})
}
