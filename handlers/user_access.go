package handlers

import (
	"context"
	"github.com/ecoprohcm/DMS_BackendServer/models"
	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type UserAccessHandler struct {
	deps *HandlerDependencies
}

func NewUserAccessHandler(deps *HandlerDependencies) *UserAccessHandler {
	return &UserAccessHandler{
		deps,
	}
}

func (h *UserAccessHandler) FindAllAccesses(c *gin.Context) {
	gwList, err := h.deps.SvcOpts.UserAccessSvc.FindAllUserAccess(c)
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

func (h *UserAccessHandler) FindUserAccessByUserID(c *gin.Context) {
	id := c.Param("id")
	accesslist, err := h.deps.SvcOpts.UserAccessSvc.FindAllUserAccessByUserID(context.Background(), id)
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

func (h *UserAccessHandler) FindUserAccessByUserIDAndTimeRange(c *gin.Context) {
	user_id := c.Param("id")
	from := c.Param("from")
	to := c.Param("to")
	fromInt, _ := strconv.ParseInt(from, 10, 64)
	toInt, _ := strconv.ParseInt(to, 10, 64)
	fromFormatted := time.Unix(fromInt, 0).Format(models.DEFAULT_TIME_FORMAT)
	toFormatted := time.Unix(toInt, 0).Format(models.DEFAULT_TIME_FORMAT)
	glList, err := h.deps.SvcOpts.UserAccessSvc.FindUserAccessesByUserIDAndTimeRange(user_id, fromFormatted, toFormatted)

	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Failed to user access",
			ErrorMsg:   err.Error(),
		})
	}
	utils.ResponseJson(c, http.StatusOK, glList)
}
