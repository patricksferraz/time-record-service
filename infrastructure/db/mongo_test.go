package db_test

import (
	"context"
	"testing"
	"time"

	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/infrastructure/db"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/utils"
	"github.com/stretchr/testify/require"
)

func TestDb_NewMongo(t *testing.T) {

	uri := utils.GetEnv("DB_URI", "mongodb://localhost")
	dbName := utils.GetEnv("DB_NAME", "test")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	database, err := db.NewMongo(ctx, uri, dbName)
	require.Nil(t, err)
	require.NotEmpty(t, database)

	err = database.Test(ctx)
	require.Nil(t, err)

	err = database.Close(ctx)
	require.Nil(t, err)

	uri = "mongodb://1.1.1.1"
	database, _ = db.NewMongo(ctx, uri, dbName)
	err = database.Test(ctx)
	require.NotNil(t, err)
}
