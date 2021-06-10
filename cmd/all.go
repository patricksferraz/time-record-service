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
	"context"
	"log"
	"os"

	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/application/grpc"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/application/grpc/pb"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/application/rest"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/infrastructure/db"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/infrastructure/external"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/utils"
	"github.com/spf13/cobra"
)

// allCmd represents the all command
var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Run both gRPC and rest servers",

	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		database, err := db.NewMongo(ctx, uri, dbName)
		if err != nil {
			log.Fatal(err)
		}

		err = database.Test(ctx)
		if err != nil {
			log.Fatal(err)
		}

		if utils.GetEnv("DB_MIGRATE", "false") == "true" {
			err = database.Migrate()
			if err != nil {
				log.Fatal(err)
			}
		}

		defer database.Close(ctx)

		authServiceAddr := os.Getenv("AUTH_SERVICE_ADDR")
		conn, err := external.ConnectAuthService(authServiceAddr)
		if err != nil {
			log.Fatal(err)
		}

		defer conn.Close()
		authdb := pb.NewAuthServiceClient(conn)

		go rest.StartRestServer(database, authdb, restPort)
		grpc.StartGrpcServer(database, authdb, grpcPort)
	},
}

func init() {
	rootCmd.AddCommand(allCmd)
	allCmd.Flags().IntVarP(&grpcPort, "grpcPort", "g", 50051, "gRPC Server port")
	allCmd.Flags().IntVarP(&restPort, "restPort", "r", 8080, "rest server port")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// allCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// allCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
