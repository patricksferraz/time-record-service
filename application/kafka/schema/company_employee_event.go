package schema

import (
	"encoding/json"

	"github.com/asaskevich/govalidator"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type CompanyEmployeeEvent struct {
	ID         string `json:"id,omitempty" valid:"uuid"`
	CompanyID  string `json:"company_id,omitempty" valid:"uuid"`
	EmployeeID string `json:"employee_id,omitempty" valid:"uuid"`
}

func NewCompanyEmployeeEvent() *CompanyEmployeeEvent {
	return &CompanyEmployeeEvent{}
}

func (e *CompanyEmployeeEvent) isValid() error {
	_, err := govalidator.ValidateStruct(e)
	return err
}

func (e *CompanyEmployeeEvent) ParseJson(data []byte) error {
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
