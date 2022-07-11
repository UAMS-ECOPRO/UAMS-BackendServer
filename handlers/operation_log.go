package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ecoprohcm/DMS_BackendServer/models"
	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"github.com/gin-gonic/gin"
)

type OperationLogHandler struct {
	deps *HandlerDependencies
}

func NewOperationLogHandler(deps *HandlerDependencies) *OperationLogHandler {
	return &OperationLogHandler{
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
// @Router /v1/operation_logs [get]
func (h *OperationLogHandler) FindAllOperationLog(c *gin.Context) {
	glList, err := h.deps.SvcOpts.OperationLogSvc.GetAllOperationLogs(c)
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
// @Router /v1/operation_logs/{id} [get]
func (h *OperationLogHandler) FindOperationLogByGatewayID(c *gin.Context) {
	id := c.Param("gateway_id")
	gl, err := h.deps.SvcOpts.OperationLogSvc.GetOperationLogByGatewayID(c, id)
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
// @Router /v1/operation_logs/gateway_id/:gateway_id/period/:from/:to [get]
func (h *OperationLogHandler) FindOperationLogsByGatewayIDAndTime(c *gin.Context) {
	gatewayId := c.Param("gateway_id")
	from := c.Param("from")
	to := c.Param("to")
	fromInt, _ := strconv.ParseInt(from, 10, 64)
	toInt, _ := strconv.ParseInt(to, 10, 64)
	fromFormatted := time.Unix(fromInt, 0).Format(models.DEFAULT_TIME_FORMAT)
	toFormatted := time.Unix(toInt, 0).Format(models.DEFAULT_TIME_FORMAT)
	glList, err := h.deps.SvcOpts.OperationLogSvc.FindOperationLogsByGatewayIDAndTime(gatewayId, fromFormatted, toFormatted)

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
// @Router /v1/operation_logs/period/:from/:to [get]
func (h *OperationLogHandler) FindOperationLogsByTime(c *gin.Context) {
	from := c.Param("from")
	to := c.Param("to")
	fromInt, _ := strconv.ParseInt(from, 10, 64)
	toInt, _ := strconv.ParseInt(to, 10, 64)
	fromFormatted := time.Unix(fromInt, 0).Format(models.DEFAULT_TIME_FORMAT)
	toFormatted := time.Unix(toInt, 0).Format(models.DEFAULT_TIME_FORMAT)
	glList, err := h.deps.SvcOpts.OperationLogSvc.FindOperationLogsByTime(fromFormatted, toFormatted)

	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Failed to get gateway logs",
			ErrorMsg:   err.Error(),
		})
	}
	utils.ResponseJson(c, http.StatusOK, glList)
}
