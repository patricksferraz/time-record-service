package schema

import (
	"encoding/json"

	"github.com/asaskevich/govalidator"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type EmployeeEvent struct {
	ID       string    `json:"id,omitempty" valid:"uuid"`
	Employee *Employee `json:"employee,omitempty" valid:"-"`
}

func NewEmployeeEvent(id, pis string) *EmployeeEvent {
	return &EmployeeEvent{}
}

func (e *EmployeeEvent) isValid() error {
	_, err := govalidator.ValidateStruct(e)
	return err
}

func (e *EmployeeEvent) ParseJson(data []byte) error {
	err := json.Unmarshal(data, e)
	if err != nil {
		return err
	}

	err = e.isValid()
	if err != nil {
		return err
	}

	return nil
}
