package repository

import (
	"context"

	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/domain/entity"
)

type EmployeeRepositoryInterface interface {
	FindEmployee(ctx context.Context, id string) (*entity.Employee, error)
}
