package schema

import (
	"github.com/asaskevich/govalidator"
	pisvalidatior "github.com/patricksferraz/pisvalidator"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)

	govalidator.TagMap["pis"] = govalidator.Validator(func(str string) bool {
		return pisvalidatior.ValidatePis(str)
	})
}

type Employee struct {
	Base      `json:",inline" valid:"required"`
	Pis       string `json:"pis,omitempty" valid:"pis"`
	CompanyID string `json:"company_id" valid:"uuid"`
}

func NewEmployee(id, pis string) *Employee {
	return &Employee{}
}
