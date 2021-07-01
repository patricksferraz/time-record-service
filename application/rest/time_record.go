package rest

import (
	"net/http"

	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/domain/entity"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/domain/service"
	"dev.azure.com/c4ut/TimeClock/_git/time-record-service/logger"
	"github.com/gin-gonic/gin"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmlogrus"
)

type TimeRecordRestService struct {
	TimeRecordService *service.TimeRecordService
	AuthMiddleware    *AuthMiddleware
}

// RegisterTimeRecord godoc
// @Security ApiKeyAuth
// @Summary register a new time record
// @ID registerTimeRecord
// @Tags Time Record
// @Description Router for registration a new time record
// @Accept json
// @Produce json
// @Param body body TimeRecordRequest true "JSON body for register a new time record"
// @Success 200 {object} ID
// @Failure 400 {object} HTTPError
// @Failure 403 {object} HTTPError
// @Router /time-records [post]
func (t *TimeRecordRestService) RegisterTimeRecord(ctx *gin.Context) {
	var req TimeRecordRequest

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

	timeRecordID, err := t.TimeRecordService.RegisterTimeRecord(ctx, req.Time, req.Description, req.EmployeeID, t.AuthMiddleware.Claims.EmployeeID)
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

	ctx.JSON(http.StatusOK, ID{ID: *timeRecordID})
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
func (t *TimeRecordRestService) ApproveTimeRecord(ctx *gin.Context) {
	var req ID

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

	err := t.TimeRecordService.ApproveTimeRecord(ctx, req.ID, t.AuthMiddleware.Claims.EmployeeID)
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
// @Param body body RefuseRequest true "JSON body for refuse a pending time record"
// @Success 200 {object} HTTPResponse
// @Failure 400 {object} HTTPError
// @Failure 403 {object} HTTPError
// @Router /time-records/{id}/refuse [post]
func (t *TimeRecordRestService) RefuseTimeRecord(ctx *gin.Context) {
	var uri ID
	var body RefuseRequest

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

	err := t.TimeRecordService.RefuseTimeRecord(ctx, uri.ID, body.RefusedReason, t.AuthMiddleware.Claims.EmployeeID)
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
func (t *TimeRecordRestService) FindTimeRecord(ctx *gin.Context) {
	var req ID

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

	timeRecord, err := t.TimeRecordService.FindTimeRecord(ctx, req.ID)
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
// @Summary search time records by employee id
// @ID searchTimeRecords
// @Tags Time Record
// @Description Search for employee time records by `id`
// @Accept json
// @Produce json
// @Param body query TimeRecordsRequest true "JSON body for search time records"
// @Success 200 {array} TimeRecord
// @Failure 400 {object} HTTPError
// @Failure 403 {object} HTTPError
// @Router /time-records [get]
func (t *TimeRecordRestService) SearchTimeRecords(ctx *gin.Context) {
	var body TimeRecordsRequest

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

	timeRecords, err := t.TimeRecordService.SearchTimeRecords(ctx, body.EmployeeID, body.FromDate, body.ToDate)
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

	ctx.JSON(http.StatusOK, timeRecords)
}

func NewTimeRecordRestService(service *service.TimeRecordService, authMiddleware *AuthMiddleware) *TimeRecordRestService {
	return &TimeRecordRestService{
		TimeRecordService: service,
		AuthMiddleware:    authMiddleware,
	}
}
