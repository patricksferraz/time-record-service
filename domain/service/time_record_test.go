package service_test

import (
	"context"
	"log"
	"testing"
	"time"

	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/domain/service"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/infrastructure/db"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/infrastructure/repository"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/utils"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"syreclabs.com/go/faker"
)

func TestService_RegisterTimeRecord(t *testing.T) {

	ctx := context.Background()
	uri := utils.GetEnv("DB_URI", "mongodb://localhost")
	dbName := utils.GetEnv("DB_NAME", "test")
	db, err := db.NewMongo(ctx, uri, dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(ctx)
	defer db.Database.Drop(ctx)

	timeRecordRepository := repository.NewTimeRecordRepository(db)
	timeRecordService := service.NewTimeRecordService(timeRecordRepository)

	_time := time.Now()
	description := faker.Lorem().Sentence(10)
	employeeID := uuid.NewV4().String()
	timeRecordID, err := timeRecordService.RegisterTimeRecord(ctx, _time, description, employeeID, employeeID)

	require.Nil(t, err)
	require.NotEmpty(t, uuid.FromStringOrNil(*timeRecordID))

	_, err = timeRecordService.RegisterTimeRecord(ctx, _time.AddDate(0, 0, 1), description, employeeID, employeeID)
	require.NotNil(t, err)
}

func TestService_ApproveTimeRecord(t *testing.T) {

	ctx := context.Background()
	uri := utils.GetEnv("DB_URI", "mongodb://localhost")
	dbName := utils.GetEnv("DB_NAME", "test")
	db, err := db.NewMongo(ctx, uri, dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(ctx)
	defer db.Database.Drop(ctx)

	timeRecordRepository := repository.NewTimeRecordRepository(db)
	timeRecordService := service.NewTimeRecordService(timeRecordRepository)

	_time := time.Now().AddDate(0, 0, -1)
	description := faker.Lorem().Sentence(10)
	employeeID := uuid.NewV4().String()
	timeRecordID, _ := timeRecordService.RegisterTimeRecord(ctx, _time, description, employeeID, employeeID)

	approvedBy := uuid.NewV4().String()
	err = timeRecordService.ApproveTimeRecord(ctx, "", approvedBy)
	require.NotNil(t, err)
	err = timeRecordService.ApproveTimeRecord(ctx, *timeRecordID, "")
	require.NotNil(t, err)

	err = timeRecordService.ApproveTimeRecord(ctx, *timeRecordID, approvedBy)
	require.Nil(t, err)
}

func TestService_RefuseTimeRecord(t *testing.T) {

	ctx := context.Background()
	uri := utils.GetEnv("DB_URI", "mongodb://localhost")
	dbName := utils.GetEnv("DB_NAME", "test")
	db, err := db.NewMongo(ctx, uri, dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(ctx)
	defer db.Database.Drop(ctx)

	timeRecordRepository := repository.NewTimeRecordRepository(db)
	timeRecordService := service.NewTimeRecordService(timeRecordRepository)

	_time := time.Now().AddDate(0, 0, -1)
	description := faker.Lorem().Sentence(10)
	employeeID := uuid.NewV4().String()
	timeRecordID, _ := timeRecordService.RegisterTimeRecord(ctx, _time, description, employeeID, employeeID)

	auditedBy := uuid.NewV4().String()
	err = timeRecordService.RefuseTimeRecord(ctx, "", auditedBy, description)
	require.NotNil(t, err)
	err = timeRecordService.RefuseTimeRecord(ctx, *timeRecordID, "", description)
	require.NotNil(t, err)
	err = timeRecordService.RefuseTimeRecord(ctx, *timeRecordID, auditedBy, "")
	require.NotNil(t, err)

	err = timeRecordService.RefuseTimeRecord(ctx, *timeRecordID, description, auditedBy)
	require.Nil(t, err)
}

func TestService_FindTimeRecord(t *testing.T) {

	ctx := context.Background()
	uri := utils.GetEnv("DB_URI", "mongodb://localhost")
	dbName := utils.GetEnv("DB_NAME", "test")
	db, err := db.NewMongo(ctx, uri, dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(ctx)
	defer db.Database.Drop(ctx)

	timeRecordRepository := repository.NewTimeRecordRepository(db)
	timeRecordService := service.NewTimeRecordService(timeRecordRepository)

	_time := time.Now().AddDate(0, 0, -1)
	description := faker.Lorem().Sentence(10)
	employeeID := uuid.NewV4().String()
	timeRecordID, _ := timeRecordService.RegisterTimeRecord(ctx, _time, description, employeeID, employeeID)

	timeRecord, err := timeRecordService.FindTimeRecord(ctx, *timeRecordID)
	require.Nil(t, err)
	require.Equal(t, _time.Unix(), timeRecord.Time.Unix())
	require.Equal(t, timeRecord.Description, description)
	require.Equal(t, timeRecord.EmployeeID, employeeID)
	_, err = timeRecordService.FindTimeRecord(ctx, "")
	require.NotNil(t, err)
}

func TestService_SearchTimeRecords(t *testing.T) {

	ctx := context.Background()
	uri := utils.GetEnv("DB_URI", "mongodb://localhost")
	dbName := utils.GetEnv("DB_NAME", "test")
	db, err := db.NewMongo(ctx, uri, dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(ctx)
	defer db.Database.Drop(ctx)

	timeRecordRepository := repository.NewTimeRecordRepository(db)
	timeRecordService := service.NewTimeRecordService(timeRecordRepository)

	_time := time.Now()
	description := faker.Lorem().Sentence(10)
	employeeID := uuid.NewV4().String()
	timeRecordID, _ := timeRecordService.RegisterTimeRecord(ctx, _time, description, employeeID, employeeID)
	timeRecord, _ := timeRecordService.FindTimeRecord(ctx, *timeRecordID)

	fromDate := _time.AddDate(0, 0, -1)
	toDate := _time.AddDate(0, 0, 1)
	trs, err := timeRecordService.SearchTimeRecords(ctx, employeeID, fromDate, toDate)
	require.Nil(t, err)
	require.Len(t, trs, 1)
	require.NotEmpty(t, trs)
	require.Equal(t, *timeRecordID, trs[0].ID)
	require.Equal(t, timeRecord.Time.Unix(), trs[0].Time.Unix())
	require.Equal(t, timeRecord.Description, trs[0].Description)
	require.Equal(t, timeRecord.RegularTime, trs[0].RegularTime)
	require.Equal(t, timeRecord.Status, trs[0].Status)
	require.Equal(t, timeRecord.ApprovedBy, trs[0].ApprovedBy)
	require.Equal(t, timeRecord.RefusedBy, trs[0].RefusedBy)
	require.Equal(t, timeRecord.RefusedReason, trs[0].RefusedReason)
	require.Equal(t, timeRecord.CreatedAt.Unix(), trs[0].CreatedAt.Unix())
	require.Equal(t, timeRecord.UpdatedAt.Unix(), trs[0].UpdatedAt.Unix())

	trs, err = timeRecordService.SearchTimeRecords(ctx, "", fromDate, toDate)
	require.Nil(t, err)
	require.Empty(t, trs)
}
