package repository

import (
	"context"

	"github.com/c-4u/time-record-service/domain/entity"
)

type RepositoryInterface interface {
	CreateEmployee(ctx context.Context, employee *entity.Employee) error
	FindEmployee(ctx context.Context, id string) (*entity.Employee, error)
	SaveEmployee(ctx context.Context, employee *entity.Employee) error

	RegisterTimeRecord(ctx context.Context, timeRecord *entity.TimeRecord) error
	SaveTimeRecord(ctx context.Context, timeRecord *entity.TimeRecord) error
	FindTimeRecord(ctx context.Context, id string) (*entity.TimeRecord, error)
	SearchTimeRecords(ctx context.Context, filter *entity.Filter) (*string, []*entity.TimeRecord, error)
}
