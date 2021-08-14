package kafka

import (
	"context"
	"fmt"

	"github.com/c-4u/time-record-service/application/kafka/schema"
	"github.com/c-4u/time-record-service/domain/service"
	"github.com/c-4u/time-record-service/infrastructure/external"
	"github.com/c-4u/time-record-service/infrastructure/external/topic"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaProcessor struct {
	Service *service.Service
	K       *external.Kafka
}

func NewKafkaProcessor(service *service.Service, kafka *external.Kafka) *KafkaProcessor {
	return &KafkaProcessor{
		Service: service,
		K:       kafka,
	}
}

func (p *KafkaProcessor) Consume() {
	p.K.Consumer.SubscribeTopics(p.K.Topics, nil)
	for {
		msg, err := p.K.Consumer.ReadMessage(-1)
		if err == nil {
			// fmt.Println(string(msg.Value))
			p.processMessage(msg)
		}
	}
}

func (p *KafkaProcessor) processMessage(msg *ckafka.Message) {
	switch _topic := *msg.TopicPartition.Topic; _topic {
	case topic.NEW_EMPLOYEE:
		// TODO: add fault tolerance
		err := p.createEmployee(msg)
		if err != nil {
			fmt.Println("creation error ", err)
		}
	case topic.NEW_COMPANY:
		err := p.createCompany(msg)
		if err != nil {
			fmt.Println("creation error ", err)
		}
	default:
		fmt.Println("not a valid topic", string(msg.Value))
	}
}

func (p *KafkaProcessor) createEmployee(msg *ckafka.Message) error {
	employeeEvent := &schema.EmployeeEvent{}
	err := employeeEvent.ParseJson(msg.Value)
	if err != nil {
		return err
	}

	err = p.Service.CreateEmployee(context.TODO(), employeeEvent.Employee.ID, employeeEvent.Employee.Pis, employeeEvent.Employee.CompanyID)
	if err != nil {
		return err
	}

	return nil
}

func (p *KafkaProcessor) createCompany(msg *ckafka.Message) error {
	companyEvent := &schema.CompanyEvent{}
	err := companyEvent.ParseJson(msg.Value)
	if err != nil {
		return err
	}

	err = p.Service.CreateCompany(context.TODO(), companyEvent.Company.ID)
	if err != nil {
		return err
	}

	return nil
}
