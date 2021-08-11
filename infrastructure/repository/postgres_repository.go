package repository

import (
	"context"
	"fmt"

	"github.com/c-4u/time-record-service/domain/entity"
	"github.com/c-4u/time-record-service/infrastructure/db"
)

type PostgresRepository struct {
	P *db.Postgres
}

func NewPostgresRepository(db *db.Postgres) *PostgresRepository {
	return &PostgresRepository{
		P: db,
	}
}

func (r *PostgresRepository) CreateEmployee(ctx context.Context, employee *entity.Employee) error {
	err := r.P.Db.Create(employee).Error
	return err
}

func (r *PostgresRepository) FindEmployee(ctx context.Context, id string) (*entity.Employee, error) {
	var employee entity.Employee
	r.P.Db.First(&employee, "id = ?", id)

	if employee.ID == "" {
		return nil, fmt.Errorf("no employee found")
	}

	return &employee, nil
}

func (r *PostgresRepository) SaveEmployee(ctx context.Context, employee *entity.Employee) error {
	err := r.P.Db.Save(employee).Error
	return err
}

func (r *PostgresRepository) RegisterTimeRecord(ctx context.Context, timeRecord *entity.TimeRecord) error {
	err := r.P.Db.Create(timeRecord).Error
	return err
}

func (r *PostgresRepository) SaveTimeRecord(ctx context.Context, timeRecord *entity.TimeRecord) error {
	err := r.P.Db.Save(timeRecord).Error
	return err
}

func (r *PostgresRepository) FindTimeRecord(ctx context.Context, id string) (*entity.TimeRecord, error) {
	var timeRecord entity.TimeRecord
	r.P.Db.Preload("Employee").First(&timeRecord, "id = ?", id)

	if timeRecord.ID == "" {
		return nil, fmt.Errorf("no timeRecord found")
	}

	return &timeRecord, nil
}

func (r *PostgresRepository) SearchTimeRecords(ctx context.Context, filter *entity.Filter) (*string, []*entity.TimeRecord, error) {
	var timeRecords []*entity.TimeRecord

	q := r.P.Db.Order("token desc").Limit(filter.PageSize).Preload("Employee")

	if !filter.FromDate.IsZero() {
		q = q.Where("time >= ?", filter.FromDate)
	}
	if !filter.ToDate.IsZero() {
		q = q.Where("time <= ?", filter.ToDate)
	}

	if filter.Status != 0 {
		q = q.Where("status = ?", filter.Status)
	}
	if filter.EmployeeID != "" {
		q = q.Where("employee_id = ?", filter.EmployeeID)
	}
	if filter.ApprovedBy != "" {
		q = q.Where("approved_by = ?", filter.ApprovedBy)
	}
	if filter.RefusedBy != "" {
		q = q.Where("refused_by = ?", filter.RefusedBy)
	}
	if filter.CreatedBy != "" {
		q = q.Where("created_by = ?", filter.CreatedBy)
	}
	if filter.PageToken != "" {
		q = q.Where("token < ?", filter.PageToken)
	}

	err := q.Find(&timeRecords).Error
	if err != nil {
		return nil, nil, err
	}

	length := len(timeRecords)
	var nextPageToken string
	if length == filter.PageSize {
		nextPageToken = *timeRecords[length-1].Token
	}

	return &nextPageToken, timeRecords, nil
}
