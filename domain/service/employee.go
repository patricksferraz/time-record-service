package service

import (
	"context"

	"github.com/c-4u/time-record-service/application/grpc/pb"
	"github.com/c-4u/time-record-service/domain/entity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type EmployeeService struct {
	service pb.EmployeeServiceClient
}

func NewEmployeeService(cc *grpc.ClientConn) *EmployeeService {
	return &EmployeeService{
		service: pb.NewEmployeeServiceClient(cc),
	}
}

func (s *EmployeeService) FindEmployee(ctx context.Context, id, accessToken string) (*entity.Employee, error) {
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", accessToken)

	req := &pb.FindEmployeeRequest{
		EmployeeId: id,
	}

	res, err := s.service.FindEmployee(ctx, req)
	if err != nil {
		return nil, err
	}

	employee := &entity.Employee{
		Pis: res.Pis,
	}
	employee.ID = res.Id

	return employee, nil
}
