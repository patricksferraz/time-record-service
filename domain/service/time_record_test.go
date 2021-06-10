package service_test

import (
	"context"
	"log"
	"testing"
	"time"

	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/domain/model"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/domain/service"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/infrastructure/db"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/infrastructure/repository"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/utils"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"syreclabs.com/go/faker"
)

func TestService_Register(t *testing.T) {

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
	timeRecord, err := timeRecordService.Register(ctx, _time, description, employeeID)

	require.Nil(t, err)
	require.NotEmpty(t, uuid.FromStringOrNil(timeRecord.ID))
	require.Equal(t, timeRecord.Time, _time)
	require.Equal(t, timeRecord.Status, model.APPROVED)
	require.Equal(t, timeRecord.Description, description)
	require.Equal(t, timeRecord.RegularTime, true)
	require.Equal(t, timeRecord.EmployeeID, employeeID)

	_, err = timeRecordService.Register(ctx, _time.AddDate(0, 0, 1), description, employeeID)
	require.NotNil(t, err)
}

func TestService_Approve(t *testing.T) {

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
	timeRecord, _ := timeRecordService.Register(ctx, _time, description, employeeID)

	approvedBy := uuid.NewV4().String()
	_, err = timeRecordService.Approve(ctx, "", approvedBy)
	require.NotNil(t, err)
	_, err = timeRecordService.Approve(ctx, timeRecord.ID, "")
	require.NotNil(t, err)

	timeRecord, err = timeRecordService.Approve(ctx, timeRecord.ID, approvedBy)
	require.Nil(t, err)
	require.Equal(t, timeRecord.ApprovedBy, approvedBy)
}

func TestService_Refuse(t *testing.T) {

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
	timeRecord, _ := timeRecordService.Register(ctx, _time, description, employeeID)

	auditedBy := uuid.NewV4().String()
	_, err = timeRecordService.Refuse(ctx, "", auditedBy, description)
	require.NotNil(t, err)
	_, err = timeRecordService.Refuse(ctx, timeRecord.ID, "", description)
	require.NotNil(t, err)
	_, err = timeRecordService.Refuse(ctx, timeRecord.ID, auditedBy, "")
	require.NotNil(t, err)

	timeRecord, err = timeRecordService.Refuse(ctx, timeRecord.ID, description, auditedBy)
	require.Nil(t, err)
	require.Equal(t, timeRecord.RefusedBy, auditedBy)
}

func TestService_Find(t *testing.T) {

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
	timeRecord, _ := timeRecordService.Register(ctx, _time, description, employeeID)

	_, err = timeRecordService.Find(ctx, timeRecord.ID)
	require.Nil(t, err)
	require.True(t, timeRecord.Time.Equal(_time))
	require.Equal(t, timeRecord.Description, description)
	require.Equal(t, timeRecord.EmployeeID, employeeID)
	_, err = timeRecordService.Find(ctx, "")
	require.NotNil(t, err)
}

func TestService_FindAllByEmployeeID(t *testing.T) {

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
	timeRecord, _ := timeRecordService.Register(ctx, _time, description, employeeID)

	fromDate := _time.AddDate(0, 0, -1)
	toDate := _time.AddDate(0, 0, 1)
	trs, err := timeRecordService.FindAllByEmployeeID(ctx, timeRecord.EmployeeID, fromDate, toDate)
	require.Nil(t, err)
	require.Len(t, trs, 1)
	require.NotEmpty(t, trs)
	require.Equal(t, timeRecord.ID, trs[0].ID)
	require.Equal(t, timeRecord.Time.Unix(), trs[0].Time.Unix())
	require.Equal(t, timeRecord.Description, trs[0].Description)
	require.Equal(t, timeRecord.RegularTime, trs[0].RegularTime)
	require.Equal(t, timeRecord.Status, trs[0].Status)
	require.Equal(t, timeRecord.ApprovedBy, trs[0].ApprovedBy)
	require.Equal(t, timeRecord.RefusedBy, trs[0].RefusedBy)
	require.Equal(t, timeRecord.RefusedReason, trs[0].RefusedReason)
	require.Equal(t, timeRecord.CreatedAt.Unix(), trs[0].CreatedAt.Unix())
	require.Equal(t, timeRecord.UpdatedAt.Unix(), trs[0].UpdatedAt.Unix())

	trs, err = timeRecordService.FindAllByEmployeeID(ctx, "", fromDate, toDate)
	require.Nil(t, err)
	require.Empty(t, trs)
}
