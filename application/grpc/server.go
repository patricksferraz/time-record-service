package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/c-4u/time-record-service/application/grpc/pb"
	"github.com/c-4u/time-record-service/domain/service"
	"github.com/c-4u/time-record-service/infrastructure/db"
	"github.com/c-4u/time-record-service/infrastructure/repository"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.elastic.co/apm/module/apmgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartGrpcServer(database *db.Postgres, authConn *grpc.ClientConn, port int) {

	authService := service.NewAuthService(authConn)
	interceptor := NewAuthInterceptor(authService)
	repository := repository.NewPostgresRepository(database)
	service := service.NewService(repository)
	grpcService := NewGrpcService(service, interceptor)

	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			apmgrpc.NewUnaryServerInterceptor(apmgrpc.WithRecovery()),
			interceptor.Unary(),
		),
		grpc.StreamInterceptor(interceptor.Stream()),
	)

	reflection.Register(grpcServer)
	pb.RegisterTimeRecordServiceServer(grpcServer, grpcService)

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
