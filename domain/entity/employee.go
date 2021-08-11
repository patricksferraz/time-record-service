package entity

import (
	"time"

	"github.com/asaskevich/govalidator"
	pisvalidatior "github.com/patricksferraz/pisvalidator"
	uuid "github.com/satori/go.uuid"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)

	govalidator.TagMap["pis"] = govalidator.Validator(func(str string) bool {
		return pisvalidatior.ValidatePis(str)
	})
}

type Employee struct {
	Base        `json:",inline" valid:"required"`
	Pis         string        `json:"pis" gorm:"column:pis;type:varchar(25);unique" valid:"pis"`
	TimeRecords []*TimeRecord `json:"time_records,omitempty" gorm:"ForeignKey:EmployeeID" valid:"-"`
}

func NewEmployee(id, pis string, createdAt time.Time) (*Employee, error) {
	entity := &Employee{
		Pis: pis,
	}

	if id == "" {
		entity.ID = uuid.NewV4().String()
	} else {
		entity.ID = id
	}

	entity.CreatedAt = createdAt

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
