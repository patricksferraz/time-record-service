package entity

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type Filter struct {
	FromDate   time.Time        `json:"from_date" valid:"optional"`
	ToDate     time.Time        `json:"to_date" valid:"optional"`
	Status     TimeRecordStatus `json:"status" valid:"timeRecordStatus,optional"`
	EmployeeID string           `json:"employee_id" valid:"uuid,optional"`
	ApprovedBy string           `json:"approved_by" valid:"uuid,optional"`
	RefusedBy  string           `json:"refused_by" valid:"uuid,optional"`
	CreatedBy  string           `json:"created_by" valid:"uuid,optional"`
	PageSize   int              `json:"page_size" valid:"optional"`
	PageToken  string           `json:"page_token" valid:"optional"`
}

func (e *Filter) isValid() error {
	if e.PageToken != "" {
		if !govalidator.IsMongoID(e.PageToken) {
			return errors.New("page token must be a valid token")
		}
	}
	_, err := govalidator.ValidateStruct(e)
	return err
}

func NewFilter(fromDate, toDate time.Time, status int, employeeID, approvedBy, refusedBy, createdBy string, pageSize int, pageToken string) (*Filter, error) {

	if pageSize == 0 {
		pageSize = 10
	}

	entity := &Filter{
		FromDate:   fromDate,
		ToDate:     toDate,
		Status:     TimeRecordStatus(status),
		EmployeeID: employeeID,
		ApprovedBy: approvedBy,
		RefusedBy:  refusedBy,
		CreatedBy:  createdBy,
		PageSize:   pageSize,
		PageToken:  pageToken,
	}

	err := entity.isValid()
	if err != nil {
		return nil, err
	}

	return entity, nil
}
