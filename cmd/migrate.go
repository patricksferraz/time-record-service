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
	"errors"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/infrastructure/db"
	_ "dev.azure.com/c4ut/TimeClock/_git/time-record-service/infrastructure/db/migrations"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/utils"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	migrate "github.com/xakep666/mongo-migrate"
)

var n int

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate [up|down]",
	Short: "A brief description of your command",
	Args: func(cmd *cobra.Command, args []string) error {

		if len(args) < 1 {
			return errors.New("requires up or down argument")
		}
		return nil

	},
	Run: func(cmd *cobra.Command, args []string) {

		ctx := context.Background()
		db, err := db.NewMongo(ctx, uri, dbName)
		if err != nil {
			log.Fatal(err)
		}
		migrate.SetDatabase(db.Database)

		switch args[0] {
		case "up":
			err = migrate.Up(n)
			if err != nil {
				log.Fatal(err)
			}
		case "down":
			err = migrate.Down(n)
			if err != nil {
				log.Fatal(err)
			}
		default:
			log.Fatal("requires up or down argument")
		}

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

	dUri := utils.GetEnv("DB_URI", "mongodb://localhost")
	dDbName := utils.GetEnv("DB_NAME", "time_record_service")

	rootCmd.AddCommand(migrateCmd)
	migrateCmd.Flags().IntVarP(&n, "n", "n", migrate.AllAvailable, "amount of migrations to UP or DOWN")
	migrateCmd.Flags().StringVarP(&uri, "uri", "u", dUri, "gRPC Server port")
	migrateCmd.Flags().StringVarP(&dbName, "dbName", "", dDbName, "gRPC Server port")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// migrateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// migrateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
