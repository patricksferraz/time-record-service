package entity_test

import (
	"testing"
	"time"

	"github.com/c-4u/time-record-service/domain/entity"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"syreclabs.com/go/faker"
)

func TestModel_NewTimeRecord(t *testing.T) {

	now := time.Now()
	_, offset := now.Zone()
	description := faker.Lorem().Sentence(10)
	employeeID := uuid.NewV4().String()
	timeRecord, err := entity.NewTimeRecord(now, description, employeeID, employeeID)

	require.Nil(t, err)
	require.NotEmpty(t, uuid.FromStringOrNil(timeRecord.ID))
	require.Equal(t, timeRecord.Time, now)
	require.Equal(t, timeRecord.Status, entity.APPROVED)
	require.Equal(t, timeRecord.Description, description)
	require.Equal(t, timeRecord.RegularTime, true)
	require.Equal(t, timeRecord.TzOffset, offset)
	require.Equal(t, timeRecord.EmployeeID, employeeID)

	yesterday := now.AddDate(0, 0, -1)
	_, offset = yesterday.Zone()
	timeRecord, err = entity.NewTimeRecord(yesterday, description, employeeID, employeeID)

	require.Nil(t, err)
	require.NotEmpty(t, uuid.FromStringOrNil(timeRecord.ID))
	require.Equal(t, timeRecord.Time, yesterday)
	require.Equal(t, timeRecord.Status, entity.PENDING)
	require.Equal(t, timeRecord.Description, description)
	require.Equal(t, timeRecord.RegularTime, false)
	require.Equal(t, timeRecord.TzOffset, offset)
	require.Equal(t, timeRecord.EmployeeID, employeeID)

	tomorrow := now.AddDate(0, 0, 1)
	_, err = entity.NewTimeRecord(tomorrow, description, employeeID, employeeID)
	require.NotNil(t, err)
	_, err = entity.NewTimeRecord(time.Time{}, "", employeeID, employeeID)
	require.NotNil(t, err)
	_, err = entity.NewTimeRecord(tomorrow, "", employeeID, employeeID)
	require.NotNil(t, err)
	_, err = entity.NewTimeRecord(tomorrow, description, "", employeeID)
	require.NotNil(t, err)
	_, err = entity.NewTimeRecord(tomorrow, description, employeeID, "")
	require.NotNil(t, err)
}

func TestModel_ChangeStatusOfATimeRecord(t *testing.T) {

	yesterday := time.Now().AddDate(0, 0, -1)
	description := faker.Lorem().Sentence(10)
	employeeID := uuid.NewV4().String()
	timeRecord, _ := entity.NewTimeRecord(yesterday, description, employeeID, employeeID)

	auditedBy := uuid.NewV4().String()
	err := timeRecord.Approve(auditedBy)
	require.Nil(t, err)
	require.Equal(t, timeRecord.Status, entity.APPROVED)
	require.True(t, timeRecord.CreatedAt.Before(timeRecord.UpdatedAt))

	err = timeRecord.Approve(auditedBy)
	require.NotNil(t, err)

	err = timeRecord.Refuse(auditedBy, description)
	require.NotNil(t, err)

	timeRecord, _ = entity.NewTimeRecord(yesterday, description, employeeID, employeeID)

	err = timeRecord.Refuse(auditedBy, description)
	require.Nil(t, err)
	require.Equal(t, timeRecord.Status, entity.REFUSED)
	require.Equal(t, timeRecord.RefusedReason, description)
	require.True(t, timeRecord.CreatedAt.Before(timeRecord.UpdatedAt))

	err = timeRecord.Refuse(auditedBy, description)
	require.NotNil(t, err)

	err = timeRecord.Approve(auditedBy)
	require.NotNil(t, err)

	timeRecord, _ = entity.NewTimeRecord(yesterday, description, employeeID, employeeID)

	err = timeRecord.Refuse(auditedBy, "")
	require.NotNil(t, err)

	timeRecord, _ = entity.NewTimeRecord(yesterday, description, employeeID, employeeID)

	err = timeRecord.Refuse("", description)
	require.NotNil(t, err)

	timeRecord, _ = entity.NewTimeRecord(yesterday, description, employeeID, employeeID)

	err = timeRecord.Approve(employeeID)
	require.NotNil(t, err)

	timeRecord, _ = entity.NewTimeRecord(yesterday, description, employeeID, employeeID)

	err = timeRecord.Approve("")
	require.NotNil(t, err)

	timeRecord, _ = entity.NewTimeRecord(yesterday, description, employeeID, employeeID)

	err = timeRecord.Refuse(employeeID, description)
	require.NotNil(t, err)
}
