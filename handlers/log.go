package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ecoprohcm/DMS_BackendServer/models"
	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"github.com/gin-gonic/gin"
)

type GatewayLogHandler struct {
	deps *HandlerDependencies
}

func NewGatewayLogHandler(deps *HandlerDependencies) *GatewayLogHandler {
	return &GatewayLogHandler{
		deps,
	}
}

// Find all gateway logs info
// @Summary Find All GatewayLog
// @Schemes
// @Description find all gateway logs info
// @Produce json
// @Success 200 {array} []models.GatewayLog
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/gateway_logs [get]
func (h *GatewayLogHandler) FindAllGatewayLog(c *gin.Context) {
	glList, err := h.deps.SvcOpts.LogSvc.FindAllGatewayLog(c)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get all gateway logs failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, glList)
}

// Find gateway log info by id
// @Summary Find GatewayLog By ID
// @Schemes
// @Description find gateway log info by id
// @Produce json
// @Param        id	path	string	true	"GatewayLog ID"
// @Success 200 {object} models.GatewayLog
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/gateway_logs/{id} [get]
func (h *GatewayLogHandler) FindGatewayLogByID(c *gin.Context) {
	id := c.Param("id")
	gl, err := h.deps.SvcOpts.LogSvc.FindGatewayLogByID(c, id)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get gateway log failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, gl)
}

// Find gateway log info by Gateway id
// @Summary Find GatewayLog By Gateway ID
// @Schemes
// @Description find gateway log info by Gateway id
// @Produce json
// @Param        id	path	string	true	"GatewayLog ID"
// @Success 200 {object} models.GatewayLog
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/gateway_logs/gateway_id/{id} [get]
func (h *GatewayLogHandler) FindGatewayByGatewayID(c *gin.Context) {
	id := c.Param("id")
	gl, err := h.deps.SvcOpts.LogSvc.FindGatewayByGatewayID(c, id)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get gateway log failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, gl)
}

// Find Gateway logs by period of time
// @Summary Find Gateway logs by period of time
// @Schemes
// @Description find Gateway logs by period of time
// @Produce json
// @Param        id	path	string	true	"GatewayLog ID"
// @Param 		 from path  string  true    "From Unix time"
// @Param 		 to path    string  true    "To Unix time"
// @Success 200 {object} []models.GatewayLog
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/gateway_logs/gateway_id/{gateway_id}/period/{from}/{to} [get]
func (h *GatewayLogHandler) FindGatewayLogsByGatewayIDAndTime(c *gin.Context) {
	gatewayId := c.Param("id")
	from := c.Param("from")
	to := c.Param("to")
	fromInt, _ := strconv.ParseInt(from, 10, 64)
	toInt, _ := strconv.ParseInt(to, 10, 64)
	fromFormatted := time.Unix(fromInt, 0).Format(models.DEFAULT_TIME_FORMAT)
	toFormatted := time.Unix(toInt, 0).Format(models.DEFAULT_TIME_FORMAT)
	glList, err := h.deps.SvcOpts.LogSvc.FindGatewayLogsByGatewayIDAndTime(gatewayId, fromFormatted, toFormatted)

	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Failed to get gateway logs",
			ErrorMsg:   err.Error(),
		})
	}
	utils.ResponseJson(c, http.StatusOK, glList)
}

// Find Gateway logs by period of time
// @Summary Find Gateway logs by period of time
// @Schemes
// @Description find Gateway logs by period of time
// @Produce json
// @Param        id	path	string	true	"GatewayLog ID"
// @Param 		 from path  string  true    "From Unix time"
// @Param 		 to path    string  true    "To Unix time"
// @Success 200 {object} []models.GatewayLog
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/gateway_logs/period/{from}/{to} [get]
func (h *GatewayLogHandler) FindGatewayLogByPeriod(c *gin.Context) {
	from := c.Param("from")
	to := c.Param("to")
	fromInt, _ := strconv.ParseInt(from, 10, 64)
	toInt, _ := strconv.ParseInt(to, 10, 64)
	fromFormatted := time.Unix(fromInt, 0).Format(models.DEFAULT_TIME_FORMAT)
	toFormatted := time.Unix(toInt, 0).Format(models.DEFAULT_TIME_FORMAT)
	glList, err := h.deps.SvcOpts.LogSvc.FindGatewayLogsByTime(fromFormatted, toFormatted)

	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Failed to get gateway logs",
			ErrorMsg:   err.Error(),
		})
	}
	utils.ResponseJson(c, http.StatusOK, glList)
}

// Delete Gateway logs in time range
// @Summary Delete GatewayLog In Time Range
// @Schemes
// @Description delete gateway logs in time range
// @Produce json
// @Param 		 from path  string  true    "From Unix time"
// @Param 		 to path    string  true    "To Unix time"
// @Success 200 {boolean} true
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/gateway_logs/period/{from}/{to} [delete]
func (h *GatewayLogHandler) DeleteGatewayLogInTimeRange(c *gin.Context) {
	from := c.Param("fromTime")
	to := c.Param("toTime")
	fromInt, _ := strconv.ParseInt(from, 10, 64)
	toInt, _ := strconv.ParseInt(to, 10, 64)
	fromFormatted := time.Unix(fromInt, 0).Format(models.DEFAULT_TIME_FORMAT)
	toFormatted := time.Unix(toInt, 0).Format(models.DEFAULT_TIME_FORMAT)
	isSuccess, err := h.deps.SvcOpts.LogSvc.DeleteGatewayLogInTimeRange(fromFormatted, toFormatted)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Failed to delete gateway logs",
			ErrorMsg:   "There is no record in this time range",
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, isSuccess)
}
