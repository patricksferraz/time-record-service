package rest

import (
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/patricksferraz/time-record-service/application/rest/docs"
	_service "github.com/patricksferraz/time-record-service/domain/service"
	"github.com/patricksferraz/time-record-service/infrastructure/db"
	"github.com/patricksferraz/time-record-service/infrastructure/external"
	"github.com/patricksferraz/time-record-service/infrastructure/repository"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.elastic.co/apm/module/apmgin"
	"google.golang.org/grpc"
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
func StartRestServer(database *db.Postgres, authConn *grpc.ClientConn, kafkaProducer *external.KafkaProducer, port int) {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"POST", "OPTIONS", "GET", "PUT"},
		AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Accept", "Origin", "Cache-Control", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowAllOrigins:  true,
		AllowCredentials: true,
	}))
	r.Use(apmgin.Middleware(r))

	authClient := external.NewAuthClient(authConn)
	authMiddlerare := NewAuthMiddleware(authClient)
	repository := repository.NewRepository(database, kafkaProducer)
	service := _service.NewService(repository)
	restService := NewRestService(service, authMiddlerare)

	v1 := r.Group("api/v1/time-records")
	{
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		authorized := v1.Group("", authMiddlerare.Require())
		{
			authorized.POST("", restService.RegisterTimeRecord)
			authorized.POST("/:id/approve", restService.ApproveTimeRecord)
			authorized.POST("/:id/refuse", restService.RefuseTimeRecord)

			authorized.GET("", restService.SearchTimeRecords)
			authorized.GET("/:id", restService.FindTimeRecord)
			authorized.GET("/export", restService.ExportTimeRecords)
		}
	}

	addr := fmt.Sprintf("0.0.0.0:%d", port)
	err := r.Run(addr)
	if err != nil {
		log.Fatal("cannot start rest server", err)
	}

	log.Printf("rest server has been started on port %d", port)
}
