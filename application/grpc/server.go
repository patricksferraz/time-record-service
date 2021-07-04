package grpc

import (
	"fmt"
	"log"
	"net"

	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/application/grpc/pb"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/domain/service"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/infrastructure/db"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/infrastructure/repository"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.elastic.co/apm/module/apmgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartGrpcServer(database *db.Mongo, authConn *grpc.ClientConn, employeeConn *grpc.ClientConn, port int) {

	authService := service.NewAuthService(authConn)
	interceptor := NewAuthInterceptor(authService)
	timeRecordRepository := repository.NewTimeRecordRepository(database)
	employeeRepository := service.NewEmployeeService(employeeConn)
	timeRecordService := service.NewTimeRecordService(timeRecordRepository, employeeRepository)
	timeRecordGrpcService := NewTimeRecordGrpcService(timeRecordService, interceptor)

	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			apmgrpc.NewUnaryServerInterceptor(apmgrpc.WithRecovery()),
			interceptor.Unary(),
		),
		grpc.StreamInterceptor(interceptor.Stream()),
	)

	reflection.Register(grpcServer)
	pb.RegisterTimeRecordServiceServer(grpcServer, timeRecordGrpcService)

	address := fmt.Sprintf("0.0.0.0:%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start grpc server", err)
	}

	log.Printf("gRPC server has been started on port %d", port)

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start grpc server", err)
	}
}
