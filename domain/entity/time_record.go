package entity

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/c-4u/time-record-service/utils"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type TimeRecord struct {
	Base          `bson:",inline" valid:"-"`
	Time          time.Time        `json:"time,omitempty" gorm:"column:time;not null;unique_index:idx_employee_time" bson:"time" valid:"required"`
	Status        TimeRecordStatus `json:"status" gorm:"column:status;not null" bson:"status" valid:"timeRecordStatus"`
	Description   string           `json:"description,omitempty" gorm:"column:description;type:varchar(255)" bson:"description,omitempty" valid:"-"`
	RefusedReason string           `json:"refused_reason,omitempty" gorm:"column:refused_reason;type:varchar(255)" bson:"refused_reason,omitempty" valid:"-"`
	RegularTime   bool             `json:"regular_time" gorm:"column:regular_time" bson:"regular_time" valid:"-"`
	TzOffset      int              `json:"tz_offset" bson:"tz_offset" valid:"int,optional"`
	EmployeeID    *string          `json:"employee_id,omitempty" gorm:"column:employee_id;type:uuid;not null;unique_index:idx_employee_time" bson:"employee_id" valid:"uuid"`
	Employee      *Employee        `json:"-" valid:"-"`
	ApprovedBy    *string          `json:"approved_by,omitempty" gorm:"column:approved_by;type:uuid" bson:"approved_by,omitempty" valid:"-"`
	Approver      *Employee        `json:"-" valid:"-"`
	RefusedBy     *string          `json:"refused_by,omitempty" gorm:"column:refused_by;type:uuid" bson:"refused_by,omitempty" valid:"-"`
	Refuser       *Employee        `json:"-" valid:"-"`
	CreatedBy     *string          `json:"created_by,omitempty" gorm:"column:created_by;type:uuid;not null" bson:"created_by" valid:"uuid"`
	Creater       *Employee        `json:"-" valid:"-"`
	Token         *string          `json:"-" gorm:"column:token;not null" bson:"token" valid:"-"`
}

func NewTimeRecord(_time time.Time, description string, employee, creater *Employee) (*TimeRecord, error) {

	_, offset := _time.Zone()
	token := primitive.NewObjectID().Hex()
	timeRecord := TimeRecord{
		Time:        _time,
		Status:      APPROVED,
		Description: description,
		RegularTime: true,
		TzOffset:    offset,
		EmployeeID:  &employee.ID,
		Employee:    employee,
		CreatedBy:   &creater.ID,
		Creater:     creater,
		Token:       &token,
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

	// TODO: change 5 for company tolerance
	if t.Time.After(time.Now().Add(time.Minute * 5)) {
		return errors.New("the registration time must not be longer than the current time")
	}

	if !t.RegularTime && t.Description == "" {
		return errors.New("the description must not be empty when the registration is done in an irregular period")
	}

	_, err := govalidator.ValidateStruct(t)
	return err
}

func (t *TimeRecord) Approve(approver *Employee) error {

	if *t.EmployeeID == approver.ID {
		return errors.New("the employee who recorded the time cannot be the same person who approves")
	}

	if t.Status == APPROVED {
		return errors.New("the time record has already been approved")
	}

	if t.Status == REFUSED {
		return errors.New("the refused time record cannot be approved")
	}

	t.Status = APPROVED
	t.UpdatedAt = time.Now()
	t.ApprovedBy = &approver.ID
	t.Approver = approver
	err := t.isValid()
	return err
}

func (t *TimeRecord) Refuse(refuser *Employee, refusedReason string) error {

	if *t.EmployeeID == refuser.ID {
		return errors.New("the employee who recorded the time cannot be the same person who refuses")
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
	t.RefusedBy = &refuser.ID
	t.Refuser = refuser
	t.RefusedReason = refusedReason
	err := t.isValid()
	return err
}
