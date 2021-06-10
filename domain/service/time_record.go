package service

import (
	"context"
	"time"

	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/domain/model"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/domain/repository"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/logger"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmlogrus"
)

type TimeRecordService struct {
	TimeRecordRepository repository.TimeRecordRepositoryInterface
}

func (p *TimeRecordService) Register(ctx context.Context, _time time.Time, description, employeeID string) (*model.TimeRecord, error) {
	span, ctx := apm.StartSpan(ctx, "Register", "time record domain service")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	timeRecord, err := model.NewTimeRecord(_time, description, employeeID)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}
	log.WithField("timeRecord", timeRecord).Info("timeRecord created")

	err = p.TimeRecordRepository.Register(ctx, timeRecord)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}
	log.WithField("timeRecord", timeRecord).Info("timeRecord registered")

	return timeRecord, nil
}

func (p *TimeRecordService) Approve(ctx context.Context, id, employeeID string) (*model.TimeRecord, error) {
	span, ctx := apm.StartSpan(ctx, "Approve", "time record domain service")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	timeRecord, err := p.TimeRecordRepository.Find(ctx, id)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}
	log.WithField("timeRecord", timeRecord).Info("timeRecord finded")

	err = timeRecord.Approve(employeeID)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}
	log.WithField("timeRecord", timeRecord).Info("timeRecord approved")

	err = p.TimeRecordRepository.Save(ctx, timeRecord)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}

	return timeRecord, nil
}

func (p *TimeRecordService) Refuse(ctx context.Context, id, refusedReason, employeeID string) (*model.TimeRecord, error) {
	span, ctx := apm.StartSpan(ctx, "Refuse", "time record domain service")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	timeRecord, err := p.TimeRecordRepository.Find(ctx, id)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}
	log.WithField("timeRecord", timeRecord).Info("timeRecord finded")

	err = timeRecord.Refuse(employeeID, refusedReason)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}
	log.WithField("timeRecord", timeRecord).Info("timeRecord refused")

	err = p.TimeRecordRepository.Save(ctx, timeRecord)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}

	return timeRecord, nil
}

func (p *TimeRecordService) Find(ctx context.Context, id string) (*model.TimeRecord, error) {
	span, ctx := apm.StartSpan(ctx, "Find", "time record domain service")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	timeRecord, err := p.TimeRecordRepository.Find(ctx, id)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}
	log.WithField("timeRecord", timeRecord).Info("timeRecord finded")
	return timeRecord, nil
}

func (p *TimeRecordService) FindAllByEmployeeID(ctx context.Context, employeeID string, fromDate, toDate time.Time) ([]*model.TimeRecord, error) {
	span, ctx := apm.StartSpan(ctx, "FindAllByEmployeeID", "time record domain service")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	timeRecords, err := p.TimeRecordRepository.FindAllByEmployeeID(ctx, employeeID, fromDate, toDate)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}
	log.WithField("timeRecords", timeRecords).Info("timeRecords finded")
	return timeRecords, nil
}

func NewTimeRecordService(timeRecordRepository repository.TimeRecordRepositoryInterface) *TimeRecordService {
	return &TimeRecordService{
		TimeRecordRepository: timeRecordRepository,
	}
}
