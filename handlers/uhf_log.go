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

func (h *UHFStatusLogHandler) GetAllUHFStatusLogs(c *gin.Context) {
	dlslList, err := h.deps.SvcOpts.UHFStatusLogSvc.GetAllUHFStatusLogs(c)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get all UHF status logs failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, dlslList)
}

func (h *UHFStatusLogHandler) GetUHFStatusLogByUHFAddress(c *gin.Context) {
	doorId := c.Param("uhf_address")
	gl, err := h.deps.SvcOpts.UHFStatusLogSvc.GetUHFStatusLogByDoorID(c, doorId)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get UHF log failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, gl)
}

func (h *UHFStatusLogHandler) GetUHFStatusLogInTimeRange(c *gin.Context) {
	gateway_id := c.Param("gateway_id")
	uhf_address := c.Param("uhf_address")
	from := c.Param("from")
	to := c.Param("to")
	fromInt, _ := strconv.ParseInt(from, 10, 64)
	toInt, _ := strconv.ParseInt(to, 10, 64)
	fromFormatted := time.Unix(fromInt, 0).Format(models.DEFAULT_TIME_FORMAT)
	toFormatted := time.Unix(toInt, 0).Format(models.DEFAULT_TIME_FORMAT)
	dlslList, err := h.deps.SvcOpts.UHFStatusLogSvc.GetUHFStatusLogInTimeRange(fromFormatted, toFormatted, gateway_id, uhf_address)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Failed to get UHF status logs",
			ErrorMsg:   err.Error(),
		})
	}
	utils.ResponseJson(c, http.StatusOK, dlslList)
}

func (h *UHFStatusLogHandler) DeleteUHFStatusLogInTimeRange(c *gin.Context) {
	from := c.Param("fromTime")
	to := c.Param("toTime")
	fromInt, _ := strconv.ParseInt(from, 10, 64)
	toInt, _ := strconv.ParseInt(to, 10, 64)
	fromFormatted := time.Unix(fromInt, 0).Format(models.DEFAULT_TIME_FORMAT)
	toFormatted := time.Unix(toInt, 0).Format(models.DEFAULT_TIME_FORMAT)
	isSuccess, err := h.deps.SvcOpts.UHFStatusLogSvc.DeleteUHFStatusLogInTimeRange(fromFormatted, toFormatted)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Failed to delete UHF status logs",
			ErrorMsg:   err.Error(),
		})
	}
	utils.ResponseJson(c, http.StatusOK, isSuccess)
}

func (h *UHFStatusLogHandler) DeleteUHFStatusLogByDoorID(c *gin.Context) {
	doorId := c.Param("id")
	isSuccess, err := h.deps.SvcOpts.UHFStatusLogSvc.DeleteUHFStatusLogUHFID(doorId)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Failed to delete UHF status logs",
			ErrorMsg:   err.Error(),
		})
	}
	utils.ResponseJson(c, http.StatusOK, isSuccess)
}
