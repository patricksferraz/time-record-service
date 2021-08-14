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

	"github.com/c-4u/time-record-service/application/kafka"
	"github.com/c-4u/time-record-service/infrastructure/db"
	"github.com/c-4u/time-record-service/infrastructure/external"
	"github.com/c-4u/time-record-service/infrastructure/external/topic"
	"github.com/c-4u/time-record-service/utils"
	"github.com/spf13/cobra"
)

// kafkaCmd represents the kafka command
func kafkaCmd() *cobra.Command {
	var servers string
	var groupId string
	var dsn string
	var dsnType string

	kafkaCmd := &cobra.Command{
		Use:   "kafka",
		Short: "Run kafka Service",

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

			k, err := external.NewKafka(servers, groupId, []string{topic.NEW_EMPLOYEE, topic.NEW_COMPANY})
			if err != nil {
				log.Fatal("cannot start kafka processor", err)
			}
			kafka.StartKafkaProcessor(database, servers, groupId, k)
		},
	}

	dDsn := os.Getenv("DSN")
	sDsnType := os.Getenv("DSN_TYPE")
	dServers := utils.GetEnv("KAFKA_BOOTSTRAP_SERVERS", "kafka:9094")
	dGroupId := utils.GetEnv("KAFKA_CONSUMER_GROUP_ID", "time-record-service")

	kafkaCmd.Flags().StringVarP(&dsn, "dsn", "d", dDsn, "dsn")
	kafkaCmd.Flags().StringVarP(&dsnType, "dsnType", "t", sDsnType, "dsn type")
	kafkaCmd.Flags().StringVarP(&servers, "servers", "s", dServers, "kafka servers")
	kafkaCmd.Flags().StringVarP(&groupId, "groupId", "i", dGroupId, "kafka group id")

	return kafkaCmd
}

func init() {
	rootCmd.AddCommand(kafkaCmd())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// kafkaCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// kafkaCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
