package kafka

import (
	"fmt"

	"github.com/c-4u/time-record-service/domain/service"
	"github.com/c-4u/time-record-service/infrastructure/db"
	"github.com/c-4u/time-record-service/infrastructure/external"
	"github.com/c-4u/time-record-service/infrastructure/repository"
)

func StartKafkaProcessor(database *db.Postgres, servers, groupId string, kafka *external.Kafka) {
	repository := repository.NewPostgresRepository(database)
	service := service.NewService(repository)

	fmt.Println("kafka consumer has been started")
	processor := NewKafkaProcessor(service, kafka)
	processor.Consume()
}
