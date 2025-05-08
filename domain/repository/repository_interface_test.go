package repository

// import (
// 	"context"
// 	"testing"

// 	"github.com/patricksferraz/time-record-service/domain/entity"
// 	"github.com/stretchr/testify/require"
// )

// type repository struct{}

// func (r *repository) RegisterTimeRecord(ctx context.Context, timeRecord *entity.TimeRecord) (*string, error) {
// 	return new(string), nil
// }

// func (r *repository) SaveTimeRecord(ctx context.Context, timeRecord *entity.TimeRecord) error {
// 	return nil
// }

// func (r *repository) FindTimeRecord(ctx context.Context, id string) (*entity.TimeRecord, error) {
// 	return &entity.TimeRecord{}, nil
// }

// func (r *repository) SearchTimeRecords(ctx context.Context, filter *entity.Filter) (*string, []*entity.TimeRecord, error) {
// 	var timeRecords []*entity.TimeRecord
// 	return nil, timeRecords, nil
// }

// func TestRepository_TimeRecordRepositoryInterface(t *testing.T) {
// 	require.Implements(t, (*TimeRecordRepositoryInterface)(nil), new(repository))
// }
