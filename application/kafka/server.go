package kafka

import (
	"fmt"

	"github.com/c-4u/time-record-service/domain/service"
	"github.com/c-4u/time-record-service/infrastructure/db"
	"github.com/c-4u/time-record-service/infrastructure/external"
	"github.com/c-4u/time-record-service/infrastructure/repository"
)

func StartKafkaServer(database *db.Postgres, kafkaProducer *external.KafkaProducer, kafkaConsumer *external.KafkaConsumer) {
	repository := repository.NewRepository(database, kafkaProducer)
	service := service.NewService(repository)

	fmt.Println("kafka pocessor has been started")
	processor := NewKafkaProcessor(service, kafkaConsumer)
	processor.Consume()
}
