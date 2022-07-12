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

// Find all User Access log
// @Summary Find User Access log
// @Schemes
// @Description find all user access log
// @Produce json
// @Success 200 {array} []models.UserAccess
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/user_accesses [get]
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

// Find all User Access log by User ID
// @Summary Find User Access log by User ID
// @Schemes
// @Description find all user access log by user id
// @Produce json
// @Success 200 {array} []models.UserAccess
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/user_accesses/user_id/{id} [get]
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

// Find all User Access log by User ID and Time Range
// @Summary Find User Access log by User ID and Time Range
// @Schemes
// @Description find all user access log by user id and Time Range
// @Produce json
// @Success 200 {array} []models.UserAccess
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/user_accesses/user_id/{id}/period/{from}/{to} [get]
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

// Find all User Access log by User ID and Area Id
// @Summary Find User Access log by User ID and Area Id
// @Schemes
// @Description find all user access log by user id and Area Id
// @Produce json
// @Success 200 {array} []models.UserAccess
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/user_accesses/user_id/{id}/area_id/{area_id} [get]
func (h *UserAccessHandler) FindUserAccessByUserIDAndAreaID(c *gin.Context) {
	id := c.Param("id")
	area_id := c.Param("area_id")
	accesslist, err := h.deps.SvcOpts.UserAccessSvc.FindAllUserAccessByUserIDAndAreaID(context.Background(), id, area_id)
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

// Find all User Access log by User ID and Area Id and TimeRange
// @Summary Find User Access log by User ID and Area Id and TimeRange
// @Schemes
// @Description find all user access log by user id and Area Id and TimeRange
// @Produce json
// @Success 200 {array} []models.UserAccess
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/user_accesses/user_id/{id}/area_id/{area_id}/period/{from}/{to} [get]
func (h *UserAccessHandler) FindAllUserAccessByUserIDAndAreaIDinTimeRange(c *gin.Context) {
	user_id := c.Param("id")
	area_id := c.Param("area_id")
	from := c.Param("from")
	to := c.Param("to")
	fromInt, _ := strconv.ParseInt(from, 10, 64)
	toInt, _ := strconv.ParseInt(to, 10, 64)
	fromFormatted := time.Unix(fromInt, 0).Format(models.DEFAULT_TIME_FORMAT)
	toFormatted := time.Unix(toInt, 0).Format(models.DEFAULT_TIME_FORMAT)
	glList, err := h.deps.SvcOpts.UserAccessSvc.FindAllUserAccessByUserIDAndAreaIDinTimeRange(context.Background(), user_id, area_id, fromFormatted, toFormatted)

	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Failed to user access",
			ErrorMsg:   err.Error(),
		})
	}
	utils.ResponseJson(c, http.StatusOK, glList)
}

// Find all User Access log by Area Id
// @Summary Find User Access log by Area Id
// @Schemes
// @Description find all user access log by Area Id
// @Produce json
// @Success 200 {array} []models.UserAccess
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/user_accesses/area_id/{area_id} [get]
func (h *UserAccessHandler) FindAllUserAccessByAreaID(c *gin.Context) {
	area_id := c.Param("area_id")
	glList, err := h.deps.SvcOpts.UserAccessSvc.FindAllUserAccessByAreaID(context.Background(), area_id)

	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Failed to user access",
			ErrorMsg:   err.Error(),
		})
	}
	utils.ResponseJson(c, http.StatusOK, glList)
}

// Find all User Access log by Area Id and TimeRange
// @Summary Find User Access log by Area Id and TimeRange
// @Schemes
// @Description find all user access log by Area Id and TimeRange
// @Produce json
// @Success 200 {array} []models.UserAccess
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/user_accesses/area_id/{area_id}/period/{from}/{to} [get]
func (h *UserAccessHandler) FindAllUserAccessByAreaIDAndTimeRange(c *gin.Context) {
	area_id := c.Param("area_id")
	from := c.Param("from")
	to := c.Param("to")
	fromInt, _ := strconv.ParseInt(from, 10, 64)
	toInt, _ := strconv.ParseInt(to, 10, 64)
	fromFormatted := time.Unix(fromInt, 0).Format(models.DEFAULT_TIME_FORMAT)
	toFormatted := time.Unix(toInt, 0).Format(models.DEFAULT_TIME_FORMAT)
	glList, err := h.deps.SvcOpts.UserAccessSvc.FindAllUserAccessByAreaIDAndTimeRange(context.Background(), area_id, fromFormatted, toFormatted)

	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Failed to user access",
			ErrorMsg:   err.Error(),
		})
	}
	utils.ResponseJson(c, http.StatusOK, glList)
}

// Find all User Access log by TimeRange
// @Summary Find User Access log by TimeRange
// @Schemes
// @Description find all user access log by TimeRange
// @Produce json
// @Success 200 {array} []models.UserAccess
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/user_accesses/period/{from}/{to} [get]
func (h *UserAccessHandler) FindAllUserAccessTimeRange(c *gin.Context) {
	from := c.Param("from")
	to := c.Param("to")
	fromInt, _ := strconv.ParseInt(from, 10, 64)
	toInt, _ := strconv.ParseInt(to, 10, 64)
	fromFormatted := time.Unix(fromInt, 0).Format(models.DEFAULT_TIME_FORMAT)
	toFormatted := time.Unix(toInt, 0).Format(models.DEFAULT_TIME_FORMAT)
	glList, err := h.deps.SvcOpts.UserAccessSvc.FindAllUserAccessTimeRange(context.Background(), fromFormatted, toFormatted)

	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Failed to user access",
			ErrorMsg:   err.Error(),
		})
	}
	utils.ResponseJson(c, http.StatusOK, glList)
}
