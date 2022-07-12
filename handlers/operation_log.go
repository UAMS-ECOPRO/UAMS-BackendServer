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

// Find all operation logs info
// @Summary Find All OperationLog
// @Schemes
// @Description find all operation logs info
// @Produce json
// @Success 200 {array} []models.OperationLog
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

// Find operation log info by gateway_id
// @Summary Find Operation By GatewayID
// @Schemes
// @Description find operation log info by gateway_id
// @Produce json
// @Param        id	path	string	true	"GatewayLog ID"
// @Success 200 {object} models.OperationLog
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/operation_logs/gateway_id/{gateway_id} [get]
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

// Find operation log info by id
// @Summary Find Operation By ID
// @Schemes
// @Description find operation log info by id
// @Produce json
// @Param        id	path	string	true	"ID"
// @Success 200 {object} models.OperationLog
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/operation_logs/{id} [get]
func (h *OperationLogHandler) FindAllOperationLogByID(c *gin.Context) {
	id := c.Param("id")
	ol, err := h.deps.SvcOpts.OperationLogSvc.GetOperationLogByID(c, id)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get gateway log failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, ol)
}

// Find operation logs by period of time
// @Summary Find operation logs by period of time
// @Schemes
// @Description find operation logs by period of time
// @Produce json
// @Param        id	path	string	true	"Gateway ID"
// @Param 		 from path  string  true    "From Unix time"
// @Param 		 to path    string  true    "To Unix time"
// @Success 200 {object} []models.OperationLog
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/operation_logs/gateway_id/{gateway_id}/period/{from}/{to} [get]
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

// Find Operation logs by period of time
// @Summary Find Operation logs by period of time
// @Schemes
// @Description find Operation logs by period of time
// @Produce json
// @Param        id	path	string	true	"Gateway ID"
// @Param 		 from path  string  true    "From Unix time"
// @Param 		 to path    string  true    "To Unix time"
// @Success 200 {object} []models.OperationLog
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/operation_logs/period/{from}/{to} [get]
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
