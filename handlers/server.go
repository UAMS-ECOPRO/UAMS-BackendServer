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
		//v1R.GET("/doorlock/status/:id", hOpts.DoorlockHandler.GetDoorlockStatusByID)
		// v1R.GET("/doorlock/status/serial/:id", hOpts.DoorlockHandler.GetDoorlockStatusBySerialID)
		v1R.PATCH("/uhf", hOpts.UHFHandler.UpdateUHF)
		//v1R.PATCH("/doorlock/cmd", hOpts.DoorlockHandler.UpdateDoorlockCmd)
		//v1R.PATCH("/doorlock/state/cmd", hOpts.DoorlockHandler.UpdateDoorlockStateCmd)
		v1R.DELETE("/uhf", hOpts.UHFHandler.DeleteUHF)

		// Gateway log routes
		v1R.GET("/gatewayLogs", hOpts.LogHandler.FindAllGatewayLog)
		v1R.GET("/gatewayLog/:id", hOpts.LogHandler.FindGatewayLogByID)
		v1R.GET("/gatewayLogs/period/:id/date/:from/:to", hOpts.LogHandler.FindGatewayLogsByTime)
		v1R.POST("/gatewayLogs/period", hOpts.LogHandler.UpdateGatewayLogCleanPeriod)
		v1R.GET("/gatewayLogs/:id/period", hOpts.LogHandler.FindGatewayLogsTypeByTime)

		v1R.GET("/uhfLogs", hOpts.UHFStatusLogHandler.GetAllUHFStatusLogs)
		v1R.GET("/uhfLog/:uhf_address", hOpts.UHFStatusLogHandler.GetUHFStatusLogByUHFAddress)
		v1R.GET("/uhfLogs/period/:gateway_id/:uhf_address/date/:from/:to", hOpts.UHFStatusLogHandler.GetUHFStatusLogInTimeRange)
		v1R.GET("/uhfLogs/:id/period", hOpts.UHFStatusLogHandler.GetUHFStatusLogInTimeRange)
		// Secret key routes

		v1R.GET("/accesses", hOpts.ActionHandler.FindAllAccesses)
		v1R.GET("/access", hOpts.ActionHandler.FindAccess)
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
