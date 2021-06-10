package repository

import (
	"context"
	"time"

	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/domain/model"
)

type TimeRecordRepositoryInterface interface {
	Register(ctx context.Context, timeRecord *model.TimeRecord) error
	Save(ctx context.Context, timeRecord *model.TimeRecord) error
	Find(ctx context.Context, id string) (*model.TimeRecord, error)
	FindAllByEmployeeID(ctx context.Context, employeeID string, fromDate, toDate time.Time) ([]*model.TimeRecord, error)
}
