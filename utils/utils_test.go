package utils_test

import (
	"os"
	"testing"
	"time"

	"github.com/patricksferraz/time-record-service/utils"
	"github.com/stretchr/testify/require"
	"syreclabs.com/go/faker"
)

func TestUtils_DateEqual(t *testing.T) {

	today := time.Now()

	result := utils.DateEqual(today, today.Add(time.Second))
	require.True(t, result)
	result = utils.DateEqual(today, today.AddDate(0, 0, 1))
	require.False(t, result)
}

func TestUtils_GetEnv(t *testing.T) {

	key := faker.Lorem().Word()
	defaultVal := faker.Lorem().Word()

	result := utils.GetEnv(key, defaultVal)
	require.Equal(t, result, defaultVal)

	otherVal := faker.Lorem().Word()
	os.Setenv(key, otherVal)

	result = utils.GetEnv(key, defaultVal)
	require.Equal(t, result, otherVal)
}
