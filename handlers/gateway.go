package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ecoprohcm/DMS_BackendServer/models"
	"github.com/ecoprohcm/DMS_BackendServer/mqttSvc"
	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"github.com/gin-gonic/gin"
)

type GatewayHandler struct {
	deps *HandlerDependencies
}

func NewGatewayHandler(deps *HandlerDependencies) *GatewayHandler {
	return &GatewayHandler{
		deps,
	}
}

// Find all gateways and doorlocks info
// @Summary Find All Gateway
// @Schemes
// @Description find all gateways info
// @Produce json
// @Success 200 {array} []models.Gateway
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/gateways [get]
func (h *GatewayHandler) FindAllGateway(c *gin.Context) {
	gwList, err := h.deps.SvcOpts.GatewaySvc.FindAllGateway(c)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get all gateways failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, gwList)
}

// Find gateway and doorlock info by id
// @Summary Find Gateway By ID
// @Schemes
// @Description find gateway and doorlock info by gateway id
// @Produce json
// @Param        id	path	string	true	"Gateway ID"
// @Success 200 {object} models.Gateway
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/gateway/{id} [get]
func (h *GatewayHandler) FindGatewayByID(c *gin.Context) {
	id := c.Param("id")

	gw, err := h.deps.SvcOpts.GatewaySvc.FindGatewayByID(c, id)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get gateway failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, gw)
}

// Update gateway
// @Summary Update Gateway By Gateway ID
// @Schemes
// @Description Update gateway, must have "gateway_id" field. Send updated info to MQTT broker
// @Accept  json
// @Produce json
// @Param	data	body	models.SwagUpateGateway	true	"Fields need to update a gateway"
// @Success 200 {boolean} true
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/gateway [patch]
func (h *GatewayHandler) UpdateGateway(c *gin.Context) {
	gw := &models.Gateway{}
	err := c.ShouldBind(gw)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	isSuccess, err := h.deps.SvcOpts.GatewaySvc.UpdateGateway(c.Request.Context(), gw)
	if err != nil || !isSuccess {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update gateway failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	updated_gw, _ := h.deps.SvcOpts.GatewaySvc.FindGatewayByMacID(c.Request.Context(), gw.GatewayID)
	new_gw_log := &models.GatewayLog{}
	new_gw_log.StateType = "ConnectState"
	new_gw_log.GatewayID = gw.GatewayID
	new_gw_log.StateValue = updated_gw.ConnectState
	new_gw_log.LogTime = time.Now()
	h.deps.SvcOpts.LogSvc.CreateGatewayLog(c.Request.Context(), new_gw_log)

	t := h.deps.MqttClient.Publish(mqttSvc.TOPIC_SV_GATEWAY_U, 1, false, mqttSvc.ServerUpdateGatewayPayload(gw))
	if err := mqttSvc.HandleMqttErr(t); err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update gateway mqtt failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	utils.ResponseJson(c, http.StatusOK, isSuccess)
}

// Delete gateway
// @Summary Delete Gateway By Gateway ID
// @Schemes
// @Description Delete gateway using "" field. Send deleted info to MQTT broker
// @Accept  json
// @Produce json
// @Param	data	body	object{gateway_id=string}	true	"Gateway ID"
// @Success 200 {boolean} true
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/gateway [delete]
func (h *GatewayHandler) DeleteGateway(c *gin.Context) {
	dgw := &models.DeleteGateway{}
	err := c.ShouldBind(dgw)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	_, err1 := h.deps.SvcOpts.GatewaySvc.FindGatewayByMacID(c.Request.Context(), dgw.GatewayID)
	if err1 != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Find Gateway failed",
			ErrorMsg:   err1.Error(),
		})
		return
	}

	dls, err := h.deps.SvcOpts.UHFSvc.FindAllUHFByGatewayID(c.Request.Context(), dgw.GatewayID)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Find UHF failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	//delete gateway first
	isSuccess, err := h.deps.SvcOpts.GatewaySvc.DeleteGateway(c.Request.Context(), dgw.GatewayID)
	if err != nil || !isSuccess {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Delete gateway failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	t := h.deps.MqttClient.Publish(mqttSvc.TOPIC_SV_GATEWAY_D, 1, false, mqttSvc.ServerDeleteGatewayPayload(dgw.GatewayID))
	if err := mqttSvc.HandleMqttErr(t); err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Delete gateway mqtt failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	// delete UHFs belong to this gateway
	for i := 0; i < len(dls); i++ {
		isSuccess, err := h.deps.SvcOpts.UHFSvc.DeleteUHF(c.Request.Context(), strconv.FormatUint(uint64(dls[i].ID), 10))
		if err != nil || !isSuccess {
			utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Msg:        "Delete doorlock failed",
				ErrorMsg:   err.Error(),
			})
			return
		}
	}

	utils.ResponseJson(c, http.StatusOK, isSuccess)
}
