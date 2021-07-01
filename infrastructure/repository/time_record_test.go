package repository_test

import (
	"context"
	"testing"
	"time"

	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/domain/entity"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/infrastructure/db"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/infrastructure/repository"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/utils"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"syreclabs.com/go/faker"
)

func TestRepository_Register(t *testing.T) {

	ctx := context.Background()
	uri := utils.GetEnv("DB_URI", "mongodb://localhost")
	dbName := utils.GetEnv("DB_NAME", "test")
	_db, _ := db.NewMongo(ctx, uri, dbName)
	repository := repository.NewTimeRecordRepository(_db)

	now := time.Now()
	description := faker.Lorem().Sentence(10)
	employeeID := uuid.NewV4().String()
	timeRecord, _ := entity.NewTimeRecord(now, description, employeeID, employeeID)

	id, err := repository.RegisterTimeRecord(ctx, timeRecord)
	require.Nil(t, err)
	require.Equal(t, id, timeRecord.ID)
}

func TestRepository_Save(t *testing.T) {

	ctx := context.Background()
	uri := utils.GetEnv("DB_URI", "mongodb://localhost")
	dbName := utils.GetEnv("DB_NAME", "test")
	_db, _ := db.NewMongo(ctx, uri, dbName)
	repository := repository.NewTimeRecordRepository(_db)

	now := time.Now()
	description := faker.Lorem().Sentence(10)
	employeeID := uuid.NewV4().String()
	timeRecord, _ := entity.NewTimeRecord(now, description, employeeID, employeeID)

	repository.RegisterTimeRecord(ctx, timeRecord)

	timeRecord.Description = faker.Lorem().Sentence(10)
	err := repository.SaveTimeRecord(ctx, timeRecord)
	require.Nil(t, err)
}

func TestRepository_Find(t *testing.T) {

	ctx := context.Background()
	uri := utils.GetEnv("DB_URI", "mongodb://localhost")
	dbName := utils.GetEnv("DB_NAME", "test")
	_db, _ := db.NewMongo(ctx, uri, dbName)
	repository := repository.NewTimeRecordRepository(_db)

	// NOTE: time.Time is in nanoseconds and mongodb in milliseconds
	y, m, d := time.Now().Date()
	now := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	description := faker.Lorem().Sentence(10)
	employeeID := uuid.NewV4().String()
	timeRecord, _ := entity.NewTimeRecord(now, description, employeeID, employeeID)

	repository.RegisterTimeRecord(ctx, timeRecord)

	timeRecordDb, err := repository.FindTimeRecord(ctx, timeRecord.ID)
	require.Nil(t, err)
	require.Equal(t, timeRecord.ID, timeRecordDb.ID)
	require.True(t, timeRecord.Time.Equal(timeRecordDb.Time))
	require.Equal(t, timeRecord.Status, timeRecordDb.Status)
	require.Equal(t, timeRecord.Description, timeRecordDb.Description)
	require.Equal(t, timeRecord.RegularTime, timeRecordDb.RegularTime)
	require.Equal(t, timeRecord.EmployeeID, timeRecordDb.EmployeeID)
	require.Equal(t, timeRecord.ApprovedBy, timeRecordDb.ApprovedBy)
	require.NotEmpty(t, timeRecordDb.CreatedAt)
	require.Empty(t, timeRecordDb.UpdatedAt)
}

func TestRepository_FindAllByEmployeeID(t *testing.T) {

	ctx := context.Background()
	uri := utils.GetEnv("DB_URI", "mongodb://localhost")
	dbName := utils.GetEnv("DB_NAME", "test")
	_db, _ := db.NewMongo(ctx, uri, dbName)
	repository := repository.NewTimeRecordRepository(_db)

	// NOTE: time.Time is in nanoseconds and mongodb in milliseconds
	y, m, d := time.Now().Date()
	now := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	description := faker.Lorem().Sentence(10)
	employeeID := uuid.NewV4().String()
	timeRecord, _ := entity.NewTimeRecord(now, description, employeeID, employeeID)

	repository.RegisterTimeRecord(ctx, timeRecord)

	timeRecordsDb, err := repository.SearchTimeRecords(ctx, timeRecord.EmployeeID, now, now)
	require.Nil(t, err)
	require.Equal(t, timeRecord.ID, timeRecordsDb[0].ID)
	require.True(t, timeRecord.Time.Equal(timeRecordsDb[0].Time))
	require.Equal(t, timeRecord.Status, timeRecordsDb[0].Status)
	require.Equal(t, timeRecord.Description, timeRecordsDb[0].Description)
	require.Equal(t, timeRecord.RegularTime, timeRecordsDb[0].RegularTime)
	require.Equal(t, timeRecord.EmployeeID, timeRecordsDb[0].EmployeeID)
	require.Equal(t, timeRecord.ApprovedBy, timeRecordsDb[0].ApprovedBy)
	require.NotEmpty(t, timeRecordsDb[0].CreatedAt)
	require.Empty(t, timeRecordsDb[0].UpdatedAt)
}
