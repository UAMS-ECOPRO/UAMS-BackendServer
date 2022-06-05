// Package mqttSvc provides mqtt connections, configs,
// mqtt topics, subscribe callbacks,
// mqtt error handlers, mqtt payload parsings
package mqttSvc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"strconv"
	"time"

	logger "github.com/ecoprohcm/DMS_BackendServer/logs"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/ecoprohcm/DMS_BackendServer/models"
	"github.com/google/uuid"
	"github.com/tidwall/gjson"
)

func NewTlsConfig() *tls.Config {
	certpool := x509.NewCertPool()
	wd, _ := os.Getwd()
	ca, err := ioutil.ReadFile(filepath.Join(wd, "certs", "ca.pem"))
	if err != nil {
		logger.LogWithoutFields(logger.MQTT, logger.FatalLevel, err.Error())
	}
	certpool.AppendCertsFromPEM(ca)
	return &tls.Config{
		RootCAs: certpool,
	}
}

// TODO: Guarantee mqtt req/res
// var DoorlockStateCheck = make(chan bool)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	logger.LogfWithoutFields(logger.MQTT, logger.DebugLevel,
		"Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	logger.LogWithoutFields(logger.MQTT, logger.DebugLevel, "Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	logger.LogfWithoutFields(logger.MQTT, logger.ErrorLevel, "Connect lost: %v\n", err)
}

// Define mqtt connections and configs
func MqttClient(
	clientID string,
	host string,
	port string,
	optSvc *models.ServiceOptions,
) mqtt.Client {

	mqtt.ERROR = logger.NewMqttLogger("MQTT ERROR", logger.ErrorLevel)
	mqtt.CRITICAL = logger.NewMqttLogger("MQTT CRITICAL", logger.FatalLevel)
	mqtt.WARN = logger.NewMqttLogger("MQTT WARNING", logger.WarnLevel)
	//mqtt.DEBUG = logger.NewMqttLogger("[MQTT-DEBUG]", logger.DebugLevel)

	opts := mqtt.NewClientOptions()
	// Setup server LWT message

	opts.AddBroker(fmt.Sprintf("ssl://%s:%s", host, port))
	opts.SetClientID(clientID) // Need to be unique per client
	tlsConfig := NewTlsConfig()
	opts.SetTLSConfig(tlsConfig)
	// opts.SetUsername("emqx") // Use this when we want to improve security
	// opts.SetPassword("public") // Use this when we want to improve security
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		logger.LogWithoutFields(logger.MQTT, logger.PanicLevel, token.Error())
	}
	subGateway(client, optSvc)

	return client
}

type GatewaySubscriber = mqtt.MessageHandler

// Define all subscribe logic callbacks for payloads that received from gateway
func subGateway(client mqtt.Client, optSvc *models.ServiceOptions) {

	topicSubscriberMap := map[string]GatewaySubscriber{}
	topicSubscriberMap[TOPIC_GW_SHUTDOWN] = gwShutDownSubscriber(client, optSvc)
	topicSubscriberMap[TOPIC_GW_BOOTUP] = gwBootupSubscriber(client, optSvc)
	topicSubscriberMap[TOPIC_GW_LOG_C] = gwLogCreateSubscriber(client, optSvc)
	topicSubscriberMap[TOPIC_GW_DOORLOCK_U] = gwDoorlockUpdateSubscriber(client, optSvc)
	topicSubscriberMap[TOPIC_GW_DOORLOCK_C] = gwDoorlockCreateSubscriber(client, optSvc)
	topicSubscriberMap[TOPIC_GW_DOORLOCK_D] = gwDoorlockDeleteSubscriber(client, optSvc)
	topicSubscriberMap[TOPIC_GW_LASTWILL] = gwLastWillSubscriber(client, optSvc)
	topicSubscriberMap[TOPIC_GW_UHF_CONNECT_STATE] = gwUHFConnectStateSubscriber(client, optSvc)

	for topic, subscriber := range topicSubscriberMap {
		t := client.Subscribe(topic, 1, subscriber)
		if err := HandleMqttErr(t); err == nil {
			logger.LogfWithoutFields(logger.MQTT, logger.InfoLevel, "[MQTT-INFO] Subscribed to topic %s", topic)
		}
	}
}

