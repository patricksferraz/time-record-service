package rest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patricksferraz/time-record-service/domain/entity"
	"github.com/patricksferraz/time-record-service/domain/service"
	"github.com/patricksferraz/time-record-service/logger"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmlogrus"
)

type RestService struct {
	Service        *service.Service
	AuthMiddleware *AuthMiddleware
}

func NewRestService(service *service.Service, authMiddleware *AuthMiddleware) *RestService {
	return &RestService{
		Service:        service,
		AuthMiddleware: authMiddleware,
	}
}

// RegisterTimeRecord godoc
// @Security ApiKeyAuth
// @Summary register a new time record
// @ID registerTimeRecord
// @Tags Time Record
// @Description Router for registration a new time record
// @Accept json
// @Produce json
// @Param body body RegisterTimeRecordRequest true "JSON body for register a new time record"
// @Success 200 {object} RegisterTimeRecordResponse
// @Failure 400 {object} HTTPError
// @Failure 403 {object} HTTPError
// @Router /time-records [post]
func (t *RestService) RegisterTimeRecord(ctx *gin.Context) {
	var req RegisterTimeRecordRequest

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		ctx.JSON(
			http.StatusBadRequest,
			HTTPError{
				Code:  http.StatusBadRequest,
				Error: err.Error(),
			},
		)
		return
	}
	log.WithField("json", req).Info("handling TimeRecord request")

	timeRecordID, err := t.Service.RegisterTimeRecord(ctx, req.Time, req.Description, req.EmployeeID, req.CompanyID, t.AuthMiddleware.Claims.EmployeeID)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		ctx.JSON(
			http.StatusForbidden,
			HTTPError{
				Code:  http.StatusForbidden,
				Error: err.Error(),
			},
		)
		return
	}
	log.WithField("timeRecordID", timeRecordID).Info("timeRecordID registered")

	ctx.JSON(http.StatusOK, RegisterTimeRecordResponse{ID: *timeRecordID})
}

// ApproveTimeRecord godoc
// @Security ApiKeyAuth
// @Summary approve a pending time record
// @ID approveTimeRecord
// @Tags Time Record
// @Description Router for approve a pending time record
// @Accept json
// @Produce json
// @Param id path string true "Time Record ID"
// @Success 200 {object} HTTPResponse
// @Failure 400 {object} HTTPError
// @Failure 403 {object} HTTPError
// @Router /time-records/{id}/approve [post]
func (t *RestService) ApproveTimeRecord(ctx *gin.Context) {
	var req ApproveTimeRecordRequest

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	if err := ctx.ShouldBindUri(&req); err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		ctx.JSON(
			http.StatusBadRequest,
			HTTPError{
				Code:  http.StatusBadRequest,
				Error: err.Error(),
			},
		)
		return
	}
	log.WithField("uri", req).Info("uri id request")

	err := t.Service.ApproveTimeRecord(ctx, req.ID, t.AuthMiddleware.Claims.EmployeeID)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		ctx.JSON(
			http.StatusForbidden,
			HTTPError{
				Code:  http.StatusForbidden,
				Error: err.Error(),
			},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		HTTPResponse{
			Code:    http.StatusOK,
			Message: "successfully " + entity.APPROVED.String(),
		},
	)
}

// RefuseTimeRecord godoc
// @Security ApiKeyAuth
// @Summary refuse a pending time record
// @ID refuseTimeRecord
// @Tags Time Record
// @Description Router for refuse a pending time record
// @Accept json
// @Produce json
// @Param id path string true "Time Record ID"
// @Param body body RefuseTimeRecordRequest true "JSON body for refuse a pending time record"
// @Success 200 {object} HTTPResponse
// @Failure 400 {object} HTTPError
// @Failure 403 {object} HTTPError
// @Router /time-records/{id}/refuse [post]
func (t *RestService) RefuseTimeRecord(ctx *gin.Context) {
	var uri IDRequest
	var body RefuseTimeRecordRequest

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	if err := ctx.ShouldBindUri(&uri); err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		ctx.JSON(
			http.StatusBadRequest,
			HTTPError{
				Code:  http.StatusBadRequest,
				Error: err.Error(),
			},
		)
		return
	}
	log.WithField("uri", uri).Info("uri id request")

	if err := ctx.ShouldBindJSON(&body); err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		ctx.JSON(
			http.StatusBadRequest,
			HTTPError{
				Code:  http.StatusBadRequest,
				Error: err.Error(),
			},
		)
		return
	}
	log.WithField("body", body).Info("handling Refuse request")

	err := t.Service.RefuseTimeRecord(ctx, uri.ID, body.RefusedReason, t.AuthMiddleware.Claims.EmployeeID)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		ctx.JSON(
			http.StatusForbidden,
			HTTPError{
				Code:  http.StatusForbidden,
				Error: err.Error(),
			},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		HTTPResponse{
			Code:    http.StatusOK,
			Message: "successfully " + entity.REFUSED.String(),
		},
	)
}

