package handlers

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/ecoprohcm/DMS_BackendServer/models"
)

type HandlerOptions struct {
	AreaHandler         *AreaHandler
	CustomerHandler     *CustomerHandler
	DoorlockHandler     *DoorlockHandler
	EmployeeHandler     *EmployeeHandler
	GatewayHandler      *GatewayHandler
	LogHandler          *GatewayLogHandler
	StudentHandler      *StudentHandler
	SchedulerHandler    *SchedulerHandler
	SecretKeyHandler    *SecretKeyHandler
	UHFStatusLogHandler *UHFStatusLogHandler
	UHFHandler          *UHFHandler
	ActionHandler       *ActionHandler
	OperationLogHandler *OperationLogHandler
}

type HandlerDependencies struct {
	SvcOpts    *models.ServiceOptions
	MqttClient mqtt.Client
}
