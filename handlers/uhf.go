package handlers

import (
	"github.com/ecoprohcm/DMS_BackendServer/models"
	"github.com/ecoprohcm/DMS_BackendServer/mqttSvc"
	"net/http"

	//"github.com/ecoprohcm/DMS_BackendServer/models"
	//"github.com/ecoprohcm/DMS_BackendServer/mqttSvc"
	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"github.com/gin-gonic/gin"
)

type UHFHandler struct {
	deps *HandlerDependencies
}

func NewUHFHandler(deps *HandlerDependencies) *UHFHandler {
	return &UHFHandler{
		deps,
	}
}

// Find all doorlocks info
// @Summary Find All Doorlock
// @Schemes
// @Description find all doorlocks info
// @Produce json
// @Success 200 {array} []models.UHF
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/doorlocks [get]
func (h *UHFHandler) FindAllUHFs(c *gin.Context) {
	test := h.deps
	dlList, err := test.SvcOpts.UHFSvc.FindAllUHF(c)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get all doorlocks failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, dlList)
}

// Find doorlock info by id
// @Summary Find Doorlock By ID
// @Schemes
// @Description find doorlock info by doorlock id
// @Produce json
// @Param        id	path	string	true	"Doorlock ID"
// @Success 200 {object} models.Doorlock
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/doorlock/{id} [get]
func (h *UHFHandler) FindUHFByID(c *gin.Context) {
	id := c.Param("id")

	dl, err := h.deps.SvcOpts.UHFSvc.FindUHFByID(c, id)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get UHF failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, dl)
}

// Update doorlock
// @Summary Update Doorlock By Doorlock Address and GatewayID
// @Schemes
// @Description Update doorlock, must have "gatewayId" and "doorlockAddress" field. Send updated info to MQTT broker
// @Accept  json
// @Produce json
// @Param	data	body	models.SwagUpdateDoorlock	true	"Fields need to update a doorlock"
// @Success 200 {boolean} true
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/doorlock [patch]
func (h *UHFHandler) UpdateUHF(c *gin.Context) {
	dl := &models.UHF{}
	err := c.ShouldBind(dl)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	isSuccess, err := h.deps.SvcOpts.UHFSvc.UpdateUHF(c.Request.Context(), dl)
	if err != nil || !isSuccess {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update uhf failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	updated_UHF, err := h.deps.SvcOpts.UHFSvc.FindUHFByAddress(c.Request.Context(), dl.UHFAddress, dl.GatewayID)

	var new_UHF_status_log_state = &models.UHFStatusLog{}
	new_UHF_status_log_state.GatewayID = updated_UHF.GatewayID
	new_UHF_status_log_state.UHFAddress = updated_UHF.UHFAddress
	new_UHF_status_log_state.StateType = "Active State"
	new_UHF_status_log_state.StateValue = dl.ActiveState
	h.deps.SvcOpts.UHFStatusLogSvc.CreateUHFStatusLog(c.Request.Context(), new_UHF_status_log_state)

	t := h.deps.MqttClient.Publish(mqttSvc.TOPIC_SV_UHF_U, 1, false,
		mqttSvc.ServerUpdateUHFPayload(dl))
	if err := mqttSvc.HandleMqttErr(t); err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update uhf mqtt failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	utils.ResponseJson(c, http.StatusOK, isSuccess)
}

// Delete doorlock
// @Summary Delete Doorlock By ID
// @Schemes
// @Description Delete doorlock using "id" field. Send deleted info to MQTT broker
// @Accept  json
// @Produce json
// @Param	data	body	object{id=int}	true	"Doorlock Delete payload"
// @Success 200 {boolean} true
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/doorlock [delete]
func (h *UHFHandler) DeleteUHF(c *gin.Context) {
	dl := &models.UHFDelete{}
	err := c.ShouldBind(dl)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	uhf, err := h.deps.SvcOpts.UHFSvc.FindUHFByID(c.Request.Context(), dl.ID)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Find doorlock fail",
			ErrorMsg:   err.Error(),
		})
		return
	}

	t := h.deps.MqttClient.Publish(mqttSvc.TOPIC_SV_UHF_D, 1, false,
		mqttSvc.ServerDeleteUHFPayload(uhf))
	if err := mqttSvc.HandleMqttErr(t); err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Delete doorlock mqtt failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	isSuccess, err := h.deps.SvcOpts.UHFSvc.DeleteUHF(c.Request.Context(), dl.ID)
	if err != nil || !isSuccess {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Delete UHF failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, isSuccess)

}