// MQTT subscriber for gateway
func gwUHFConnectStateSubscriber(client mqtt.Client, optSvc *models.ServiceOptions) mqtt.MessageHandler {
	return func(c mqtt.Client, msg mqtt.Message) {
		var payloadStr = string(msg.Payload())
		gwId := gjson.Get(payloadStr, "gateway_id")
		uhf_address := gjson.Get(payloadStr, "message.uhf_address")
		uhf_connect_state := gjson.Get(payloadStr, "message.uhf_connect_state")
		//logger.LogfWithFields(logger.MQTT, logger.InfoLevel, logger.LoggerFields{
		//	"GwMsg": gwMsg.String(),
		//}, "Receive gateway shutdown message with ID %s", gwId.Strin
		uhf, error := optSvc.UHFSvc.FindUHFByAddress(context.Background(), uhf_address.String(), gwId.String())
		if error != nil {
			return
		}
		uhf.ConnectState = uhf_connect_state.String()
		optSvc.UHFSvc.UpdateUHF(context.Background(), uhf)
		return
	}
}

func gwShutDownSubscriber(client mqtt.Client, optSvc *models.ServiceOptions) mqtt.MessageHandler {
	return func(c mqtt.Client, msg mqtt.Message) {
		var payloadStr = string(msg.Payload())
		gwId := gjson.Get(payloadStr, "gateway_id")
		gwMsg := gjson.Get(payloadStr, "message")
		logger.LogfWithFields(logger.MQTT, logger.InfoLevel, logger.LoggerFields{
			"GwMsg": gwMsg.String(),
		}, "Receive gateway shutdown message with ID %s", gwId.String())
		optSvc.GatewaySvc.DeleteGateway(context.Background(), gwId.String())
	}
}

func gwBootupSubscriber(client mqtt.Client, optSvc *models.ServiceOptions) mqtt.MessageHandler {
	return func(c mqtt.Client, msg mqtt.Message) {
		var payloadStr = string(msg.Payload())
		gwId := gjson.Get(payloadStr, "gateway_id")
		gw_string := gwId.String()
		logger.LogfWithFields(logger.MQTT, logger.DebugLevel, logger.LoggerFields{
			"payload": payloadStr,
		}, "Gateway bootup with ID %s", gw_string)

		checkGw, _ := optSvc.GatewaySvc.FindGatewayByMacID(context.Background(), gwId.String())

		if checkGw == nil {
			newGw := &models.Gateway{}
			newGw.GatewayID = gwId.String()
			newGw.ConnectState = true
			newGw.SoftwareVersion = gjson.Get(payloadStr, "message.software_version").String()
			optSvc.GatewaySvc.CreateGateway(context.Background(), newGw)
			uhfs := []models.UHF{}
			gw := []models.GwNetwork{}
			t := client.Publish(TOPIC_SV_SYNC, 1, false, ServerBootupSystemPayload(gwId.String(), uhfs, gw))
			HandleMqttErr(t)
			return
		}
		checkGw.SoftwareVersion = gjson.Get(payloadStr, "message.software_version").String()
		optSvc.GatewaySvc.UpdateGateway(context.Background(), checkGw)
		uhfs := checkGw.UHFs
		gw_networks := checkGw.GwNetworks
		t := client.Publish(TOPIC_SV_SYNC, 1, false, ServerBootupSystemPayload(gwId.String(), uhfs, gw_networks))
		HandleMqttErr(t)
	}
}

func gwLogCreateSubscriber(client mqtt.Client, optSvc *models.ServiceOptions) mqtt.MessageHandler {
	return func(c mqtt.Client, msg mqtt.Message) {
		var payloadStr = string(msg.Payload())
		logMsg := gjson.Get(payloadStr, "message").String()
		gatewayId := gjson.Get(payloadStr, "gateway_id")
		logType := gjson.Get(logMsg, "log_type")
		content := gjson.Get(logMsg, "log_data")
		logTime := gjson.Get(logMsg, "log_time")
		logger.LogfWithFields(logger.MQTT, logger.DebugLevel, logger.LoggerFields{
			"logPayload": logMsg,
		}, "Receive gw:%s logs message", gatewayId.String())
		logTimeInt, e := strconv.ParseInt(logTime.String(), 10, 64)
		if e != nil {
			fmt.Println(e.Error())
			return
		}
		formatLogTime := time.Unix(logTimeInt, 0)
		fmt.Printf(" %s: %s \n", msg.Topic(), payloadStr)
		optSvc.LogSvc.CreateGatewayLog(context.Background(), &models.GatewayLog{
			GatewayID: gatewayId.String(),
			LogType:   logType.String(),
			Content:   content.String(),
			LogTime:   formatLogTime,
		})
	}
}

