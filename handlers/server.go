package handlers

import (
	logger "github.com/ecoprohcm/DMS_BackendServer/logs"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	hOpts *HandlerOptions,
) *gin.Engine {
	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(logger.GinLogger())
	r.Use(CORSMiddleware())
	v1R := r.Group("/v1")
	{
		// Gateway routes
		v1R.GET("/gateways", hOpts.GatewayHandler.FindAllGateway)
		v1R.GET("/gateway/:id", hOpts.GatewayHandler.FindGatewayByID)
		v1R.GET("/gateway/gateway_id/:gateway_id", hOpts.GatewayHandler.FindGatewayByGatewayID)
		v1R.PATCH("/gateway", hOpts.GatewayHandler.UpdateGateway)
		v1R.DELETE("/gateway", hOpts.GatewayHandler.DeleteGateway)

		// Area routes
		v1R.GET("/areas", hOpts.AreaHandler.FindAllArea)
		v1R.GET("/area/:id", hOpts.AreaHandler.FindAreaByID)
		v1R.POST("/area", hOpts.AreaHandler.CreateArea)
		v1R.PATCH("/area", hOpts.AreaHandler.UpdateArea)
		v1R.DELETE("/area", hOpts.AreaHandler.DeleteArea)

		// UHF routes
		v1R.GET("/uhfs", hOpts.UHFHandler.FindAllUHFs)
		v1R.GET("/uhf/:id", hOpts.UHFHandler.FindUHFByID)
		v1R.PATCH("/uhf", hOpts.UHFHandler.UpdateUHF)
		v1R.DELETE("/uhf", hOpts.UHFHandler.DeleteUHF)

		// Gateway log routes
		v1R.GET("/gateway_logs", hOpts.LogHandler.FindAllGatewayLog)
		v1R.GET("/gateway_logs/:id", hOpts.LogHandler.FindGatewayLogByID)
		v1R.GET("/gateway_logs/gateway_id/:id", hOpts.LogHandler.FindGatewayByGatewayID)
		v1R.GET("/gateway_logs/gateway_id/:id/period/:from/:to", hOpts.LogHandler.FindGatewayLogsByGatewayIDAndTime)
		v1R.GET("/gateway_logs/period/:from/:to", hOpts.LogHandler.FindGatewayLogByPeriod)
		v1R.DELETE("/gateway_logs/period/:fromTime/:toTime", hOpts.LogHandler.DeleteGatewayLogInTimeRange)

		v1R.GET("/uhf_logs", hOpts.UHFStatusLogHandler.GetAllUHFStatusLogs)
		v1R.GET("/uhf_logs/:id", hOpts.UHFStatusLogHandler.GetUHFStatusLogsByID)
		v1R.GET("/uhf_logs/gateway_id/:gateway_id/uhf_address/:uhf_address", hOpts.UHFStatusLogHandler.GetUHFStatusLogByUHFAddress)
		v1R.GET("/uhf_logs/gateway_id/:gateway_id/uhf_address/:uhf_address/period/:from/:to", hOpts.UHFStatusLogHandler.GetUHFStatusLogBYGatewayIDAndUHFAddressInTimeRange)
		v1R.GET("/uhf_logs/period/:from/:to", hOpts.UHFStatusLogHandler.GetUHFStatusLogInTimeRange)
		v1R.DELETE("/uhf_logs/period/:fromTime/:toTime", hOpts.UHFStatusLogHandler.DeleteUHFLogInTimeRange)
		// Operation log
		v1R.GET("/operation_logs", hOpts.OperationLogHandler.FindAllOperationLog)
		v1R.GET("/operation_logs/:id", hOpts.OperationLogHandler.FindAllOperationLogByID)
		v1R.GET("/operation_logs/gateway_id/:gateway_id", hOpts.OperationLogHandler.FindOperationLogByGatewayID)
		v1R.GET("/operation_logs/gateway_id/:gateway_id/period/:from/:to", hOpts.OperationLogHandler.FindOperationLogsByGatewayIDAndTime)
		v1R.GET("/operation_logs/period/:from/:to", hOpts.OperationLogHandler.FindOperationLogsByTime)
		v1R.DELETE("/operation_logs/period/:fromTime/:toTime", hOpts.OperationLogHandler.DeleteOperationLogInTimeRange)

		v1R.GET("/user_accesses", hOpts.UserAccessHandler.FindAllAccesses)
		v1R.GET("/user_accesses/user_id/:id", hOpts.UserAccessHandler.FindUserAccessByUserID)
		v1R.GET("/user_accesses/user_id/:id/area_id/:area_id", hOpts.UserAccessHandler.FindUserAccessByUserIDAndAreaID)
		v1R.GET("/user_accesses/user_id/:id/period/:from/:to", hOpts.UserAccessHandler.FindUserAccessByUserIDAndTimeRange)
		v1R.GET("/user_accesses/user_id/:id/area_id/:area_id/period/:from/:to", hOpts.UserAccessHandler.FindAllUserAccessByUserIDAndAreaIDinTimeRange)
		v1R.GET("/user_accesses/area_id/:area_id", hOpts.UserAccessHandler.FindAllUserAccessByAreaID)
		v1R.GET("/user_accesses/area_id/:area_id/period/:from/:to", hOpts.UserAccessHandler.FindAllUserAccessByAreaIDAndTimeRange)
		v1R.GET("/user_accesses/period/:from/:to", hOpts.UserAccessHandler.FindAllUserAccessTimeRange)
		v1R.DELETE("/user_accesses/period/:fromTime/:toTime", hOpts.UserAccessHandler.DeleteUserAccessTimeRange)

		v1R.GET("/package_accesses", hOpts.PackageAccessHandler.FindAllAccesses)
		v1R.GET("/package_accesses/package_id/:id", hOpts.PackageAccessHandler.FindPackageAccessByPackageID)
		v1R.GET("/package_accesses/package_id/:id/area_id/:area_id", hOpts.PackageAccessHandler.FindPackageAccessByPackageIDAndAreaID)
		v1R.GET("/package_accesses/package_id/:id/period/:from/:to", hOpts.PackageAccessHandler.FindPackageAccessByPackageIDAndTimeRange)
		v1R.GET("/package_accesses/package_id/:id/area_id/:area_id/period/:from/:to", hOpts.PackageAccessHandler.FindAllPackageAccessByPackageIDAndAreaIDinTimeRange)
		v1R.GET("/package_accesses/area_id/:area_id", hOpts.PackageAccessHandler.FindAllPackageAccessByAreaID)
		v1R.GET("/package_accesses/area_id/:area_id/period/:from/:to", hOpts.PackageAccessHandler.FindAllPackageAccessByAreaIDAndTimeRange)
		v1R.GET("/package_accesses/period/:from/:to", hOpts.PackageAccessHandler.FindAllPackageAccessTimeRange)
		v1R.DELETE("/package_accesses/period/:fromTime/:toTime", hOpts.PackageAccessHandler.DeletePackageAccessTimeRange)
	}
	return r
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Accept, Origin, Cache-Control, X-Requested-With, User-Agent, Accept-Language, Accept-Encoding")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
