package entity

import (
	"github.com/asaskevich/govalidator"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type Claims struct {
	EmployeeID string   `valid:"uuid"`
	Roles      []string `valid:"-"`
}

func (c *Claims) isValid() error {
	_, err := govalidator.ValidateStruct(c)
	return err
}

func NewClaims(employeeID string, roles []string) (*Claims, error) {

	claims := Claims{
		EmployeeID: employeeID,
		Roles:      roles,
	}

	err := claims.isValid()
	if err != nil {
		return nil, err
	}

	return &claims, nil
}
