package entity

import (
	"time"

	"github.com/asaskevich/govalidator"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type Base struct {
	ID        string    `json:"id,omitempty" bson:"_id" valid:"uuid"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at" valid:"required"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty" valid:"-"`
}
