// Package mqttSvc provides mqtt connections, configs,
// mqtt topics, subscribe callbacks,
// mqtt error handlers, mqtt payload parsings
package mqttSvc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	_ "fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	logger "github.com/ecoprohcm/DMS_BackendServer/logs"
	"strconv"

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
	//opts.SetWill(TOPIC_SV_LASTWILL, string(`{"status":"shutdown"}`), 0, false)
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
	//topicSubscriberMap[TOPIC_GW_LOG_C] = gwLogCreateSubscriber(client, optSvc)
	topicSubscriberMap[TOPIC_GW_DOORLOCK_C] = gwDoorlockCreateSubscriber(client, optSvc)
	topicSubscriberMap[TOPIC_GW_DOORLOCK_D] = gwDoorlockDeleteSubscriber(client, optSvc)
	topicSubscriberMap[TOPIC_GW_UHF_CONNECT_STATE] = gwUHFConnectStateSubscriber(client, optSvc)
	topicSubscriberMap[TOPIC_GW_UHF_SCAN] = gwUHFScanSubscriber(client, optSvc)
	topicSubscriberMap[TOPIC_GW_TAG] = gwAccessSubscriber(client, optSvc)
	topicSubscriberMap[TOPIC_GW_LOG] = gwSystemSubscriber(client, optSvc)
	topicSubscriberMap[TOPIC_GW_LASTWILL] = gwLastWillSubscriber(client, optSvc)

	for topic, subscriber := range topicSubscriberMap {
		t := client.Subscribe(topic, 1, subscriber)
		if err := HandleMqttErr(t); err == nil {
			logger.LogfWithoutFields(logger.MQTT, logger.InfoLevel, "[MQTT-INFO] Subscribed to topic %s", topic)
		}
	}
}

func gwUHFConnectStateSubscriber(client mqtt.Client, optSvc *models.ServiceOptions) mqtt.MessageHandler {
	return func(c mqtt.Client, msg mqtt.Message) {
		var payloadStr = string(msg.Payload())
		gwId := gjson.Get(payloadStr, "gateway_id")
		uhf_address := gjson.Get(payloadStr, "message.uhf_address")
		uhf_connect_state := gjson.Get(payloadStr, "message.connection_state")
		logger.LogfWithFields(logger.MQTT, logger.InfoLevel, logger.LoggerFields{
			"GwMsg": payloadStr,
		}, "Connect state of  ID %s", gwId.String())
		uhf, error := optSvc.UHFSvc.FindUHFByAddress(context.Background(), uhf_address.String(), gwId.String())
		if error != nil {
			return
		}
		uhf.ConnectState = uhf_connect_state.String()
		optSvc.UHFSvc.UpdateUHF(context.Background(), uhf)

		new_uhf_log := &models.UHFStatusLog{}
		new_uhf_log.GatewayID = gwId.String()
		new_uhf_log.UHFAddress = uhf_address.String()
		new_uhf_log.StateType = "Connect State"
		new_uhf_log.StateValue = uhf_connect_state.String()
		optSvc.UHFStatusLogSvc.CreateUHFStatusLog(context.Background(), new_uhf_log)
		return
	}
}

func gwLastWillSubscriber(client mqtt.Client, optSvc *models.ServiceOptions) mqtt.MessageHandler {
	return func(c mqtt.Client, msg mqtt.Message) {
		var payloadStr = string(msg.Payload())
		gwId := gjson.Get(payloadStr, "gateway_id")
		logger.LogfWithoutFields(logger.MQTT, logger.DebugLevel, "Gateway ID %s has disconnected", gwId.String())
		gw, _ := optSvc.GatewaySvc.FindGatewayByMacID(context.Background(), gwId.String())
		if gw != nil {
			gw.ConnectState = "disconnected"
			_, err := optSvc.GatewaySvc.UpdateGatewayConnectState(context.Background(), gw.GatewayID, gw.ConnectState)
			if err != nil {
				logger.LogfWithoutFields(logger.MQTT, logger.ErrorLevel,
					"Update connect_state for gateway ID %s failed, err %s", gwId.String(), err.Error())
			}
			new_gateway_log := &models.GatewayLog{}
			new_gateway_log.GatewayID = gwId.String()
			new_gateway_log.StateType = "Connect State"
			new_gateway_log.StateValue = "Disconnected"
			new_gateway_log.LogTime = time.Now()
			optSvc.LogSvc.CreateGatewayLog(context.Background(), new_gateway_log)
		}
	}
}

