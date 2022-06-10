package handlers

import (
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
