package entity

import (
	"errors"
	"time"

	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/utils"
	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type TimeRecord struct {
	Base          `bson:",inline" valid:"-"`
	Time          time.Time        `json:"time,omitempty" bson:"time" valid:"required"`
	Status        TimeRecordStatus `json:"status" bson:"status" valid:"timeRecordStatus,optional"`
	Description   string           `json:"description,omitempty" bson:"description,omitempty" valid:"-"`
	RefusedReason string           `json:"refused_reason,omitempty" bson:"refused_reason,omitempty" valid:"-"`
	RegularTime   bool             `json:"regular_time" bson:"regular_time" valid:"-"`
	TzOffset      int              `json:"tz_offset" bson:"tz_offset" valid:"int"`
	EmployeeID    string           `json:"employee_id,omitempty" bson:"employee_id" valid:"uuid"`
	ApprovedBy    string           `json:"approved_by,omitempty" bson:"approved_by,omitempty" valid:"-"`
	RefusedBy     string           `json:"refused_by,omitempty" bson:"refused_by,omitempty" valid:"-"`
	CreatedBy     string           `json:"created_by,omitempty" bson:"created_by" valid:"uuid"`
}

func NewTimeRecord(_time time.Time, description, employeeID, createdBy string) (*TimeRecord, error) {

	_, offset := _time.Zone()
	timeRecord := TimeRecord{
		Time:        _time,
		Status:      APPROVED,
		Description: description,
		RegularTime: true,
		TzOffset:    offset,
		EmployeeID:  employeeID,
		CreatedBy:   createdBy,
	}

	loc := _time.Location()
	if !utils.DateEqual(_time, time.Now().In(loc)) {
		timeRecord.Status = PENDING
		timeRecord.RegularTime = false
	}

	timeRecord.ID = uuid.NewV4().String()
	timeRecord.CreatedAt = time.Now()

	err := timeRecord.isValid()
	if err != nil {
		return nil, err
	}

	return &timeRecord, nil
}

func (t *TimeRecord) isValid() error {

	if t.Time.After(time.Now()) {
		return errors.New("the registration time must not be longer than the current time")
	}

	if t.EmployeeID == t.ApprovedBy {
		return errors.New("the employee who recorded the time cannot be the same person who approves")
	}

	if !t.RegularTime && t.Description == "" {
		return errors.New("the description must not be empty when the registration is done in an irregular period")
	}

	if t.EmployeeID == t.RefusedBy {
		return errors.New("the employee who recorded the time cannot be the same person who refuses")
	}

	_, err := govalidator.ValidateStruct(t)
	return err
}

func (t *TimeRecord) Approve(approvedBy string) error {

	if !govalidator.IsUUID(approvedBy) {
		return errors.New("the approved id must be a valid uuid")
	}

	if t.Status == APPROVED {
		return errors.New("the time record has already been approved")
	}

	if t.Status == REFUSED {
		return errors.New("the refused time record cannot be approved")
	}

	t.Status = APPROVED
	t.UpdatedAt = time.Now()
	t.ApprovedBy = approvedBy
	err := t.isValid()
	return err
}

func (t *TimeRecord) Refuse(refusedBy, refusedReason string) error {

	if !govalidator.IsUUID(refusedBy) {
		return errors.New("the refused id must be a valid uuid")
	}

	if t.Status == APPROVED {
		return errors.New("the approved time record cannot be refused")
	}

	if t.Status == REFUSED {
		return errors.New("the time record has already been refused")
	}

	if refusedReason == "" {
		return errors.New("the refused reason must not be empty")
	}

	t.Status = REFUSED
	t.UpdatedAt = time.Now()
	t.RefusedBy = refusedBy
	t.RefusedReason = refusedReason
	err := t.isValid()
	return err
}
