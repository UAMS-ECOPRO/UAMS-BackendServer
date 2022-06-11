package handlers

import (
	"context"
	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ActionHandler struct {
	deps *HandlerDependencies
}

func NewActionHandler(deps *HandlerDependencies) *ActionHandler {
	return &ActionHandler{
		deps,
	}
}

func (h *ActionHandler) FindAllAccesses(c *gin.Context) {
	gwList, err := h.deps.SvcOpts.ActionSvc.FindAllActions(c)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get all accesses failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, gwList)
}

func (h *ActionHandler) FindAccess(c *gin.Context) {
	//uhf_address := c.Param('uhf_address')
	//gateway_id := c.Param('gateway_id')
	epc := c.Param("epc")
	//time := c.Param('time')
	accesslist, err := h.deps.SvcOpts.ActionSvc.FindAllActionsByEPC(context.Background(), epc)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get all accesses failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, accesslist)
}
