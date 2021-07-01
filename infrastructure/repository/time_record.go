package repository

import (
	"context"
	"time"

	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/domain/entity"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/infrastructure/db"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/infrastructure/db/collection"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/logger"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmlogrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TimeRecordRepository struct {
	M *db.Mongo
}

func (t *TimeRecordRepository) RegisterTimeRecord(ctx context.Context, timeRecord *entity.TimeRecord) (*string, error) {
	span, ctx := apm.StartSpan(ctx, "Register", "repository")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	collection := t.M.Database.Collection(collection.TimeRecordCollection)
	res, err := collection.InsertOne(ctx, timeRecord)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}
	log.WithField("result", res).Info("InsertOne database result")

	return &timeRecord.ID, nil
}

func (t *TimeRecordRepository) SaveTimeRecord(ctx context.Context, timeRecord *entity.TimeRecord) error {
	span, ctx := apm.StartSpan(ctx, "Save", "repository")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	collection := t.M.Database.Collection(collection.TimeRecordCollection)
	res, err := collection.ReplaceOne(ctx, bson.M{"_id": timeRecord.ID}, timeRecord)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return err
	}
	log.WithField("result", res).Info("ReplaceOne database result")

	return nil
}

func (t *TimeRecordRepository) FindTimeRecord(ctx context.Context, id string) (*entity.TimeRecord, error) {
	span, ctx := apm.StartSpan(ctx, "Find", "repository")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	var timeRecord *entity.TimeRecord
	collection := t.M.Database.Collection(collection.TimeRecordCollection)
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&timeRecord)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}
	log.WithField("timeRecord", timeRecord).Info("timeRecord finded")

	return timeRecord, err
}

func (t *TimeRecordRepository) SearchTimeRecords(ctx context.Context, employeeID string, fromDate, toDate time.Time) ([]*entity.TimeRecord, error) {
	span, ctx := apm.StartSpan(ctx, "FindAllByEmployeeID", "repository")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	var timeRecords []*entity.TimeRecord
	collection := t.M.Database.Collection(collection.TimeRecordCollection)

	findOpts := options.Find()
	findOpts.SetSort(bson.M{"time": -1})
	cur, err := collection.Find(
		ctx,
		bson.M{
			"employee_id": employeeID,
			"time": bson.M{
				"$gte": fromDate,
				"$lte": toDate,
			},
		},
		findOpts,
	)
	log.WithField("findOpts", findOpts).Info("database findOpts")

	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}

	for cur.Next(ctx) {
		var timeRecord *entity.TimeRecord
		err := cur.Decode(&timeRecord)
		if err != nil {
			log.WithError(err)
			apm.CaptureError(ctx, err).Send()
			return nil, err
		}
		timeRecords = append(timeRecords, timeRecord)
	}
	log.WithField("timeRecords", timeRecords).Info("timeRecords finded")

	if err := cur.Err(); err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}

	cur.Close(ctx)

	return timeRecords, nil
}

func NewTimeRecordRepository(database *db.Mongo) *TimeRecordRepository {
	return &TimeRecordRepository{
		M: database,
	}
}
