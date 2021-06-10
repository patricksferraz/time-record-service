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
	"path/filepath"
	"runtime"

	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/application/grpc/pb"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/application/rest"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/infrastructure/db"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/infrastructure/external"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/utils"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var restPort int

// restCmd represents the rest command
var restCmd = &cobra.Command{
	Use:   "rest",
	Short: "Run rest Service",

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

		rest.StartRestServer(database, authdb, restPort)
	},
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

	rootCmd.AddCommand(restCmd)
	restCmd.Flags().IntVarP(&restPort, "port", "p", 8080, "rest server port")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// restCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// restCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
