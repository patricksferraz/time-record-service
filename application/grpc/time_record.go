package grpc

import (
	"context"

	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/application/grpc/pb"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/domain/entity"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/domain/service"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/logger"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmlogrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TimeRecordGrpcService struct {
	pb.UnimplementedTimeRecordServiceServer
	TimeRecordService *service.TimeRecordService
	AuthInterceptor   *AuthInterceptor
}

func (t *TimeRecordGrpcService) RegisterTimeRecord(ctx context.Context, in *pb.RegisterTimeRecordRequest) (*pb.RegisterTimeRecordResponse, error) {
	span, ctx := apm.StartSpan(ctx, "RegisterTimeRecord", "gRPC application")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))
	log.WithField("in", in).Info("handling RegisterTimeRecord request")

	timeRecordID, err := t.TimeRecordService.RegisterTimeRecord(ctx, in.Time.AsTime(), in.Description, in.EmployeeId, t.AuthInterceptor.Claims.EmployeeID)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return &pb.RegisterTimeRecordResponse{}, err
	}
	log.WithField("timeRecordID", *timeRecordID).Info("timeRecord registered")

	return &pb.RegisterTimeRecordResponse{
		Id: *timeRecordID,
	}, nil
}

func (t *TimeRecordGrpcService) ApproveTimeRecord(ctx context.Context, in *pb.ApproveTimeRecordRequest) (*pb.StatusResponse, error) {
	span, ctx := apm.StartSpan(ctx, "ApproveTimeRecord", "gRPC application")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))
	log.WithField("in", in).Info("handling ApproveTimeRecord request")

	err := t.TimeRecordService.ApproveTimeRecord(ctx, in.Id, t.AuthInterceptor.Claims.EmployeeID)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return &pb.StatusResponse{
			Code:    uint32(status.Code(err)),
			Message: "not updated",
			Error:   err.Error(),
		}, err
	}

	return &pb.StatusResponse{
		Code:    uint32(codes.OK),
		Message: "successfully " + entity.APPROVED.String(),
	}, nil
}

func (t *TimeRecordGrpcService) RefuseTimeRecord(ctx context.Context, in *pb.RefuseTimeRecordRequest) (*pb.StatusResponse, error) {
	span, ctx := apm.StartSpan(ctx, "RefuseTimeRecord", "gRPC application")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))
	log.WithField("in", in).Info("handling RefuseTimeRecord request")

	err := t.TimeRecordService.RefuseTimeRecord(ctx, in.Id, in.RefusedReason, t.AuthInterceptor.Claims.EmployeeID)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return &pb.StatusResponse{
			Code:    uint32(status.Code(err)),
			Message: "not updated",
			Error:   err.Error(),
		}, err
	}

	return &pb.StatusResponse{
		Code:    uint32(codes.OK),
		Message: "successfully " + entity.REFUSED.String(),
	}, nil
}

func (t *TimeRecordGrpcService) FindTimeRecord(ctx context.Context, in *pb.FindTimeRecordRequest) (*pb.FindTimeRecordResponse, error) {
	span, ctx := apm.StartSpan(ctx, "FindTimeRecord", "gRPC application")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))
	log.WithField("in", in).Info("handling FindTimeRecord request")

	timeRecord, err := t.TimeRecordService.FindTimeRecord(ctx, in.Id)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return &pb.FindTimeRecordResponse{}, err
	}
	log.WithField("timeRecord", timeRecord).Info("timeRecord finded")

	return &pb.FindTimeRecordResponse{
		TimeRecord: &pb.TimeRecord{
			Id:            timeRecord.ID,
			Time:          timestamppb.New(timeRecord.Time),
			Status:        pb.TimeRecord_Status(timeRecord.Status),
			Description:   timeRecord.Description,
			RefusedReason: timeRecord.RefusedReason,
			RegularTime:   timeRecord.RegularTime,
			TzOffset:      int64(timeRecord.TzOffset),
			EmployeeId:    timeRecord.EmployeeID,
			ApprovedBy:    timeRecord.ApprovedBy,
			RefusedBy:     timeRecord.RefusedBy,
			CreatedAt:     timestamppb.New(timeRecord.CreatedAt),
			UpdatedAt:     timestamppb.New(timeRecord.UpdatedAt),
		},
	}, nil
}

func (t *TimeRecordGrpcService) SearchTimeRecords(in *pb.SearchTimeRecordsRequest, stream pb.TimeRecordService_SearchTimeRecordsServer) error {
	span, ctx := apm.StartSpan(stream.Context(), "SearchTimeRecords", "gRPC application")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))
	log.WithField("in", in).Info("handling SearchTimeRecords request")

	timeRecords, err := t.TimeRecordService.SearchTimeRecords(ctx, in.EmployeeId, in.FromDate.AsTime(), in.ToDate.AsTime())
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return err
	}
	log.WithField("timeRecords", timeRecords).Info("timeRecords searched")

	for _, timeRecord := range timeRecords {
		stream.Send(&pb.TimeRecord{
			Id:            timeRecord.ID,
			Time:          timestamppb.New(timeRecord.Time),
			Status:        pb.TimeRecord_Status(timeRecord.Status),
			Description:   timeRecord.Description,
			RefusedReason: timeRecord.RefusedReason,
			RegularTime:   timeRecord.RegularTime,
			TzOffset:      int64(timeRecord.TzOffset),
			EmployeeId:    timeRecord.EmployeeID,
			ApprovedBy:    timeRecord.ApprovedBy,
			RefusedBy:     timeRecord.RefusedBy,
			CreatedAt:     timestamppb.New(timeRecord.CreatedAt),
			UpdatedAt:     timestamppb.New(timeRecord.UpdatedAt),
		})
	}

	return nil
}

func NewTimeRecordGrpcService(service *service.TimeRecordService, authInterceptor *AuthInterceptor) *TimeRecordGrpcService {
	return &TimeRecordGrpcService{
		TimeRecordService: service,
		AuthInterceptor:   authInterceptor,
	}
}