//func gwSystemSubscriber(client mqtt.Client, optSvc *models.ServiceOptions) mqtt.MessageHandler {
//	return func(c mqtt.Client, msg mqtt.Message) {
//		var payloadStr = string(msg.Payload())
//		gwId := gjson.Get(payloadStr, "gateway_id")
//		debug_mode := gjson.Get(payloadStr, "message.debug_mode")
//
//		_, error := optSvc.GatewaySvc.FindGatewayByMacID(context.Background(), gwId.String())
//		if error != nil {
//			return
//		}
//		new_system_log := &models.SystemLog{}
//		new_system_log.GatewayID = gwId.String()
//		new_system_log.LogType = "DEBUG_MODE"
//		new_system_log.Content = debug_mode.String()
//		optSvc.SystemLogSvc.CreateSystemLog(context.Background(), new_system_log)
//		return
//	}
//}

func gwSystemSubscriber(client mqtt.Client, optSvc *models.ServiceOptions) mqtt.MessageHandler {
	return func(c mqtt.Client, msg mqtt.Message) {
		var payloadStr = string(msg.Payload())
		gwId := gjson.Get(payloadStr, "gateway_id")
		log := gjson.Get(payloadStr, "message.log").String()

		_, error := optSvc.GatewaySvc.FindGatewayByMacID(context.Background(), gwId.String())
		if error != nil {
			return
		}
		new_operation_log := &models.OperationLog{}
		new_operation_log.GatewayID = gwId.String()
		new_operation_log.Content = "DEBUG_MODE"
		new_operation_log.Content = log
		optSvc.OperationLogSvc.CreateOperationLog(context.Background(), new_operation_log)
		return
	}
}

func gwAccessSubscriber(client mqtt.Client, optSvc *models.ServiceOptions) mqtt.MessageHandler {
	return func(c mqtt.Client, msg mqtt.Message) {
		var payloadStr = string(msg.Payload())
		var ecps_string []string

		gwId := gjson.Get(payloadStr, "gateway_id")
		uhf_address := gjson.Get(payloadStr, "message.uhf_address")
		epcs := gjson.Get(payloadStr, "message.epcs")
		epcs_string := epcs.String()

		err := json.Unmarshal([]byte(epcs_string), &ecps_string)
		if err != nil {
			return
		}
		_, error := optSvc.UHFSvc.FindUHFByAddress(context.Background(), uhf_address.String(), gwId.String())
		if error != nil {
			return
		}
		for _, ecp := range ecps_string {
			var new_access = &models.UserAccess{}
			new_access.UserID = gwId.String()
			new_access.Random = uhf_address.String()
			new_access.Group = ecp
			optSvc.UserAccessSvc.CreateUserAccess(context.Background(), new_access)
		}
	}
}

