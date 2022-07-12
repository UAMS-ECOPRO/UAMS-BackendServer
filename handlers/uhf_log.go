package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ecoprohcm/DMS_BackendServer/models"
	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"github.com/gin-gonic/gin"
)

type UHFStatusLogHandler struct {
	deps *HandlerDependencies
}

func NewUHFStatusLogHandler(deps *HandlerDependencies) *UHFStatusLogHandler {
	return &UHFStatusLogHandler{
		deps,
	}
}

// Find all UHF Status logs info
// @Summary Find All UHF Status logs
// @Schemes
// @Description find all UHF Status logs
// @Produce json
// @Success 200 {array} []models.UHFStatusLog
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/uhf_logs [get]
func (h *UHFStatusLogHandler) GetAllUHFStatusLogs(c *gin.Context) {
	uhflList, err := h.deps.SvcOpts.UHFStatusLogSvc.GetAllUHFStatusLogs(c)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get all UHF logs failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, uhflList)
}

// Find all UHF Status logs info by Gateway ID and UHF Address
// @Summary Find All UHF Status logs by Gateway ID and UHF Address
// @Schemes
// @Description find all UHF Status logs by Gateway ID and UHF Address
// @Produce json
// @Success 200 {array} []models.UHFStatusLog
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/uhf_logs/gateway_id/{gateway_id}/uhf_address/{uhf_address} [get]
func (h *UHFStatusLogHandler) GetUHFStatusLogByUHFAddress(c *gin.Context) {
	uhf_address := c.Param("uhf_address")
	gateway_id := c.Param("gateway_id")
	uhfl, err := h.deps.SvcOpts.UHFStatusLogSvc.GetUHFStatusLogByUHFAddress(c, uhf_address, gateway_id)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get UHF log failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, uhfl)
}

// Find all UHF Status logs info by ID
// @Summary Find All UHF Status logs by ID
// @Schemes
// @Description find all UHF Status logs by ID
// @Produce json
// @Success 200 {array} []models.UHFStatusLog
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/uhf_logs/{id} [get]
func (h *UHFStatusLogHandler) GetUHFStatusLogsByID(c *gin.Context) {
	Id := c.Param("id")
	uhfl, err := h.deps.SvcOpts.UHFStatusLogSvc.GetUHFStatusLogByID(c, Id)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get UHF log failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, uhfl)
}

// Find all UHF Status logs info by Timerange
// @Summary Find All UHF Status logs by Timerange
// @Schemes
// @Description find all UHF Status logs by Timerange
// @Produce json
// @Success 200 {array} []models.UHFStatusLog
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/uhf_logs/period/{from}/{to} [get]
func (h *UHFStatusLogHandler) GetUHFStatusLogInTimeRange(c *gin.Context) {
	from := c.Param("from")
	to := c.Param("to")
	fromInt, _ := strconv.ParseInt(from, 10, 64)
	toInt, _ := strconv.ParseInt(to, 10, 64)
	fromFormatted := time.Unix(fromInt, 0).Format(models.DEFAULT_TIME_FORMAT)
	toFormatted := time.Unix(toInt, 0).Format(models.DEFAULT_TIME_FORMAT)
	dlslList, err := h.deps.SvcOpts.UHFStatusLogSvc.GetUHFStatusLogInTimeRange(fromFormatted, toFormatted)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Failed to get UHF status logs",
			ErrorMsg:   err.Error(),
		})
	}
	utils.ResponseJson(c, http.StatusOK, dlslList)
}

// Find all UHF Status logs info by GatewayID, AddressID and TimeRange
// @Summary Find All UHF Status logs by GatewayID, AddressID and TimeRange
// @Schemes
// @Description find all UHF Status logs by GatewayID, AddressID and TimeRange
// @Produce json
// @Success 200 {array} []models.UHFStatusLog
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/uhf_logs/gateway_id/{gateway_id}/uhf_address/{uhf_address}/period/{from}/{to} [get]
func (h *UHFStatusLogHandler) GetUHFStatusLogBYGatewayIDAndUHFAddressInTimeRange(c *gin.Context) {
	gateway_id := c.Param("gateway_id")
	uhf_address := c.Param("uhf_address")
	from := c.Param("from")
	to := c.Param("to")
	fromInt, _ := strconv.ParseInt(from, 10, 64)
	toInt, _ := strconv.ParseInt(to, 10, 64)
	fromFormatted := time.Unix(fromInt, 0).Format(models.DEFAULT_TIME_FORMAT)
	toFormatted := time.Unix(toInt, 0).Format(models.DEFAULT_TIME_FORMAT)
	dlslList, err := h.deps.SvcOpts.UHFStatusLogSvc.GetUHFStatusLogBYGatewayIDAndUHFAddressInTimeRange(fromFormatted, toFormatted, gateway_id, uhf_address)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Failed to get UHF status logs",
			ErrorMsg:   err.Error(),
		})
	}
	utils.ResponseJson(c, http.StatusOK, dlslList)
}
