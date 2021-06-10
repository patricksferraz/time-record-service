package model

import (
	"github.com/asaskevich/govalidator"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type EmployeeClaims struct {
	ID    string   `valid:"uuid"`
	Roles []string `valid:"-"`
}

func (a *EmployeeClaims) isValid() error {
	_, err := govalidator.ValidateStruct(a)
	return err
}

func NewEmployeeClaims(employeeID string, roles []string) (*EmployeeClaims, error) {

	employeeClaims := EmployeeClaims{
		ID:    employeeID,
		Roles: roles,
	}

	err := employeeClaims.isValid()
	if err != nil {
		return nil, err
	}

	return &employeeClaims, nil
}
