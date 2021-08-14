package schema

import (
	"github.com/asaskevich/govalidator"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type Company struct {
	Base `json:",inline" valid:"required"`
}

func NewCompany() *Company {
	return &Company{}
}