// FindTimeRecord godoc
// @Security ApiKeyAuth
// @Summary find a time record
// @ID findTimeRecord
// @Tags Time Record
// @Description Router for find a time record
// @Accept json
// @Produce json
// @Param id path string true "Time Record ID"
// @Success 200 {object} TimeRecord
// @Failure 400 {object} HTTPError
// @Failure 403 {object} HTTPError
// @Router /time-records/{id} [get]
func (t *RestService) FindTimeRecord(ctx *gin.Context) {
	var req FindTimeRecordRequest

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	if err := ctx.ShouldBindUri(&req); err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		ctx.JSON(
			http.StatusBadRequest,
			HTTPError{
				Code:  http.StatusBadRequest,
				Error: err.Error(),
			},
		)
		return
	}
	log.WithField("uri", req).Info("uri id request")

	timeRecord, err := t.Service.FindTimeRecord(ctx, req.ID)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		ctx.JSON(
			http.StatusForbidden,
			HTTPError{
				Code:  http.StatusForbidden,
				Error: err.Error(),
			},
		)
		return
	}
	log.WithField("timeRecord", timeRecord).Info("timeRecord finded")

	ctx.JSON(http.StatusOK, timeRecord)
}

// SearchTimeRecords godoc
// @Security ApiKeyAuth
// @Summary search time records by filter
// @ID searchTimeRecords
// @Tags Time Record
// @Description Search for employee time records by `filter`
// @Accept json
// @Produce json
// @Param body query SearchTimeRecordsRequest true "JSON body for search time records"
// @Success 200 {array} SearchTimeRecordsResponse
// @Failure 400 {object} HTTPError
// @Failure 403 {object} HTTPError
// @Router /time-records [get]
func (t *RestService) SearchTimeRecords(ctx *gin.Context) {
	var body SearchTimeRecordsRequest

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	if err := ctx.ShouldBindQuery(&body); err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		ctx.JSON(
			http.StatusBadRequest,
			HTTPError{
				Code:  http.StatusBadRequest,
				Error: err.Error(),
			},
		)
		return
	}
	log.WithField("query", body).Info("query TimeRecords request")

	nextPageToken, timeRecords, err := t.Service.SearchTimeRecords(ctx, body.FromDate, body.ToDate, body.Status, body.EmployeeID, body.ApprovedBy, body.RefusedBy, body.CreatedBy, body.CompanyID, body.PageSize, body.PageToken)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		ctx.JSON(
			http.StatusForbidden,
			HTTPError{
				Code:  http.StatusForbidden,
				Error: err.Error(),
			},
		)
		return
	}
	log.WithField("timeRecords", timeRecords).Info("timeRecords searched")

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"next_page_token": *nextPageToken,
			"time_records":    timeRecords,
		},
	)
}

// ExportTimeRecords godoc
// @Security ApiKeyAuth
// @Summary export time records by filter
// @ID exportTimeRecords
// @Tags Time Record
// @Description Export for employee time records by `filter`
// @Accept json
// @Produce json
// @Param body query ExportTimeRecordsRequest true "JSON body for search time records"
// @Success 200 {array} ExportTimeRecordsResponse
// @Failure 400 {object} HTTPError
// @Failure 403 {object} HTTPError
// @Router /time-records/export [get]
func (t *RestService) ExportTimeRecords(ctx *gin.Context) {
	var body ExportTimeRecordsRequest

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	if err := ctx.ShouldBindQuery(&body); err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		ctx.JSON(
			http.StatusBadRequest,
			HTTPError{
				Code:  http.StatusBadRequest,
				Error: err.Error(),
			},
		)
		return
	}
	log.WithField("query", body).Info("query TimeRecords request")

	nextPageToken, registers, err := t.Service.ExportTimeRecords(ctx, body.FromDate, body.ToDate, body.Status, body.EmployeeID, body.ApprovedBy, body.RefusedBy, body.CreatedBy, body.CompanyID, body.PageSize, body.PageToken, *t.AuthMiddleware.AccessToken)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		ctx.JSON(
			http.StatusForbidden,
			HTTPError{
				Code:  http.StatusForbidden,
				Error: err.Error(),
			},
		)
		return
	}
	log.WithField("registers", registers).Info("registers exported")

	// TODO: adds async export
	if body.AsFile {
		file, err := ioutil.TempFile("/tmp", *nextPageToken)
		if err != nil {
			log.WithError(err)
			apm.CaptureError(ctx, err).Send()
			ctx.JSON(
				http.StatusForbidden,
				HTTPError{
					Code:  http.StatusForbidden,
					Error: err.Error(),
				},
			)
			return
		}
		defer os.Remove(file.Name())

		for _, data := range registers {
			_, err := file.WriteString(string(*data) + "\n")
			if err != nil {
				log.WithError(err)
				apm.CaptureError(ctx, err).Send()
				ctx.JSON(
					http.StatusForbidden,
					HTTPError{
						Code:  http.StatusForbidden,
						Error: err.Error(),
					},
				)
				return
			}
		}

		b, err := ioutil.ReadFile(file.Name())
		if err != nil {
			log.WithError(err)
			apm.CaptureError(ctx, err).Send()
			ctx.JSON(
				http.StatusForbidden,
				HTTPError{
					Code:  http.StatusForbidden,
					Error: err.Error(),
				},
			)
			return
		}

		_t := time.Now()
		m := http.DetectContentType(b[:512])
		filename := fmt.Sprintf("%d%02d%02d%02d%02d%02d%s.txt", _t.Year(), _t.Month(), _t.Day(), _t.Hour(), _t.Minute(), _t.Second(), *nextPageToken)
		ctx.Header("Content-Disposition", "attachment; filename="+filename)
		ctx.Data(http.StatusOK, m, b)
	} else {
		ctx.JSON(
			http.StatusOK,
			gin.H{
				"next_page_token": *nextPageToken,
				"registers":       registers,
			},
		)
	}

}
