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

type PackageAccessHandler struct {
	deps *HandlerDependencies
}

func NewPackageAccessHandler(deps *HandlerDependencies) *PackageAccessHandler {
	return &PackageAccessHandler{
		deps,
	}
}

// Find all Package Access log
// @Summary Find Package Access log
// @Schemes
// @Description find all package access log
// @Produce json
// @Success 200 {array} []models.PackageAccess
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/package_accesses [get]
func (h *PackageAccessHandler) FindAllAccesses(c *gin.Context) {
	gwList, err := h.deps.SvcOpts.PackageAccessSvc.FindAllPackageAccess(c)
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

// Find all Package Access log by Package ID
// @Summary Find Package Access log by Package ID
// @Schemes
// @Description find all package access log by Package ID
// @Produce json
// @Success 200 {array} []models.PackageAccess
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/package_accesses/package_id/{id} [get]
func (h *PackageAccessHandler) FindPackageAccessByPackageID(c *gin.Context) {
	id := c.Param("id")
	accesslist, err := h.deps.SvcOpts.PackageAccessSvc.FindAllPackageAccessByPackageID(context.Background(), id)
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

// Find all Package Access log by Package ID and TimeRange
// @Summary Find Package Access log by Package ID and TimeRange
// @Schemes
// @Description find all package access log by Package ID and TimeRange
// @Produce json
// @Success 200 {array} []models.PackageAccess
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/package_accesses/package_id/{id}/period/{from}/{to} [get]
func (h *PackageAccessHandler) FindPackageAccessByPackageIDAndTimeRange(c *gin.Context) {
	package_id := c.Param("id")
	from := c.Param("from")
	to := c.Param("to")
	fromInt, _ := strconv.ParseInt(from, 10, 64)
	toInt, _ := strconv.ParseInt(to, 10, 64)
	fromFormatted := time.Unix(fromInt, 0).Format(models.DEFAULT_TIME_FORMAT)
	toFormatted := time.Unix(toInt, 0).Format(models.DEFAULT_TIME_FORMAT)
	glList, err := h.deps.SvcOpts.PackageAccessSvc.FindPackageAccessByPackageIDAndTimeRange(context.Background(), package_id, fromFormatted, toFormatted)

	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Failed to user access",
			ErrorMsg:   err.Error(),
		})
	}
	utils.ResponseJson(c, http.StatusOK, glList)
}

// Find all Package Access log by Package ID and AreaID and TimeRange
// @Summary Find Package Access log by Package ID and AreaID and TimeRange
// @Schemes
// @Description find all package access log by Package ID and AreaID and TimeRange
// @Produce json
// @Success 200 {array} []models.PackageAccess
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/package_accesses/package_id/{id}/area_id/{area_id}/period/{from}/{to} [get]
func (h *PackageAccessHandler) FindAllPackageAccessByPackageIDAndAreaIDinTimeRange(c *gin.Context) {
	package_id := c.Param("id")
	area_id := c.Param("area_id")
	from := c.Param("from")
	to := c.Param("to")
	fromInt, _ := strconv.ParseInt(from, 10, 64)
	toInt, _ := strconv.ParseInt(to, 10, 64)
	fromFormatted := time.Unix(fromInt, 0).Format(models.DEFAULT_TIME_FORMAT)
	toFormatted := time.Unix(toInt, 0).Format(models.DEFAULT_TIME_FORMAT)
	glList, err := h.deps.SvcOpts.PackageAccessSvc.FindAllPackageAccessByPackageIDAndAreaIDinTimeRange(context.Background(), package_id, area_id, fromFormatted, toFormatted)

	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Failed to user access",
			ErrorMsg:   err.Error(),
		})
	}
	utils.ResponseJson(c, http.StatusOK, glList)
}

// Find all Package Access log by Package ID and AreaID
// @Summary Find Package Access log by Package ID and AreaID
// @Schemes
// @Description find all package access log by Package ID and AreaID
// @Produce json
// @Success 200 {array} []models.PackageAccess
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/package_accesses/package_id/{id}/area_id/{area_id} [get]
func (h *PackageAccessHandler) FindPackageAccessByPackageIDAndAreaID(c *gin.Context) {
	package_id := c.Param("id")
	area_id := c.Param("area_id")
	glList, err := h.deps.SvcOpts.PackageAccessSvc.FindPackageAccessByPackageIDAndAreaID(context.Background(), package_id, area_id)

	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Failed to user access",
			ErrorMsg:   err.Error(),
		})
	}
	utils.ResponseJson(c, http.StatusOK, glList)
}

// Find all Package Access log by AreaID
// @Summary Find Package Access log by AreaID
// @Schemes
// @Description find all package access log by  and AreaID
// @Produce json
// @Success 200 {array} []models.PackageAccess
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/package_accesses/area_id/{area_id} [get]
func (h *PackageAccessHandler) FindAllPackageAccessByAreaID(c *gin.Context) {
	area_id := c.Param("area_id")
	glList, err := h.deps.SvcOpts.PackageAccessSvc.FindAllPackageAccessByAreaID(context.Background(), area_id)

	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Failed to user access",
			ErrorMsg:   err.Error(),
		})
	}
	utils.ResponseJson(c, http.StatusOK, glList)
}

// Find all Package Access log by AreaID and TimeRange
// @Summary Find Package Access log by AreaID and TimeRange
// @Schemes
// @Description find all package access log by AreaID and TimeRange
// @Produce json
// @Success 200 {array} []models.PackageAccess
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/package_accesses/area_id/{area_id}/period/{from}/{to} [get]
func (h *PackageAccessHandler) FindAllPackageAccessByAreaIDAndTimeRange(c *gin.Context) {
	area_id := c.Param("area_id")

	from := c.Param("from")
	to := c.Param("to")
	fromInt, _ := strconv.ParseInt(from, 10, 64)
	toInt, _ := strconv.ParseInt(to, 10, 64)
	fromFormatted := time.Unix(fromInt, 0).Format(models.DEFAULT_TIME_FORMAT)
	toFormatted := time.Unix(toInt, 0).Format(models.DEFAULT_TIME_FORMAT)
	glList, err := h.deps.SvcOpts.PackageAccessSvc.FindAllPackageAccessByAreaIDAndTimeRange(context.Background(), area_id, fromFormatted, toFormatted)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Failed to user access",
			ErrorMsg:   err.Error(),
		})
	}
	utils.ResponseJson(c, http.StatusOK, glList)
}

// Find all Package Access log by TimeRange
// @Summary Find Package Access log by TimeRange
// @Schemes
// @Description find all package access log by TimeRange
// @Produce json
// @Success 200 {array} []models.PackageAccess
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/package_accesses/period/{from}/{to} [get]
func (h *PackageAccessHandler) FindAllPackageAccessTimeRange(c *gin.Context) {
	from := c.Param("from")
	to := c.Param("to")
	fromInt, _ := strconv.ParseInt(from, 10, 64)
	toInt, _ := strconv.ParseInt(to, 10, 64)
	fromFormatted := time.Unix(fromInt, 0).Format(models.DEFAULT_TIME_FORMAT)
	toFormatted := time.Unix(toInt, 0).Format(models.DEFAULT_TIME_FORMAT)
	glList, err := h.deps.SvcOpts.PackageAccessSvc.FindAllPackageAccessTimeRange(context.Background(), fromFormatted, toFormatted)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Failed to user access",
			ErrorMsg:   err.Error(),
		})
	}
	utils.ResponseJson(c, http.StatusOK, glList)
}
