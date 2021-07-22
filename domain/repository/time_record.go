package repository

import (
	"context"

	"github.com/c-4u/time-record-service/domain/entity"
)

type TimeRecordRepositoryInterface interface {
	RegisterTimeRecord(ctx context.Context, timeRecord *entity.TimeRecord) (*string, error)
	SaveTimeRecord(ctx context.Context, timeRecord *entity.TimeRecord) error
	FindTimeRecord(ctx context.Context, id string) (*entity.TimeRecord, error)
	SearchTimeRecords(ctx context.Context, filter *entity.Filter) (*string, []*entity.TimeRecord, error)
}