func gwUHFScanSubscriber(client mqtt.Client, optSvc *models.ServiceOptions) mqtt.MessageHandler {
	return func(c mqtt.Client, msg mqtt.Message) {
		var payloadStr = string(msg.Payload())
		uhf_list := []map[string]string{}
		gwId := gjson.Get(payloadStr, "gateway_id")
		uhfs := gjson.Get(payloadStr, "message.uhfs")
		err := json.Unmarshal([]byte(uhfs.String()), &uhf_list)
		if err != nil {
			return
		}
		_, err = optSvc.GatewaySvc.FindGatewayByMacID(context.Background(), gwId.String())
		if err != nil {
			return
		}
		logger.LogfWithFields(logger.MQTT, logger.DebugLevel, logger.LoggerFields{
			"payload": payloadStr,
		}, "Gateway bootup with ID %s and uhfs %s", gwId.String(), uhfs.String())
		for _, uhf := range uhf_list {
			existing_uhf, err := optSvc.UHFSvc.FindUHFByAddress(context.Background(), uhf["uhf_address"], gwId.String())
			if err != nil {
				newUHF := &models.UHF{}
				newUHF.GatewayID = gwId.String()
				newUHF.ConnectState = uhf["connect_state"]
				newUHF.ActiveState = uhf["state"]
				newUHF.UHFAddress = uhf["uhf_address"]
				newUHF.UHFSerialNumber = uuid.New().String()
				optSvc.UHFSvc.CreateUHF(context.Background(), newUHF)
			} else {
				existing_uhf.ConnectState = uhf["connect_state"]
				existing_uhf.ActiveState = uhf["state"]
				existing_uhf.UHFAddress = uhf["uhf_address"]
				optSvc.UHFSvc.UpdateUHF(context.Background(), existing_uhf)
			}
		}
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
			newGw.ConnectState = "connected"
			newGw.SoftwareVersion = gjson.Get(payloadStr, "message.version").String()
			optSvc.GatewaySvc.CreateGateway(context.Background(), newGw)
			uhfs := []models.UHF{}
			gw := []models.GwNetwork{}
			new_gateway_log := &models.GatewayLog{}
			new_gateway_log.GatewayID = gwId.String()
			new_gateway_log.StateType = "Connect State"
			new_gateway_log.StateValue = "Connected"
			new_gateway_log.LogTime = time.Now()
			optSvc.LogSvc.CreateGatewayLog(context.Background(), new_gateway_log)
			t := client.Publish(TOPIC_SV_SYNC, 1, false, ServerBootupSystemPayload(gwId.String(), uhfs, gw))
			HandleMqttErr(t)
			return
		}
		checkGw.SoftwareVersion = gjson.Get(payloadStr, "message.version").String()
		optSvc.GatewaySvc.UpdateGateway(context.Background(), checkGw)
		uhfs := checkGw.UHFs
		gw_networks := checkGw.GwNetworks
		t := client.Publish(TOPIC_SV_SYNC, 1, false, ServerBootupSystemPayload(gwId.String(), uhfs, gw_networks))
		HandleMqttErr(t)
		new_gateway_log := &models.GatewayLog{}
		new_gateway_log.GatewayID = gwId.String()
		new_gateway_log.StateType = "Connect State"
		new_gateway_log.StateValue = "Connected"
		new_gateway_log.LogTime = time.Now()
		optSvc.LogSvc.CreateGatewayLog(context.Background(), new_gateway_log)
		return
	}
}

//func gwLogCreateSubscriber(client mqtt.Client, optSvc *models.ServiceOptions) mqtt.MessageHandler {
//	return func(c mqtt.Client, msg mqtt.Message) {
//		var payloadStr = string(msg.Payload())
//		logMsg := gjson.Get(payloadStr, "message").String()
//		gatewayId := gjson.Get(payloadStr, "gateway_id")
//		logType := gjson.Get(logMsg, "log_type")
//		content := gjson.Get(logMsg, "log_data")
//		logTime := gjson.Get(logMsg, "log_time")
//		logger.LogfWithFields(logger.MQTT, logger.DebugLevel, logger.LoggerFields{
//			"logPayload": logMsg,
//		}, "Receive gw:%s logs message", gatewayId.String())
//		logTimeInt, e := strconv.ParseInt(logTime.String(), 10, 64)
//		if e != nil {
//			fmt.Println(e.Error())
//			return
//		}
//		formatLogTime := time.Unix(logTimeInt, 0)
//		fmt.Printf(" %s: %s \n", msg.Topic(), payloadStr)
//		optSvc.LogSvc.CreateGatewayLog(context.Background(), &models.GatewayLog{
//			GatewayID: gatewayId.String(),
//			LogType:   logType.String(),
//			Content:   content.String(),
//			LogTime:   formatLogTime,
//		})
//	}
//}

