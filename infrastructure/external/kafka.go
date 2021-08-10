package external

import (
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

type Kafka struct {
	Consumer *ckafka.Consumer
	Topics   []string
}

func NewKafka(servers, groupId string, topics []string) (*Kafka, error) {
	c, err := ckafka.NewConsumer(
		&ckafka.ConfigMap{
			"bootstrap.servers": servers,
			"group.id":          groupId,
			"auto.offset.reset": "earliest",
		},
	)
	if err != nil {
		return nil, err
	}

	return &Kafka{
		Consumer: c,
		Topics:   topics,
	}, nil
}
