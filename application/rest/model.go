package rest

import "time"

type Base struct {
	ID        string    `json:"id,omitempty" binding:"required,uuid"`
	CreatedAt time.Time `json:"created_at,omitempty" time_format:"RFC3339"`
	UpdatedAt time.Time `json:"updated_at,omitempty" time_format:"RFC3339"`
}

type TimeRecordRequest struct {
	EmployeeID  string    `json:"employee_id,omitempty" binding:"required,uuid"`
	Time        time.Time `json:"time,omitempty" time_format:"RFC3339" binding:"required"`
	Description string    `json:"description,omitempty"`
}

type TimeRecord struct {
	Base          `bson:",inline" valid:"-"`
	Time          time.Time `json:"time,omitempty" time_format:"RFC3339"`
	Status        int       `json:"status,omitempty"`
	Description   string    `json:"description,omitempty"`
	RefusedReason string    `json:"refused_reason,omitempty"`
	RegularTime   bool      `json:"regular_time,omitempty"`
	EmployeeID    string    `json:"employee_id,omitempty"`
	ApprovedBy    string    `json:"approved_by,omitempty"`
	RefusedBy     string    `json:"refused_by,omitempty"`
}

type HTTPResponse struct {
	Code    int    `json:"code,omitempty" example:"200"`
	Message string `json:"message,omitempty" example:"a message"`
}

type HTTPError struct {
	Code  int    `json:"code,omitempty" example:"400"`
	Error string `json:"error,omitempty" example:"status bad request"`
}

type ID struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type RefuseRequest struct {
	RefusedReason string `json:"refused_reason"`
}

type TimeRecordsRequest struct {
	EmployeeID string    `json:"employee_id,omitempty" binding:"required,uuid"`
	FromDate   time.Time `json:"from_date" form:"from_date" binding:"required"`
	ToDate     time.Time `json:"to_date" form:"to_date" binding:"required"`
}
