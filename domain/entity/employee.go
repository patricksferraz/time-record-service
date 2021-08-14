package entity

import (
	"time"

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
	Base        `json:",inline" valid:"required"`
	Pis         string        `json:"pis" gorm:"column:pis;type:varchar(25);not null;unique" valid:"pis"`
	TimeRecords []*TimeRecord `json:"time_records,omitempty" gorm:"ForeignKey:EmployeeID" valid:"-"`
	CompanyID   string        `json:"company_id" gorm:"column:company_id;type:uuid;not null" valid:"uuid"`
	Company     *Company      `json:"-" valid:"-"`
}

func NewEmployee(id, pis string, company *Company) (*Employee, error) {
	entity := &Employee{
		Pis:       pis,
		CompanyID: company.ID,
		Company:   company,
	}
	entity.ID = id
	entity.CreatedAt = time.Now()

	err := entity.isValid()
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (e *Employee) isValid() error {
	_, err := govalidator.ValidateStruct(e)
	return err
}