func gwDoorlockUpdateSubscriber(client mqtt.Client, optSvc *models.ServiceOptions) mqtt.MessageHandler {
	return func(c mqtt.Client, msg mqtt.Message) {
		var payloadStr = string(msg.Payload())
		gatewayId := gjson.Get(payloadStr, "gateway_id").String()
		doorStateMsg := gjson.Get(payloadStr, "message").String()
		doorlockAddress := gjson.Get(doorStateMsg, "doorlock_address").String()
		state := gjson.Get(doorStateMsg, "doorlock_connect_state").String()
		lastOpenTime := gjson.Get(doorStateMsg, "last_open_time")
		activeState := gjson.Get(doorStateMsg, "doorlock_active_state").String()

		dl, _ := optSvc.DoorlockSvc.FindDoorlockByAddress(context.Background(), doorlockAddress, gatewayId)

		doorID := strconv.Itoa(int(dl.ID))

		if activeState != "" {
			dl.ActiveState = activeState
			optSvc.DoorlockSvc.UpdateDoorlock(context.Background(), dl)
		}

		if state != "" {
			optSvc.DoorlockSvc.UpdateDoorlockByAddress(context.Background(), &models.Doorlock{
				DoorlockAddress: doorlockAddress,
				ConnectState:    state,
				LastOpenTime:    uint(lastOpenTime.Uint()),
				GatewayID:       gatewayId,
			})
			optSvc.DoorlockStatusLogSvc.CreateDoorlockStatusLog(context.Background(), &models.DoorlockStatusLog{
				DoorID:     doorID,
				StateType:  "connectState",
				StateValue: state,
			})
		}

		doorState := gjson.Get(doorStateMsg, "doorlock_open_state").String()
		if doorState != "" {
			optSvc.DoorlockSvc.UpdateDoorState(context.Background(), &models.DoorlockStatus{
				GatewayID:       gatewayId,
				DoorlockAddress: doorlockAddress,
				DoorState:       doorState,
			})
			optSvc.DoorlockStatusLogSvc.CreateDoorlockStatusLog(context.Background(), &models.DoorlockStatusLog{
				DoorID:     doorID,
				StateType:  "doorState",
				StateValue: doorState,
			})
		}

		lockState := gjson.Get(doorStateMsg, "doorlock_lock_state").String()
		if lockState != "" {
			optSvc.DoorlockSvc.UpdateLockState(context.Background(), &models.DoorlockStatus{
				GatewayID:       gatewayId,
				DoorlockAddress: doorlockAddress,
				LockState:       lockState,
			})
			optSvc.DoorlockStatusLogSvc.CreateDoorlockStatusLog(context.Background(), &models.DoorlockStatusLog{
				DoorID:     doorID,
				StateType:  "lockState",
				StateValue: lockState,
			})
		}
	}
}

func gwDoorlockCreateSubscriber(client mqtt.Client, optSvc *models.ServiceOptions) mqtt.MessageHandler {
	return func(c mqtt.Client, msg mqtt.Message) {
		dl := parseDoorlockPayload(string(msg.Payload()))
		optSvc.DoorlockSvc.CreateDoorlock(context.Background(), dl)
	}
}

func gwDoorlockDeleteSubscriber(client mqtt.Client, optSvc *models.ServiceOptions) mqtt.MessageHandler {
	return func(c mqtt.Client, msg mqtt.Message) {
		var payloadStr = string(msg.Payload())
		gatewayId := gjson.Get(payloadStr, "gateway_id").String()
		doorStateMsg := gjson.Get(payloadStr, "message").String()
		doorlockAddress := gjson.Get(doorStateMsg, "doorlock_address").String()
		optSvc.DoorlockSvc.DeleteDoorlockByAddress(context.Background(), &models.Doorlock{
			DoorlockAddress: doorlockAddress,
			GatewayID:       gatewayId,
		})
	}
}

