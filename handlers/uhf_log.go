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

// Find all doorlock status logs info
// @Summary Find All DoorlockStatusLog
// @Schemes
// @Description find all doorlock status logs info
// @Produce json
// @Success 200 {array} []models.DoorlockStatusLog
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/doorlockStatusLogs [get]
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

// Find doorlock status logs by door id
// @Summary Find DoorlockStatusLog By DoorID
// @Schemes
// @Description find doorlock status logs by door id
// @Produce json
// @Param        id	path	string	true	"DoorlockStatusLog ID"
// @Success 200 {object} models.DoorlockStatusLog
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/doorlockStatusLog/{doorId} [get]
func (h *UHFStatusLogHandler) GetUHFStatusLogByUHFAddress(c *gin.Context) {
	doorId := c.Param("id")
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

// Find Doorlock status logs in time range
// @Summary Find DoorlockStatusLog In Time Range
// @Schemes
// @Description find doorlock status logs in time range
// @Produce json
// @Param 		 fromTime path  string  true    "From Unix time"
// @Param 		 toTime path    string  true    "To Unix time"
// @Success 200 {object} []models.DoorlockStatusLog
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/doorlockStatusLog/:fromTime/:toTime [get]
func (h *UHFStatusLogHandler) GetUHFStatusLogInTimeRange(c *gin.Context) {
	from := c.Param("fromTime")
	to := c.Param("toTime")
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

// Delete Doorlock status logs in time range
// @Summary Delete DoorlockStatusLog In Time Range
// @Schemes
// @Description delete doorlock status logs in time range
// @Produce json
// @Param 		 fromTime path  string  true    "From Unix time"
// @Param 		 toTime path    string  true    "To Unix time"
// @Success 200 {boolean} true
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/doorlockStatusLog/:fromTime/:toTime [delete]
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

// Delete Doorlock status logs by doorId
// @Summary Delete Doorlock status logs by doorId
// @Schemes
// @Description delete doorlock status logs by doorId
// @Produce json
// @Param        id	path	string	true	"DoorID"
// @Success 200 {boolean} true
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/doorlockStatusLog/door/:id [delete]
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
