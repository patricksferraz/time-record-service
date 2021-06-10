package model_test

import (
	"math"
	"testing"

	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/domain/model"
	"github.com/stretchr/testify/require"
	"syreclabs.com/go/faker"
)

func TestModel_TimeRecordStatus(t *testing.T) {

	status := model.PENDING
	require.Equal(t, status.String(), model.PENDING.String())
	status = model.APPROVED
	require.Equal(t, status.String(), model.APPROVED.String())
	status = model.REFUSED
	require.Equal(t, status.String(), model.REFUSED.String())

	otherStatus := model.TimeRecordStatus(faker.RandomInt(int(model.REFUSED)+1, math.MaxInt64))
	require.Equal(t, otherStatus.String(), "")
}
