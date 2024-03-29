package initializers

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/ecoprohcm/DMS_BackendServer/handlers"
	logger "github.com/ecoprohcm/DMS_BackendServer/logs"
	"github.com/ecoprohcm/DMS_BackendServer/models"
	"github.com/ecoprohcm/DMS_BackendServer/mqttSvc"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type ContextContainer struct {
	Config         Config
	Db             *gorm.DB
	MqttClient     mqtt.Client
	HandlerOptions *handlers.HandlerOptions
}

func ProvideConfig(envFilePath string) (Config, error) {
	cfg := Config{}
	err := godotenv.Load(envFilePath) //use env.local for localhost
	if err != nil {
		fmt.Printf("Error loading .env file %s", err)
		return Config{}, err
	}
	err = envconfig.Process("", &cfg)
	if err != nil {
		return cfg, err
	}

	logger.InitLogger(cfg.SvLogPath)
	return cfg, nil
}

func ProvideGormDb(config Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", config.DbUser, config.DbPass, config.DbHost, config.DbPort, config.DbName)
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("failed to connect database")
		return nil, err
	}
	models.Migrate(db)
	return db, nil
}

func ProvideSvcOptions(db *gorm.DB) *models.ServiceOptions {
	return &models.ServiceOptions{
		GatewaySvc:       models.NewGatewaySvc(db),
		AreaSvc:          models.NewAreaSvc(db),
		LogSvc:           models.NewLogSvc(db),
		UHFStatusLogSvc:  models.NewUHFStatusLogSvc(db),
		UHFSvc:           models.NewUHFSvc(db),
		UserAccessSvc:    models.NewUserAccessSvc(db),
		PackageAccessSvc: models.NewPackageAccessSvc(db),
		SystemLogSvc:     models.NewSystemLogSvc(db),
		OperationLogSvc:  models.NewOperationLogSvc(db),
	}
}

func ProvideMqttClient(config Config, svcOptions *models.ServiceOptions) mqtt.Client {
	return mqttSvc.MqttClient(
		config.MqttClient,
		config.MqttHost,
		config.MqttPort,
		svcOptions,
	)
}

func ProvideHandlerOptions(svcOptions *models.ServiceOptions, mqttClient mqtt.Client) *handlers.HandlerOptions {
	deps := &handlers.HandlerDependencies{
		SvcOpts:    svcOptions,
		MqttClient: mqttClient,
	}

	return &handlers.HandlerOptions{
		AreaHandler:          handlers.NewAreaHandler(deps),
		GatewayHandler:       handlers.NewGatewayHandler(deps),
		LogHandler:           handlers.NewGatewayLogHandler(deps),
		UHFStatusLogHandler:  handlers.NewUHFStatusLogHandler(deps),
		UHFHandler:           handlers.NewUHFHandler(deps),
		UserAccessHandler:    handlers.NewUserAccessHandler(deps),
		OperationLogHandler:  handlers.NewOperationLogHandler(deps),
		PackageAccessHandler: handlers.NewPackageAccessHandler(deps),
	}
}

func ProvideAppInfrastructure(config Config, db *gorm.DB, mqttClient mqtt.Client, handlerOpts *handlers.HandlerOptions) *ContextContainer {
	return &ContextContainer{
		Config:         config,
		Db:             db,
		MqttClient:     mqttClient,
		HandlerOptions: handlerOpts,
	}
}
