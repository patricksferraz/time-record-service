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

	"github.com/c-4u/time-record-service/application/rest"
	"github.com/c-4u/time-record-service/infrastructure/db"
	"github.com/c-4u/time-record-service/infrastructure/external"
	"github.com/c-4u/time-record-service/utils"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// restCmd represents the rest command
func restCmd() *cobra.Command {
	var restPort int
	// var uri string
	// var dbName string
	var dsn string
	var dsnType string

	restCmd := &cobra.Command{
		Use:   "rest",
		Short: "Run rest Service",

		Run: func(cmd *cobra.Command, args []string) {
			// ctx := context.Background()
			// database, err := db.NewMongo(ctx, uri, dbName)
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

			employeeServiceAddr := os.Getenv("EMPLOYEE_SERVICE_ADDR")
			employeeConn, err := external.GrpcClient(employeeServiceAddr)
			if err != nil {
				log.Fatal(err)
			}
			defer employeeConn.Close()

			rest.StartRestServer(database, authConn, employeeConn, restPort)
		},
	}

	// dUri := utils.GetEnv("DB_URI", "mongodb://localhost")
	// dDbName := utils.GetEnv("DB_NAME", "time_record_service")
	dDsn := utils.GetEnv("DSN", "dbname=time-record-service sslmode=disable user=postgres password=root host=trdb")
	dDsnType := utils.GetEnv("DSN_TYPE", "postgres")

	restCmd.Flags().IntVarP(&restPort, "port", "p", 8080, "rest server port")
	restCmd.Flags().StringVarP(&dsn, "dsn", "d", dDsn, "dsn")
	restCmd.Flags().StringVarP(&dsnType, "dsnType", "t", dDsnType, "dsn type")
	// restCmd.Flags().StringVarP(&uri, "uri", "u", dUri, "database uri")
	// restCmd.Flags().StringVarP(&dbName, "dbName", "", dDbName, "database name")

	return restCmd
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

	rootCmd.AddCommand(restCmd())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// restCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// restCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
