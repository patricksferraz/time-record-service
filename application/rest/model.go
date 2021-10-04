package rest

import (
	"time"
)

type Base struct {
	ID        string    `json:"id,omitempty" binding:"uuid"`
	CreatedAt time.Time `json:"created_at,omitempty" time_format:"RFC3339"`
	UpdatedAt time.Time `json:"updated_at,omitempty" time_format:"RFC3339"`
}

type RegisterTimeRecordRequest struct {
	EmployeeID  string    `json:"employee_id,omitempty" binding:"required,uuid"`
	CompanyID   string    `json:"company_id,omitempty" binding:"required,uuid"`
	Time        time.Time `json:"time,omitempty" time_format:"RFC3339" binding:"required"`
	Description string    `json:"description,omitempty"`
}

type RegisterTimeRecordResponse struct {
	ID string `json:"id" binding:"uuid"`
}

type TimeRecord struct {
	Base          `json:",inline"`
	Time          time.Time `json:"time,omitempty" time_format:"RFC3339"`
	Status        int       `json:"status,omitempty"`
	Description   string    `json:"description,omitempty"`
	RefusedReason string    `json:"refused_reason,omitempty"`
	RegularTime   bool      `json:"regular_time,omitempty"`
	TzOffset      int       `json:"tz_offset,omitempty"`
	EmployeeID    string    `json:"employee_id,omitempty"`
	ApprovedBy    string    `json:"approved_by,omitempty"`
	RefusedBy     string    `json:"refused_by,omitempty"`
	CreatedBy     string    `json:"created_by,omitempty"`
}

type HTTPResponse struct {
	Code    int    `json:"code,omitempty" example:"200"`
	Message string `json:"message,omitempty" example:"a message"`
}

type HTTPError struct {
	Code  int    `json:"code,omitempty" example:"400"`
	Error string `json:"error,omitempty" example:"status bad request"`
}

type ApproveTimeRecordRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type FindTimeRecordRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type IDRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type RefuseTimeRecordRequest struct {
	RefusedReason string `json:"refused_reason" binding:"required"`
}

type SearchTimeRecordsRequest struct {
	Filter `json:",inline"`
}

type ExportTimeRecordsRequest struct {
	Filter `json:",inline"`
	AsFile bool `json:"as_file" form:"as_file" default:"false"`
}

type Filter struct {
	FromDate   time.Time `json:"from_date" form:"from_date"`
	ToDate     time.Time `json:"to_date" form:"to_date"`
	Status     int       `json:"status" form:"status"`
	EmployeeID string    `json:"employee_id" form:"employee_id"`
	ApprovedBy string    `json:"approved_by" form:"approved_by"`
	RefusedBy  string    `json:"refused_by" form:"refused_by"`
	CreatedBy  string    `json:"created_by" form:"created_by"`
	PageSize   int       `json:"page_size" form:"page_size" default:"10"`
	PageToken  string    `json:"page_token" form:"page_token"`
}

type SearchTimeRecordsResponse struct {
	NextPageToken string       `json:"next_page_token"`
	TimeRecords   []TimeRecord `json:"time_records"`
}

type ExportTimeRecordsResponse struct {
	NextPageToken string   `json:"next_page_token"`
	Registers     []string `json:"registers"`
}
