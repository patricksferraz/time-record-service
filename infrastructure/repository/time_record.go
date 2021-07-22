package repository

import (
	"context"

	"github.com/c-4u/time-record-service/domain/entity"
	"github.com/c-4u/time-record-service/infrastructure/db"
	"github.com/c-4u/time-record-service/infrastructure/db/collection"
	"github.com/c-4u/time-record-service/logger"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmlogrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	res, err := collection.ReplaceOne(ctx, bson.M{"id": timeRecord.ID}, timeRecord)
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
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&timeRecord)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}
	log.WithField("timeRecord", timeRecord).Info("timeRecord finded")

	return timeRecord, err
}

func (t *TimeRecordRepository) SearchTimeRecords(ctx context.Context, filter *entity.Filter) (*string, []*entity.TimeRecord, error) {
	span, ctx := apm.StartSpan(ctx, "FindAllByEmployeeID", "repository")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	collection := t.M.Database.Collection(collection.TimeRecordCollection)

	findOpts := options.Find()
	findOpts.SetLimit(int64(filter.PageSize))
	findOpts.SetSort(bson.M{"_id": 1})

	_t := []bson.M{}
	if !filter.FromDate.IsZero() {
		_t = append(_t, bson.M{
			"time": bson.M{
				"$gte": filter.FromDate,
			}})
	}
	if !filter.ToDate.IsZero() {
		_t = append(_t, bson.M{
			"time": bson.M{
				"$lte": filter.ToDate,
			}})
	}

	f := bson.M{}
	if len(_t) > 0 {
		f["$and"] = _t
	}
	if filter.Status != 0 {
		f["status"] = filter.Status
	}
	if filter.EmployeeID != "" {
		f["employee_id"] = filter.EmployeeID
	}
	if filter.ApprovedBy != "" {
		f["approved_by"] = filter.ApprovedBy
	}
	if filter.RefusedBy != "" {
		f["refused_by"] = filter.RefusedBy
	}
	if filter.CreatedBy != "" {
		f["created_by"] = filter.CreatedBy
	}
	if filter.PageToken != "" {
		token, err := primitive.ObjectIDFromHex(filter.PageToken)
		if err != nil {
			log.Fatal(err)
		}
		f["_id"] = bson.M{"$gt": token}
	}

	cur, err := collection.Find(ctx, f, findOpts)
	log.WithField("findOpts", findOpts).Info("database findOpts")

	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, nil, err
	}

	var fullBatch bool = true
	if cur.RemainingBatchLength() < int(filter.PageSize) {
		fullBatch = false
	}

	var nextPageToken string
	var timeRecords []*entity.TimeRecord
	for cur.Next(ctx) {
		var timeRecord *entity.TimeRecord
		err := cur.Decode(&timeRecord)
		if err != nil {
			log.WithError(err)
			apm.CaptureError(ctx, err).Send()
			return nil, nil, err
		}
		if cur.RemainingBatchLength() == 0 && fullBatch {
			nextPageToken = cur.Current.Lookup("_id").ObjectID().Hex()
		}
		timeRecords = append(timeRecords, timeRecord)
	}
	log.WithField("timeRecords", timeRecords).Info("timeRecords finded")

	if err := cur.Err(); err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, nil, err
	}

	cur.Close(ctx)

	return &nextPageToken, timeRecords, nil
}

func NewTimeRecordRepository(database *db.Mongo) *TimeRecordRepository {
	return &TimeRecordRepository{
		M: database,
	}
}
