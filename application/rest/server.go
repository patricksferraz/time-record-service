package rest

import (
	"fmt"
	"log"

	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/application/grpc/pb"
	_ "dev.azure.com/c4ut/TimeClock/_git/time-record-service/application/rest/docs"
	_service "dev.azure.com/c4ut/TimeClock/_git/time-record-service/domain/service"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/infrastructure/db"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/infrastructure/repository"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.elastic.co/apm/module/apmgin"
)

// @title Time Record Swagger API
// @version 1.0
// @description Swagger API for Golang Project Time Record.
// @termsOfService http://swagger.io/terms/

// @contact.name Coding4u
// @contact.email contato@coding4u.com.br

// @BasePath /api/v1
// @query.collection.format multi

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func StartRestServer(database *db.Mongo, service pb.AuthServiceClient, port int) {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.Default())
	r.Use(apmgin.Middleware(r))

	authService := _service.NewAuthService(service)
	authMiddlerare := NewAuthMiddleware(authService)
	timeRecordRepository := repository.NewTimeRecordRepository(database)
	timeRecordService := _service.NewTimeRecordService(timeRecordRepository)
	timeRecordRestService := NewTimeRecordRestService(timeRecordService, authMiddlerare)

	v1 := r.Group("api/v1/time-records")
	{
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		authorized := v1.Group("/", authMiddlerare.Require())
		{
			authorized.POST("/", timeRecordRestService.RegisterTimeRecord)
			authorized.POST("/:id/approve", timeRecordRestService.ApproveTimeRecord)
			authorized.POST("/:id/refuse", timeRecordRestService.RefuseTimeRecord)

			authorized.GET("/", timeRecordRestService.SearchTimeRecords)
			authorized.GET("/:id", timeRecordRestService.FindTimeRecord)
		}
	}

	addr := fmt.Sprintf("0.0.0.0:%d", port)
	err := r.Run(addr)
	if err != nil {
		log.Fatal("cannot start rest server", err)
	}

	log.Printf("rest server has been started on port %d", port)
}
