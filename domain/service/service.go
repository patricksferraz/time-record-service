package service

import (
	"context"
	"time"

	"github.com/patricksferraz/time-record-service/domain/entity"
	"github.com/patricksferraz/time-record-service/domain/entity/exporter"
	"github.com/patricksferraz/time-record-service/domain/repository"
	"github.com/patricksferraz/time-record-service/infrastructure/external/topic"
)

type Service struct {
	Repository repository.RepositoryInterface
}

func NewService(timeRecordRepository repository.RepositoryInterface) *Service {
	return &Service{
		Repository: timeRecordRepository,
	}
}

func (s *Service) RegisterTimeRecord(ctx context.Context, _time time.Time, description, employeeID, companyID, createdBy string) (*string, error) {
	// span, ctx := apm.StartSpan(ctx, "Register", "time record domain service")
	// defer span.End()

	// log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))
	employee, err := s.Repository.FindEmployee(ctx, employeeID)
	if err != nil {
		return nil, err
	}

	creater, err := s.Repository.FindEmployee(ctx, createdBy)
	if err != nil {
		return nil, err
	}

	company, err := employee.GetCompany(companyID)
	if err != nil {
		return nil, err
	}

	timeRecord, err := entity.NewTimeRecord(_time, description, employee, creater, company)
	if err != nil {
		// log.WithError(err)
		// apm.CaptureError(ctx, err).Send()
		return nil, err
	}
	// log.WithField("timeRecord", timeRecord).Info("timeRecord created")

	err = s.Repository.RegisterTimeRecord(ctx, timeRecord)
	if err != nil {
		// log.WithError(err)
		// apm.CaptureError(ctx, err).Send()
		return nil, err
	}
	// log.WithField("timeRecord", timeRecord).Info("timeRecord registered")

	event, err := entity.NewTimeRecordEvent(timeRecord)
	if err != nil {
		return nil, err
	}

	msg, err := event.ToJson()
	if err != nil {
		return nil, err
	}

	err = s.Repository.PublishEvent(ctx, string(msg), topic.NEW_TIME_RECORD, timeRecord.ID)
	if err != nil {
		return nil, err
	}

	return &timeRecord.ID, nil
}

func (s *Service) ApproveTimeRecord(ctx context.Context, id, approvedBy string) error {
	timeRecord, err := s.Repository.FindTimeRecord(ctx, id)
	if err != nil {
		return err
	}

	approver, err := s.Repository.FindEmployee(ctx, approvedBy)
	if err != nil {
		return err
	}

	err = timeRecord.Approve(approver)
	if err != nil {
		return err
	}

	err = s.Repository.SaveTimeRecord(ctx, timeRecord)
	if err != nil {
		return err
	}

	event, err := entity.NewTimeRecordEvent(timeRecord)
	if err != nil {
		return err
	}

	msg, err := event.ToJson()
	if err != nil {
		return err
	}

	err = s.Repository.PublishEvent(ctx, string(msg), topic.APPROVE_TIME_RECORD, timeRecord.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) RefuseTimeRecord(ctx context.Context, id, refusedReason, refusedBy string) error {
	timeRecord, err := s.Repository.FindTimeRecord(ctx, id)
	if err != nil {
		return err
	}

	refuser, err := s.Repository.FindEmployee(ctx, refusedBy)
	if err != nil {
		return err
	}

	err = timeRecord.Refuse(refuser, refusedReason)
	if err != nil {
		return err
	}

	err = s.Repository.SaveTimeRecord(ctx, timeRecord)
	if err != nil {
		return err
	}

	event, err := entity.NewTimeRecordEvent(timeRecord)
	if err != nil {
		return err
	}

	msg, err := event.ToJson()
	if err != nil {
		return err
	}

	err = s.Repository.PublishEvent(ctx, string(msg), topic.REFUSE_TIME_RECORD, timeRecord.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) FindTimeRecord(ctx context.Context, id string) (*entity.TimeRecord, error) {
	timeRecord, err := s.Repository.FindTimeRecord(ctx, id)
	if err != nil {
		return nil, err
	}

	return timeRecord, nil
}

func (s *Service) SearchTimeRecords(ctx context.Context, fromDate, toDate time.Time, status int, employeeID, approvedBy, refusedBy, createdBy, companyID string, pageSize int, pageToken string) (*string, []*entity.TimeRecord, error) {
	filter, err := entity.NewFilter(fromDate, toDate, status, employeeID, approvedBy, refusedBy, createdBy, companyID, pageSize, pageToken)
	if err != nil {
		return nil, nil, err
	}

	nextPageToken, timeRecords, err := s.Repository.SearchTimeRecords(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	return nextPageToken, timeRecords, nil
}

func (s *Service) ExportTimeRecords(ctx context.Context, fromDate, toDate time.Time, status int, employeeID, approvedBy, refusedBy, createdBy, companyID string, pageSize int, pageToken, accessToken string) (*string, []*exporter.Register, error) {
	filter, err := entity.NewFilter(fromDate, toDate, status, employeeID, approvedBy, refusedBy, createdBy, companyID, pageSize, pageToken)
	if err != nil {
		return nil, nil, err
	}

	nextPageToken, timeRecords, err := s.Repository.SearchTimeRecords(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	exporter, err := exporter.NewExporter(exporter.SECULLUM, timeRecords)
	if err != nil {
		return nil, nil, err
	}

	return nextPageToken, exporter.Export(), nil
}

func (s *Service) CreateCompany(ctx context.Context, id string) error {
	company, err := entity.NewCompany(id)
	if err != nil {
		return err
	}

	err = s.Repository.CreateCompany(ctx, company)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) CreateEmployee(ctx context.Context, id, pis string) error {
	employee, err := entity.NewEmployee(id, pis)
	if err != nil {
		return err
	}

	err = s.Repository.CreateEmployee(ctx, employee)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) AddEmployeeToCompany(ctx context.Context, companyID, employeeID string) error {
	company, err := s.Repository.FindCompany(ctx, companyID)
	if err != nil {
		return err
	}

	employee, err := s.Repository.FindEmployee(ctx, employeeID)
	if err != nil {
		return err
	}

	err = employee.AddCompany(company)
	if err != nil {
		return err
	}

	err = s.Repository.SaveEmployee(ctx, employee)
	if err != nil {
		return err
	}

	return nil
}
