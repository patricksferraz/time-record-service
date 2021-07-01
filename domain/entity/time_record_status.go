package entity

import (
	"github.com/asaskevich/govalidator"
)

func init() {
	govalidator.TagMap["timeRecordStatus"] = govalidator.Validator(func(str string) bool {
		res := str == PENDING.String()
		res = res || str == APPROVED.String()
		res = res || str == REFUSED.String()
		return res
	})
}

type TimeRecordStatus int

const (
	PENDING TimeRecordStatus = iota
	APPROVED
	REFUSED
)

func (t TimeRecordStatus) String() string {
	switch t {
	case PENDING:
		return "PENDING"
	case APPROVED:
		return "APPROVED"
	case REFUSED:
		return "REFUSED"
	}
	return ""
}
