package repository

import (
	"context"
	"fmt"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/patricksferraz/time-record-service/domain/entity"
	"github.com/patricksferraz/time-record-service/infrastructure/db"
	"github.com/patricksferraz/time-record-service/infrastructure/external"
)

type Repository struct {
	P *db.Postgres
	K *external.KafkaProducer
}

func NewRepository(db *db.Postgres, kafkaProducer *external.KafkaProducer) *Repository {
	return &Repository{
		P: db,
		K: kafkaProducer,
	}
}

func (r *Repository) SaveEmployee(ctx context.Context, employee *entity.Employee) error {
	err := r.P.Db.Save(employee).Error
	return err
}

func (r *Repository) RegisterTimeRecord(ctx context.Context, timeRecord *entity.TimeRecord) error {
	err := r.P.Db.Create(timeRecord).Error
	return err
}

func (r *Repository) SaveTimeRecord(ctx context.Context, timeRecord *entity.TimeRecord) error {
	err := r.P.Db.Save(timeRecord).Error
	return err
}

func (r *Repository) FindTimeRecord(ctx context.Context, id string) (*entity.TimeRecord, error) {
	var timeRecord entity.TimeRecord
	r.P.Db.Preload("Employee").First(&timeRecord, "id = ?", id)

	if timeRecord.ID == "" {
		return nil, fmt.Errorf("no timeRecord found")
	}

	return &timeRecord, nil
}

func (r *Repository) SearchTimeRecords(ctx context.Context, filter *entity.Filter) (*string, []*entity.TimeRecord, error) {
	var timeRecords []*entity.TimeRecord

	q := r.P.Db.Order("token desc").Preload("Employee")

	if filter.PageSize != 0 {
		q = q.Limit(filter.PageSize)
	}
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
	if filter.CompanyID != "" {
		q = q.Where("company_id = ?", filter.CompanyID)
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

func (r *Repository) CreateEmployee(ctx context.Context, employee *entity.Employee) error {
	err := r.P.Db.Create(employee).Error
	return err
}

func (r *Repository) FindEmployee(ctx context.Context, id string) (*entity.Employee, error) {
	var employee entity.Employee
	r.P.Db.Preload("Companies").First(&employee, "id = ?", id)

	if employee.ID == "" {
		return nil, fmt.Errorf("no employee found")
	}

	return &employee, nil
}

func (r *Repository) CreateCompany(ctx context.Context, company *entity.Company) error {
	err := r.P.Db.Create(company).Error
	return err
}

func (r *Repository) FindCompany(ctx context.Context, id string) (*entity.Company, error) {
	var company entity.Company
	r.P.Db.First(&company, "id = ?", id)

	if company.ID == "" {
		return nil, fmt.Errorf("no company found")
	}

	return &company, nil
}

func (r *Repository) PublishEvent(ctx context.Context, msg, topic, key string) error {
	message := &ckafka.Message{
		TopicPartition: ckafka.TopicPartition{Topic: &topic, Partition: ckafka.PartitionAny},
		Value:          []byte(msg),
		Key:            []byte(key),
	}
	err := r.K.Producer.Produce(message, r.K.DeliveryChan)
	if err != nil {
		return err
	}
	return nil
}
