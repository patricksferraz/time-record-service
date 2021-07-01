package repository

import (
	"context"
	"time"

	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/domain/entity"
)

type TimeRecordRepositoryInterface interface {
	RegisterTimeRecord(ctx context.Context, timeRecord *entity.TimeRecord) (*string, error)
	SaveTimeRecord(ctx context.Context, timeRecord *entity.TimeRecord) error
	FindTimeRecord(ctx context.Context, id string) (*entity.TimeRecord, error)
	SearchTimeRecords(ctx context.Context, employeeID string, fromDate, toDate time.Time) ([]*entity.TimeRecord, error)
}
