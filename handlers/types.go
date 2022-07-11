package handlers

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/ecoprohcm/DMS_BackendServer/models"
)

type HandlerOptions struct {
	AreaHandler          *AreaHandler
	GatewayHandler       *GatewayHandler
	LogHandler           *GatewayLogHandler
	UHFStatusLogHandler  *UHFStatusLogHandler
	UHFHandler           *UHFHandler
	UserAccessHandler    *UserAccessHandler
	PackageAccessHandler *PackageAccessHandler
	OperationLogHandler  *OperationLogHandler
}

type HandlerDependencies struct {
	SvcOpts    *models.ServiceOptions
	MqttClient mqtt.Client
}