//func gwDoorlockUpdateSubscriber(client mqtt.Client, optSvc *models.ServiceOptions) mqtt.MessageHandler {
//	return func(c mqtt.Client, msg mqtt.Message) {
//		var payloadStr = string(msg.Payload())
//		gatewayId := gjson.Get(payloadStr, "gateway_id").String()
//		doorStateMsg := gjson.Get(payloadStr, "message").String()
//		doorlockAddress := gjson.Get(doorStateMsg, "doorlock_address").String()
//		state := gjson.Get(doorStateMsg, "doorlock_connect_state").String()
//		lastOpenTime := gjson.Get(doorStateMsg, "last_open_time")
//		activeState := gjson.Get(doorStateMsg, "doorlock_active_state").String()
//
//		dl, _ := optSvc.DoorlockSvc.FindDoorlockByAddress(context.Background(), doorlockAddress, gatewayId)
//
//		doorID := strconv.Itoa(int(dl.ID))
//
//		if activeState != "" {
//			dl.ActiveState = activeState
//			optSvc.DoorlockSvc.UpdateDoorlock(context.Background(), dl)
//		}
//
//		if state != "" {
//			optSvc.DoorlockSvc.UpdateDoorlockByAddress(context.Background(), &models.Doorlock{
//				DoorlockAddress: doorlockAddress,
//				ConnectState:    state,
//				LastOpenTime:    uint(lastOpenTime.Uint()),
//				GatewayID:       gatewayId,
//			})
//			optSvc.DoorlockStatusLogSvc.CreateDoorlockStatusLog(context.Background(), &models.DoorlockStatusLog{
//				DoorID:     doorID,
//				StateType:  "connectState",
//				StateValue: state,
//			})
//		}
//
//		doorState := gjson.Get(doorStateMsg, "doorlock_open_state").String()
//		if doorState != "" {
//			optSvc.DoorlockSvc.UpdateDoorState(context.Background(), &models.DoorlockStatus{
//				GatewayID:       gatewayId,
//				DoorlockAddress: doorlockAddress,
//				DoorState:       doorState,
//			})
//			optSvc.DoorlockStatusLogSvc.CreateDoorlockStatusLog(context.Background(), &models.DoorlockStatusLog{
//				DoorID:     doorID,
//				StateType:  "doorState",
//				StateValue: doorState,
//			})
//		}
//
//		lockState := gjson.Get(doorStateMsg, "doorlock_lock_state").String()
//		if lockState != "" {
//			optSvc.DoorlockSvc.UpdateLockState(context.Background(), &models.DoorlockStatus{
//				GatewayID:       gatewayId,
//				DoorlockAddress: doorlockAddress,
//				LockState:       lockState,
//			})
//			optSvc.DoorlockStatusLogSvc.CreateDoorlockStatusLog(context.Background(), &models.DoorlockStatusLog{
//				DoorID:     doorID,
//				StateType:  "lockState",
//				StateValue: lockState,
//			})
//		}
//	}
//}

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

//func gwLastWillSubscriber(client mqtt.Client, optSvc *models.ServiceOptions) mqtt.MessageHandler {
//	return func(c mqtt.Client, msg mqtt.Message) {
//		var payloadStr = string(msg.Payload())
//		gwId := gjson.Get(payloadStr, "gateway_id")
//		logger.LogfWithoutFields(logger.MQTT, logger.DebugLevel, "Gateway ID %s has disconnected", gwId.String())
//		gw, _ := optSvc.GatewaySvc.FindGatewayByMacID(context.Background(), gwId.String())
//		if gw != nil {
//			gw.ConnectState = "disconnected"
//			_, err := optSvc.GatewaySvc.UpdateGatewayConnectState(context.Background(), gw.GatewayID, gw.ConnectState)
//			if err != nil {
//				logger.LogfWithoutFields(logger.MQTT, logger.ErrorLevel,
//					"Update connect_state for gateway ID %s failed, err %s", gwId.String(), err.Error())
//			}
//		}
//	}
//}

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
