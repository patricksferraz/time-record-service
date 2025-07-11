package grpc

import (
	"context"

	"github.com/patricksferraz/time-record-service/domain/entity"
	"github.com/patricksferraz/time-record-service/domain/service"
	"github.com/patricksferraz/time-record-service/logger"
	"github.com/patricksferraz/time-record-service/proto/pb"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmlogrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GrpcService struct {
	pb.UnimplementedTimeRecordServiceServer
	Service         *service.Service
	AuthInterceptor *AuthInterceptor
}

func NewGrpcService(service *service.Service, authInterceptor *AuthInterceptor) *GrpcService {
	return &GrpcService{
		Service:         service,
		AuthInterceptor: authInterceptor,
	}
}

func (t *GrpcService) RegisterTimeRecord(ctx context.Context, in *pb.RegisterTimeRecordRequest) (*pb.RegisterTimeRecordResponse, error) {
	span, ctx := apm.StartSpan(ctx, "RegisterTimeRecord", "gRPC application")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))
	log.WithField("in", in).Info("handling RegisterTimeRecord request")

	timeRecordID, err := t.Service.RegisterTimeRecord(ctx, in.Time.AsTime(), in.Description, in.EmployeeId, in.CompanyId, t.AuthInterceptor.Claims.EmployeeID)
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

func (t *GrpcService) ApproveTimeRecord(ctx context.Context, in *pb.ApproveTimeRecordRequest) (*pb.StatusResponse, error) {
	span, ctx := apm.StartSpan(ctx, "ApproveTimeRecord", "gRPC application")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))
	log.WithField("in", in).Info("handling ApproveTimeRecord request")

	err := t.Service.ApproveTimeRecord(ctx, in.Id, t.AuthInterceptor.Claims.EmployeeID)
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

func (t *GrpcService) RefuseTimeRecord(ctx context.Context, in *pb.RefuseTimeRecordRequest) (*pb.StatusResponse, error) {
	span, ctx := apm.StartSpan(ctx, "RefuseTimeRecord", "gRPC application")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))
	log.WithField("in", in).Info("handling RefuseTimeRecord request")

	err := t.Service.RefuseTimeRecord(ctx, in.Id, in.RefusedReason, t.AuthInterceptor.Claims.EmployeeID)
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

func (t *GrpcService) FindTimeRecord(ctx context.Context, in *pb.FindTimeRecordRequest) (*pb.FindTimeRecordResponse, error) {
	span, ctx := apm.StartSpan(ctx, "FindTimeRecord", "gRPC application")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))
	log.WithField("in", in).Info("handling FindTimeRecord request")

	timeRecord, err := t.Service.FindTimeRecord(ctx, in.Id)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return &pb.FindTimeRecordResponse{}, err
	}

	return &pb.FindTimeRecordResponse{
		TimeRecord: &pb.TimeRecord{
			Id:            timeRecord.ID,
			Time:          timestamppb.New(timeRecord.Time),
			Status:        pb.TimeRecord_Status(timeRecord.Status),
			Description:   timeRecord.Description,
			RefusedReason: timeRecord.RefusedReason,
			RegularTime:   timeRecord.RegularTime,
			TzOffset:      int32(timeRecord.TzOffset),
			EmployeeId:    *timeRecord.EmployeeID,
			ApprovedBy:    *timeRecord.ApprovedBy,
			RefusedBy:     *timeRecord.RefusedBy,
			CreatedAt:     timestamppb.New(timeRecord.CreatedAt),
			UpdatedAt:     timestamppb.New(timeRecord.UpdatedAt),
		},
	}, nil
}

func (t *GrpcService) SearchTimeRecords(ctx context.Context, in *pb.SearchTimeRecordsRequest) (*pb.SearchTimeRecordsResponse, error) {
	span, ctx := apm.StartSpan(ctx, "SearchTimeRecords", "gRPC application")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))
	log.WithField("in", in).Info("handling SearchTimeRecords request")

	nextPageToken, timeRecords, err := t.Service.SearchTimeRecords(ctx, in.Filter.FromDate.AsTime(), in.Filter.ToDate.AsTime(), int(in.Filter.Status), in.Filter.EmployeeId, in.Filter.ApprovedBy, in.Filter.RefusedBy, in.Filter.CreatedBy, in.Filter.CompanyId, int(in.Filter.PageSize), in.Filter.PageToken)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return &pb.SearchTimeRecordsResponse{}, err
	}

	var result []*pb.TimeRecord
	for _, timeRecord := range timeRecords {
		result = append(
			result,
			&pb.TimeRecord{
				Id:            timeRecord.ID,
				Time:          timestamppb.New(timeRecord.Time),
				Status:        pb.TimeRecord_Status(timeRecord.Status),
				Description:   timeRecord.Description,
				RefusedReason: timeRecord.RefusedReason,
				RegularTime:   timeRecord.RegularTime,
				TzOffset:      int32(timeRecord.TzOffset),
				EmployeeId:    *timeRecord.EmployeeID,
				ApprovedBy:    *timeRecord.ApprovedBy,
				RefusedBy:     *timeRecord.RefusedBy,
				CompanyId:     *timeRecord.CompanyID,
				CreatedAt:     timestamppb.New(timeRecord.CreatedAt),
				UpdatedAt:     timestamppb.New(timeRecord.UpdatedAt),
			},
		)
	}

	return &pb.SearchTimeRecordsResponse{NextPageToken: *nextPageToken, TimeRecords: result}, nil
}

func (t *GrpcService) ExportTimeRecords(ctx context.Context, in *pb.ExportTimeRecordsRequest) (*pb.ExportTimeRecordsResponse, error) {
	span, ctx := apm.StartSpan(ctx, "ExportTimeRecords", "gRPC application")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))
	log.WithField("in", in).Info("handling ExportTimeRecords request")

	nextPageToken, registers, err := t.Service.ExportTimeRecords(ctx, in.Filter.FromDate.AsTime(), in.Filter.ToDate.AsTime(), int(in.Filter.Status), in.Filter.EmployeeId, in.Filter.ApprovedBy, in.Filter.RefusedBy, in.Filter.CreatedBy, in.Filter.CompanyId, int(in.Filter.PageSize), in.Filter.PageToken, *t.AuthInterceptor.AccessToken)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return &pb.ExportTimeRecordsResponse{}, err
	}

	var result []string
	for _, r := range registers {
		result = append(result, string(*r))
	}

	return &pb.ExportTimeRecordsResponse{NextPageToken: *nextPageToken, Registers: result}, nil
}
