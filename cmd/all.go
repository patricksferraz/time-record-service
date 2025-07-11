/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/joho/godotenv"
	"github.com/patricksferraz/time-record-service/application/grpc"
	"github.com/patricksferraz/time-record-service/application/kafka"
	"github.com/patricksferraz/time-record-service/application/rest"
	"github.com/patricksferraz/time-record-service/infrastructure/db"
	"github.com/patricksferraz/time-record-service/infrastructure/external"
	"github.com/patricksferraz/time-record-service/infrastructure/external/topic"
	"github.com/patricksferraz/time-record-service/utils"
	"github.com/spf13/cobra"
)

// allCmd represents the all command
func allCmd() *cobra.Command {
	var grpcPort int
	var restPort int
	var servers string
	var groupId string
	var dsn string
	var dsnType string

	allCmd := &cobra.Command{
		Use:   "all",
		Short: "Run both gRPC and rest servers",

		Run: func(cmd *cobra.Command, args []string) {
			database, err := db.NewPostgres(dsnType, dsn)
			if err != nil {
				log.Fatal(err)
			}

			if utils.GetEnv("DB_DEBUG", "false") == "true" {
				database.Debug(true)
			}

			if utils.GetEnv("DB_MIGRATE", "false") == "true" {
				database.Migrate()
			}
			defer database.Db.Close()

			authServiceAddr := os.Getenv("AUTH_SERVICE_ADDR")
			authConn, err := external.GrpcClient(authServiceAddr)
			if err != nil {
				log.Fatal(err)
			}
			defer authConn.Close()

			deliveryChan := make(chan ckafka.Event)
			kc, err := external.NewKafkaConsumer(servers, groupId, topic.CONSUMER_TOPICS)
			if err != nil {
				log.Fatal("cannot start kafka consumer", err)
			}

			kp, err := external.NewKafkaProducer(servers, deliveryChan)
			if err != nil {
				log.Fatal("cannot start kafka producer", err)
			}

			go kafka.StartKafkaServer(database, kp, kc)
			go rest.StartRestServer(database, authConn, kp, restPort)
			grpc.StartGrpcServer(database, authConn, kp, grpcPort)
		},
	}

	dDsn := os.Getenv("DSN")
	sDsnType := os.Getenv("DSN_TYPE")
	dServers := utils.GetEnv("KAFKA_BOOTSTRAP_SERVERS", "kafka:9094")
	dGroupId := utils.GetEnv("KAFKA_CONSUMER_GROUP_ID", "time-record-service")

	allCmd.Flags().StringVarP(&dsn, "dsn", "d", dDsn, "dsn")
	allCmd.Flags().StringVarP(&dsnType, "dsnType", "t", sDsnType, "dsn type")
	allCmd.Flags().StringVarP(&servers, "servers", "s", dServers, "kafka servers")
	allCmd.Flags().StringVarP(&groupId, "groupId", "i", dGroupId, "kafka group id")
	allCmd.Flags().IntVarP(&grpcPort, "grpcPort", "g", 50051, "gRPC Server port")
	allCmd.Flags().IntVarP(&restPort, "restPort", "r", 8080, "rest server port")

	return allCmd
}

func init() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	if os.Getenv("ENV") == "dev" {
		err := godotenv.Load(basepath + "/../.env")
		if err != nil {
			log.Printf("Error loading .env files")
		}
	}

	rootCmd.AddCommand(allCmd())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// allCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// allCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
