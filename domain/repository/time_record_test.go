package repository

import (
	"context"
	"testing"
	"time"

	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/domain/model"
	"github.com/stretchr/testify/require"
)

type repository struct{}

func (r *repository) Register(ctx context.Context, timeRecord *model.TimeRecord) error {
	return nil
}

func (r *repository) Save(ctx context.Context, timeRecord *model.TimeRecord) error {
	return nil
}

func (r *repository) Find(ctx context.Context, id string) (*model.TimeRecord, error) {
	return &model.TimeRecord{}, nil
}

func (r *repository) FindAllByEmployeeID(ctx context.Context, employeeID string, fromDate, toDate time.Time) ([]*model.TimeRecord, error) {
	var timeRecords []*model.TimeRecord
	return timeRecords, nil
}

func TestRepository_TimeRecordRepositoryInterface(t *testing.T) {
	require.Implements(t, (*TimeRecordRepositoryInterface)(nil), new(repository))
}
