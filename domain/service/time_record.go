package service

import (
	"context"
	"time"

	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/domain/entity"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/domain/repository"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/logger"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmlogrus"
)

type TimeRecordService struct {
	TimeRecordRepository repository.TimeRecordRepositoryInterface
}

func (p *TimeRecordService) RegisterTimeRecord(ctx context.Context, _time time.Time, description, employeeID, createdBy string) (*string, error) {
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

	timeRecordID, err := p.TimeRecordRepository.RegisterTimeRecord(ctx, timeRecord)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}
	log.WithField("timeRecord", timeRecord).Info("timeRecord registered")

	return timeRecordID, nil
}

func (p *TimeRecordService) ApproveTimeRecord(ctx context.Context, id, employeeID string) error {
	span, ctx := apm.StartSpan(ctx, "Approve", "time record domain service")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	timeRecord, err := p.TimeRecordRepository.FindTimeRecord(ctx, id)
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

	err = p.TimeRecordRepository.SaveTimeRecord(ctx, timeRecord)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return err
	}

	return nil
}

func (p *TimeRecordService) RefuseTimeRecord(ctx context.Context, id, refusedReason, employeeID string) error {
	span, ctx := apm.StartSpan(ctx, "Refuse", "time record domain service")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	timeRecord, err := p.TimeRecordRepository.FindTimeRecord(ctx, id)
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

	err = p.TimeRecordRepository.SaveTimeRecord(ctx, timeRecord)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return err
	}

	return nil
}

func (p *TimeRecordService) FindTimeRecord(ctx context.Context, id string) (*entity.TimeRecord, error) {
	span, ctx := apm.StartSpan(ctx, "Find", "time record domain service")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	timeRecord, err := p.TimeRecordRepository.FindTimeRecord(ctx, id)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}
	log.WithField("timeRecord", timeRecord).Info("timeRecord finded")
	return timeRecord, nil
}

func (p *TimeRecordService) SearchTimeRecords(ctx context.Context, fromDate, toDate time.Time, status int, employeeID, approvedBy, refusedBy, createdBy string, pageSize int, pageToken string) (*string, []*entity.TimeRecord, error) {
	span, ctx := apm.StartSpan(ctx, "FindAllByEmployeeID", "time record domain service")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	filter, err := entity.NewFilter(fromDate, toDate, status, employeeID, approvedBy, refusedBy, createdBy, pageSize, pageToken)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, nil, err
	}

	nextPageToken, timeRecords, err := p.TimeRecordRepository.SearchTimeRecords(ctx, filter)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, nil, err
	}
	log.WithField("timeRecords", timeRecords).Info("timeRecords finded")
	return nextPageToken, timeRecords, nil
}

func NewTimeRecordService(timeRecordRepository repository.TimeRecordRepositoryInterface) *TimeRecordService {
	return &TimeRecordService{
		TimeRecordRepository: timeRecordRepository,
	}
}
