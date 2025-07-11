package repository

import (
	"context"

	"github.com/patricksferraz/time-record-service/domain/entity"
)

type RepositoryInterface interface {
	CreateEmployee(ctx context.Context, employee *entity.Employee) error
	FindEmployee(ctx context.Context, id string) (*entity.Employee, error)
	SaveEmployee(ctx context.Context, employee *entity.Employee) error

	RegisterTimeRecord(ctx context.Context, timeRecord *entity.TimeRecord) error
	SaveTimeRecord(ctx context.Context, timeRecord *entity.TimeRecord) error
	FindTimeRecord(ctx context.Context, id string) (*entity.TimeRecord, error)
	SearchTimeRecords(ctx context.Context, filter *entity.Filter) (*string, []*entity.TimeRecord, error)

	CreateCompany(ctx context.Context, company *entity.Company) error
	FindCompany(ctx context.Context, id string) (*entity.Company, error)

	PublishEvent(ctx context.Context, msg, topic, key string) error
}
