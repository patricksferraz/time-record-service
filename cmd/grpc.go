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
	"github.com/patricksferraz/time-record-service/infrastructure/db"
	"github.com/patricksferraz/time-record-service/infrastructure/external"
	"github.com/patricksferraz/time-record-service/utils"
	"github.com/spf13/cobra"
)

// grpcCmd represents the grpc command
func grpcCmd() *cobra.Command {
	var grpcPort int
	var servers string
	var groupId string
	var dsn string
	var dsnType string

	grpcCmd := &cobra.Command{
		Use:   "grpc",
		Short: "Run gRPC Service",

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
			k, err := external.NewKafkaProducer(servers, deliveryChan)
			if err != nil {
				log.Fatal("cannot start kafka processor", err)
			}

			grpc.StartGrpcServer(database, authConn, k, grpcPort)
		},
	}

	dDsn := os.Getenv("DSN")
	sDsnType := os.Getenv("DSN_TYPE")
	dServers := utils.GetEnv("KAFKA_BOOTSTRAP_SERVERS", "kafka:9094")
	dGroupId := utils.GetEnv("KAFKA_CONSUMER_GROUP_ID", "time-record-service")

	grpcCmd.Flags().StringVarP(&dsn, "dsn", "d", dDsn, "dsn")
	grpcCmd.Flags().StringVarP(&dsnType, "dsnType", "t", sDsnType, "dsn type")
	grpcCmd.Flags().IntVarP(&grpcPort, "port", "p", 50051, "gRPC Server port")
	grpcCmd.Flags().StringVarP(&servers, "servers", "s", dServers, "kafka servers")
	grpcCmd.Flags().StringVarP(&groupId, "groupId", "i", dGroupId, "kafka group id")

	return grpcCmd
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

	rootCmd.AddCommand(grpcCmd())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// grpcCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// grpcCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