func gwLastWillSubscriber(client mqtt.Client, optSvc *models.ServiceOptions) mqtt.MessageHandler {
	return func(c mqtt.Client, msg mqtt.Message) {
		var payloadStr = string(msg.Payload())
		gwId := gjson.Get(payloadStr, "gateway_id")
		logger.LogfWithoutFields(logger.MQTT, logger.DebugLevel, "Gateway ID %s has disconnected", gwId.String())
		gw, _ := optSvc.GatewaySvc.FindGatewayByMacID(context.Background(), gwId.String())
		if gw != nil {
			gw.ConnectState = false
			_, err := optSvc.GatewaySvc.UpdateGatewayConnectState(context.Background(), gw.GatewayID, gw.ConnectState)
			if err != nil {
				logger.LogfWithoutFields(logger.MQTT, logger.ErrorLevel,
					"Update connect_state for gateway ID %s failed, err %s", gwId.String(), err.Error())
			}
		}
	}
}

// Util funcs
func parseDoorlockPayload(payloadStr string) *models.Doorlock {
	doorStateMsg := gjson.Get(payloadStr, "message").String()
	doorlockAdress := gjson.Get(doorStateMsg, "doorlock_address")
	active_state := gjson.Get(doorStateMsg, "doorlock_active_state")
	gatewayId := gjson.Get(payloadStr, "gateway_id")
	open_state := gjson.Get(doorStateMsg, "doorlock_open_state")
	lock_state := gjson.Get(doorStateMsg, "doorlock_lock_state")
	doorSerialId := uuid.New().String()

	dl := &models.Doorlock{
		GatewayID:       gatewayId.String(),
		DoorSerialID:    doorSerialId,
		DoorlockAddress: doorlockAdress.String(),
		ActiveState:     active_state.String(),
		DoorState:       open_state.String(),
		LockState:       lock_state.String(),
	}
	return dl
}

func getUserPassInfoFromScheduler(optSvc *models.ServiceOptions, sche models.Scheduler) (userIdPwd UserIDPassword, err bool) {

	err = false
	if sche.Role == "employee" {
		emp, _ := optSvc.EmployeeSvc.FindEmployeeByMSNV(context.Background(), sche.UserID)
		userIdPwd.UserId = emp.MSNV
		userIdPwd.RfidPass = emp.RfidPass
		userIdPwd.KeypadPass = emp.KeypadPass
		err = true
	} else if sche.Role == "student" {
		stu, _ := optSvc.StudentSvc.FindStudentByMSSV(context.Background(), sche.UserID)
		userIdPwd.UserId = stu.MSSV
		userIdPwd.RfidPass = stu.RfidPass
		userIdPwd.KeypadPass = stu.KeypadPass
		err = true
	} else if sche.Role == "customer" {
		cus, _ := optSvc.CustomerSvc.FindCustomerByCCCD(context.Background(), sche.UserID)
		userIdPwd.UserId = cus.CCCD
		userIdPwd.RfidPass = cus.RfidPass
		userIdPwd.KeypadPass = cus.KeypadPass
		err = true
	}

	return userIdPwd, err
}

func mergeInfoToScheBootUp(optSvc *models.ServiceOptions, dlList []models.Doorlock) (scheBoUpList []*SchedulerBootUp) {
	for _, dl := range dlList {
		for _, sche := range dl.Schedulers {

			uIp, _ := getUserPassInfoFromScheduler(optSvc, sche)

			scheBoUp := SchedulerBootUp{
				SchedulerId:     strconv.Itoa(int(sche.ID)),
				UserId:          uIp.UserId,
				RfidPass:        uIp.RfidPass,
				KeypadPass:      uIp.KeypadPass,
				DoorlockAddress: dl.DoorlockAddress,
				StartDate:       sche.StartDate,
				EndDate:         sche.EndDate,
				WeekDay:         strconv.Itoa(int(sche.WeekDay)),
				StartClass:      strconv.Itoa(int(sche.StartClassTime)),
				EndClass:        strconv.Itoa(int(sche.EndClassTime)),
			}
			scheBoUpList = append(scheBoUpList, &scheBoUp)
		}
	}
	return scheBoUpList
}
