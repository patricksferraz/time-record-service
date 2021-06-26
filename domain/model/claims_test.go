package model_test

import (
	"testing"

	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/domain/model"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"syreclabs.com/go/faker"
)

func TestModel_NewEmployeeClaims(t *testing.T) {

	employeeID := uuid.NewV4().String()
	count := faker.Number().NumberInt(2)

	var roles []string
	for i := 0; i < count; i++ {
		roles = append(roles, faker.Lorem().Word())
	}

	claims, err := model.NewClaims(employeeID, roles)

	require.Nil(t, err)
	require.NotEmpty(t, uuid.FromStringOrNil(claims.EmployeeID))
	require.Equal(t, claims.Roles, roles)

	_, err = model.NewClaims("", roles)
	require.NotNil(t, err)
}
