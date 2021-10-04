package entity

import (
	"fmt"
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
	Companies   []*Company    `json:"companies,omitempty" gorm:"many2many:companies_employees" valid:"-"`
}

func NewEmployee(id, pis string, company *Company) (*Employee, error) {
	entity := &Employee{Pis: pis}
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

func (e *Employee) AddCompany(company *Company) error {
	e.Companies = append(e.Companies, company)
	e.UpdatedAt = time.Now()
	err := e.isValid()
	return err
}

func (e *Employee) GetCompany(companyID string) (*Company, error) {
	for _, company := range e.Companies {
		if company.ID == companyID {
			return company, nil
		}
	}

	return nil, fmt.Errorf("employee does not belong to the company %s", companyID)
}
