package service

import (
	"context"
	"time"

	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/domain/entity"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/domain/entity/exporter"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/domain/repository"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/logger"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmlogrus"
)

type TimeRecordService struct {
	TimeRecordRepository repository.TimeRecordRepositoryInterface
	EmployeeService      *EmployeeService
}

func NewTimeRecordService(timeRecordRepository repository.TimeRecordRepositoryInterface, employeeService *EmployeeService) *TimeRecordService {
	return &TimeRecordService{
		TimeRecordRepository: timeRecordRepository,
		EmployeeService:      employeeService,
	}
}

func (s *TimeRecordService) RegisterTimeRecord(ctx context.Context, _time time.Time, description, employeeID, createdBy string) (*string, error) {
	span, ctx := apm.StartSpan(ctx, "Register", "time record domain service")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	timeRecord, err := entity.NewTimeRecord(_time, description, employeeID, createdBy)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}
	log.WithField("timeRecord", timeRecord).Info("timeRecord created")

	timeRecordID, err := s.TimeRecordRepository.RegisterTimeRecord(ctx, timeRecord)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}
	log.WithField("timeRecord", timeRecord).Info("timeRecord registered")

	return timeRecordID, nil
}

func (s *TimeRecordService) ApproveTimeRecord(ctx context.Context, id, employeeID string) error {
	span, ctx := apm.StartSpan(ctx, "Approve", "time record domain service")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	timeRecord, err := s.TimeRecordRepository.FindTimeRecord(ctx, id)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return err
	}
	log.WithField("timeRecord", timeRecord).Info("timeRecord finded")

	err = timeRecord.Approve(employeeID)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return err
	}
	log.WithField("timeRecord", timeRecord).Info("timeRecord approved")

	err = s.TimeRecordRepository.SaveTimeRecord(ctx, timeRecord)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return err
	}

	return nil
}

func (s *TimeRecordService) RefuseTimeRecord(ctx context.Context, id, refusedReason, employeeID string) error {
	span, ctx := apm.StartSpan(ctx, "Refuse", "time record domain service")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	timeRecord, err := s.TimeRecordRepository.FindTimeRecord(ctx, id)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return err
	}
	log.WithField("timeRecord", timeRecord).Info("timeRecord finded")

	err = timeRecord.Refuse(employeeID, refusedReason)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return err
	}
	log.WithField("timeRecord", timeRecord).Info("timeRecord refused")

	err = s.TimeRecordRepository.SaveTimeRecord(ctx, timeRecord)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return err
	}

	return nil
}

func (s *TimeRecordService) FindTimeRecord(ctx context.Context, id string) (*entity.TimeRecord, error) {
	span, ctx := apm.StartSpan(ctx, "Find", "time record domain service")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	timeRecord, err := s.TimeRecordRepository.FindTimeRecord(ctx, id)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}
	log.WithField("timeRecord", timeRecord).Info("timeRecord finded")
	return timeRecord, nil
}

func (s *TimeRecordService) SearchTimeRecords(ctx context.Context, fromDate, toDate time.Time, status int, employeeID, approvedBy, refusedBy, createdBy string, pageSize int, pageToken string) (*string, []*entity.TimeRecord, error) {
	span, ctx := apm.StartSpan(ctx, "FindAllByEmployeeID", "time record domain service")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	filter, err := entity.NewFilter(fromDate, toDate, status, employeeID, approvedBy, refusedBy, createdBy, pageSize, pageToken)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, nil, err
	}

	nextPageToken, timeRecords, err := s.TimeRecordRepository.SearchTimeRecords(ctx, filter)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, nil, err
	}
	log.WithField("timeRecords", timeRecords).Info("timeRecords finded")
	return nextPageToken, timeRecords, nil
}

func (s *TimeRecordService) ExportTimeRecords(ctx context.Context, fromDate, toDate time.Time, status int, employeeID, approvedBy, refusedBy, createdBy string, pageSize int, pageToken, accessToken string) (*string, []*exporter.Register, error) {
	filter, err := entity.NewFilter(fromDate, toDate, status, employeeID, approvedBy, refusedBy, createdBy, pageSize, pageToken)
	if err != nil {
		return nil, nil, err
	}

	nextPageToken, timeRecords, err := s.TimeRecordRepository.SearchTimeRecords(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	// group time records by employe id
	groupEmployeeTRs := make(map[string][]*entity.TimeRecord)
	for _, tr := range timeRecords {
		groupEmployeeTRs[tr.EmployeeID] = append(groupEmployeeTRs[tr.EmployeeID], tr)
	}

	// create employees with their own time records
	var employees []*entity.Employee
	for employeeID := range groupEmployeeTRs {
		employee, err := s.EmployeeService.FindEmployee(ctx, employeeID, accessToken)
		if err != nil {
			return nil, nil, err
		}
		employee.AddTimeRecord(groupEmployeeTRs[employeeID]...)
		employees = append(employees, employee)
	}

	exporter, err := exporter.NewExporter(exporter.SECULLUM, employees)
	if err != nil {
		return nil, nil, err
	}

	return nextPageToken, exporter.Export(), nil
}
