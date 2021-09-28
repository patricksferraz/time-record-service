package entity

import (
	"encoding/json"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type TimeRecordEvent struct {
	ID         string      `json:"id,omitempty" valid:"uuid"`
	TimeRecord *TimeRecord `json:"time_record,omitempty" valid:"-"`
}

func NewTimeRecordEvent(timeRecord *TimeRecord) (*TimeRecordEvent, error) {
	e := &TimeRecordEvent{
		ID:         uuid.NewV4().String(),
		TimeRecord: timeRecord,
	}

	if err := timeRecord.isValid(); err != nil {
		return nil, err
	}

	return e, nil
}

func (e *TimeRecordEvent) isValid() error {
	_, err := govalidator.ValidateStruct(e)
	return err
}

func (e *TimeRecordEvent) ToJson() ([]byte, error) {
	err := e.isValid()
	if err != nil {
		return nil, err
	}

	result, err := json.Marshal(e)
	if err != nil {
		return nil, nil
	}

	return result, nil
}
