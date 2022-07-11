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

// Find all UHF info
// @Summary Find All UHF
// @Schemes
// @Description find all UHF info
// @Produce json
// @Success 200 {array} []models.UHF
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/uhfs [get]
func (h *UHFHandler) FindAllUHFs(c *gin.Context) {
	test := h.deps
	dlList, err := test.SvcOpts.UHFSvc.FindAllUHF(c)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get all UHfs failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, dlList)
}

// Find UHF info by id
// @Summary Find UHF By ID
// @Schemes
// @Description find UHF info by UHF id
// @Produce json
// @Param        id	path	string	true	"UHF ID"
// @Success 200 {object} models.UHF
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/uhf/{id} [get]
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

// Update UHF
// @Summary Update UHF By UHF Address and GatewayID
// @Schemes
// @Description Update UHF, must have "gatewayId" and "UHFAddress" field. Send updated info to MQTT broker
// @Accept  json
// @Produce json
// @Param	data	body	models.SwagUpdateUHF	true	"Fields need to update a UHF"
// @Success 200 {boolean} true
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/uhf [patch]
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

	existing_uhf, _ := h.deps.SvcOpts.UHFSvc.FindUHFByAddress(c.Request.Context(), dl.UHFAddress, dl.GatewayID)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "There is no UHF",
			ErrorMsg:   err.Error(),
		})
		return
	}
	if existing_uhf.AreaId == "" {
		if dl.AreaId == "" {
			utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Msg:        "Please input AreaID",
				ErrorMsg:   "Please input AreaID",
			})
			return
		}
		_, err = h.deps.SvcOpts.AreaSvc.FindAreaByID(c.Request.Context(), dl.AreaId)
		if err != nil {
			utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Msg:        "This area ID does not exist",
				ErrorMsg:   "This area ID does not exist",
			})
			return
		}
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

// Delete UHF
// @Summary Delete UHF By ID
// @Schemes
// @Description Delete UHF using "id" field. Send deleted info to MQTT broker
// @Accept  json
// @Produce json
// @Param	data	body	object{id=int}	true	"UHF Delete payload"
// @Success 200 {boolean} true
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/uhf [delete]
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
			Msg:        "Find UHF fail",
			ErrorMsg:   err.Error(),
		})
		return
	}

	t := h.deps.MqttClient.Publish(mqttSvc.TOPIC_SV_UHF_D, 1, false,
		mqttSvc.ServerDeleteUHFPayload(uhf))
	if err := mqttSvc.HandleMqttErr(t); err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Delete UHF mqtt failed",
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
