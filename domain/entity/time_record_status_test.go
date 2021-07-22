package entity_test

import (
	"math"
	"testing"

	"github.com/c-4u/time-record-service/domain/entity"
	"github.com/stretchr/testify/require"
	"syreclabs.com/go/faker"
)

func TestModel_TimeRecordStatus(t *testing.T) {

	status := entity.PENDING
	require.Equal(t, status.String(), entity.PENDING.String())
	status = entity.APPROVED
	require.Equal(t, status.String(), entity.APPROVED.String())
	status = entity.REFUSED
	require.Equal(t, status.String(), entity.REFUSED.String())

	otherStatus := entity.TimeRecordStatus(faker.RandomInt(int(entity.REFUSED)+1, math.MaxInt64))
	require.Equal(t, otherStatus.String(), "")
}
