package entity_test

import (
	"testing"

	"github.com/c-4u/time-record-service/domain/entity"
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

	claims, err := entity.NewClaims(employeeID, roles)

	require.Nil(t, err)
	require.NotEmpty(t, uuid.FromStringOrNil(claims.EmployeeID))
	require.Equal(t, claims.Roles, roles)
}
