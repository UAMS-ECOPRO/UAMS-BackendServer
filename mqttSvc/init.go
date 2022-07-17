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
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	logger "github.com/ecoprohcm/DMS_BackendServer/logs"
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

type UHFTagInfo struct {
	EPC       string `json:"epc"`
	Mem       string `json:"mem"`
	TimeStamp string `json:"timestamp"`
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
	topicSubscriberMap[TOPIC_GW_UHF_CONNECT_STATE] = gwUHFConnectStateSubscriber(client, optSvc)
	topicSubscriberMap[TOPIC_GW_UHF_SCAN] = gwUHFScanSubscriber(client, optSvc)
	topicSubscriberMap[TOPIC_GW_TAG] = gwAccessSubscriber(client, optSvc)
	topicSubscriberMap[TOPIC_GW_LOG] = gwSystemSubscriber(client, optSvc)
	topicSubscriberMap[TOPIC_GW_LASTWILL] = gwLastWillSubscriber(client, optSvc)
	topicSubscriberMap[TOPIC_GW_GW_CONNECT_STATE] = gwGatewayConnectStateSubscriber(client, optSvc)

	for topic, subscriber := range topicSubscriberMap {
		t := client.Subscribe(topic, 1, subscriber)
		if err := HandleMqttErr(t); err == nil {
			logger.LogfWithoutFields(logger.MQTT, logger.InfoLevel, "[MQTT-INFO] Subscribed to topic %s", topic)
		}
	}
}

func gwGatewayConnectStateSubscriber(client mqtt.Client, optSvc *models.ServiceOptions) mqtt.MessageHandler {
	return func(c mqtt.Client, msg mqtt.Message) {
		var payloadStr = string(msg.Payload())
		gwId := gjson.Get(payloadStr, "gateway_id")
		gw_connect_state := gjson.Get(payloadStr, "message.connection_state")
		logger.LogfWithFields(logger.MQTT, logger.InfoLevel, logger.LoggerFields{
			"GwMsg": payloadStr,
		}, "Connect state of  ID %s", gwId.String())
		gw, error := optSvc.GatewaySvc.FindGatewayByGatewayID(context.Background(), gwId.String())
		if error != nil {
			return
		}
		gw.ConnectState = gw_connect_state.String()
		optSvc.GatewaySvc.UpdateGateway(context.Background(), gw)
		new_gw_log := &models.GatewayLog{}
		new_gw_log.GatewayID = gwId.String()
		new_gw_log.StateType = "Connect State"
		new_gw_log.StateValue = gw_connect_state.String()
		new_gw_log.LogTime = time.Now()
		optSvc.LogSvc.CreateGatewayLog(context.Background(), new_gw_log)
		return
	}
}

func gwUHFConnectStateSubscriber(client mqtt.Client, optSvc *models.ServiceOptions) mqtt.MessageHandler {
	return func(c mqtt.Client, msg mqtt.Message) {
		var payloadStr = string(msg.Payload())
		gwId := gjson.Get(payloadStr, "gateway_id")
		uhf_address := gjson.Get(payloadStr, "message.address")
		uhf_connect_state := gjson.Get(payloadStr, "message.connection_state")
		time_stamp := gjson.Get(payloadStr, "message.timestamp")
		logger.LogfWithFields(logger.MQTT, logger.InfoLevel, logger.LoggerFields{
			"GwMsg": payloadStr,
		}, "Connect state of  ID %s", gwId.String())
		uhf, error := optSvc.UHFSvc.FindUHFByAddress(context.Background(), uhf_address.String(), gwId.String())
		if error != nil {
			return
		}
		uhf.ConnectState = uhf_connect_state.String()
		optSvc.UHFSvc.UpdateUHF(context.Background(), uhf)
		time_layout := "2006-01-02 15:04:05"
		var time_stamp_converted, _ = time.ParseInLocation(time_layout, time_stamp.String(), time.Local)
		new_uhf_log := &models.UHFStatusLog{}
		new_uhf_log.GatewayID = gwId.String()
		new_uhf_log.UHFAddress = uhf_address.String()
		new_uhf_log.StateType = "Connect State"
		new_uhf_log.StateValue = uhf_connect_state.String()
		new_uhf_log.Time = time_stamp_converted
		optSvc.UHFStatusLogSvc.CreateUHFStatusLog(context.Background(), new_uhf_log)
		return
	}
}

func gwLastWillSubscriber(client mqtt.Client, optSvc *models.ServiceOptions) mqtt.MessageHandler {
	return func(c mqtt.Client, msg mqtt.Message) {
		var payloadStr = string(msg.Payload())
		gwId := gjson.Get(payloadStr, "gateway_id")
		logger.LogfWithoutFields(logger.MQTT, logger.DebugLevel, "Gateway ID %s has disconnect", gwId.String())
		gw, _ := optSvc.GatewaySvc.FindGatewayByGatewayID(context.Background(), gwId.String())
		if gw != nil {
			gw.ConnectState = "disconnect"
			_, err := optSvc.GatewaySvc.UpdateGatewayConnectState(context.Background(), gw.GatewayID, gw.ConnectState)
			if err != nil {
				logger.LogfWithoutFields(logger.MQTT, logger.ErrorLevel,
					"Update connect_state for gateway ID %s failed, err %s", gwId.String(), err.Error())
			}
			new_gateway_log := &models.GatewayLog{}
			new_gateway_log.GatewayID = gwId.String()
			new_gateway_log.StateType = "Connect State"
			new_gateway_log.StateValue = "disconnect"
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
		time_layout := "2006-01-02 15:04:05"
		var time_stamp, _ = time.ParseInLocation(time_layout, gjson.Get(payloadStr, "message.timestamp").String(), time.Local)
		_, error := optSvc.GatewaySvc.FindGatewayByGatewayID(context.Background(), gwId.String())
		if error != nil {
			return
		}
		new_operation_log := &models.OperationLog{}
		new_operation_log.GatewayID = gwId.String()
		new_operation_log.Content = "DEBUG_MODE"
		new_operation_log.Content = log
		new_operation_log.Time = time_stamp
		optSvc.OperationLogSvc.CreateOperationLog(context.Background(), new_operation_log)
		return
	}
}

func gwAccessSubscriber(client mqtt.Client, optSvc *models.ServiceOptions) mqtt.MessageHandler {
	return func(c mqtt.Client, msg mqtt.Message) {
		var payloadStr = string(msg.Payload())
		var tags_list []UHFTagInfo

		gwId := gjson.Get(payloadStr, "gateway_id")
		uhf_address := gjson.Get(payloadStr, "message.address")
		tags := gjson.Get(payloadStr, "message.tags")
		tags_string := tags.String()
		time_layout := "2006-01-02 15:04:05"

		err := json.Unmarshal([]byte(tags_string), &tags_list)
		if err != nil {
			return
		}
		existing_uhf, error := optSvc.UHFSvc.FindUHFByAddress(context.Background(), uhf_address.String(), gwId.String())
		if error != nil {
			return
		}
		if existing_uhf.AreaId == "" {
			return
		}
		for _, item := range tags_list {
			var mem = item.Mem
			var time_stamp, _ = time.ParseInLocation(time_layout, item.TimeStamp, time.Local)
			var access_type = mem[0:1]
			var access_group = mem[1:3]
			var access_id = mem[3:13]
			var access_random = mem[13:16]
			if access_type == "U" {
				var new_user_access = &models.UserAccess{}
				new_user_access.UserID = access_id
				new_user_access.Random = access_random
				new_user_access.Group = access_group
				new_user_access.AreaID = existing_uhf.AreaId
				new_user_access.Time = time_stamp
				optSvc.UserAccessSvc.CreateUserAccess(context.Background(), new_user_access)
			} else if access_type == "P" {
				var new_package_access = &models.PackageAccess{}
				new_package_access.PackageID = access_id
				new_package_access.Random = access_random
				new_package_access.Group = access_group
				new_package_access.AreaID = existing_uhf.AreaId
				new_package_access.Time = time_stamp
				optSvc.PackageAccessSvc.CreatePackageAccess(context.Background(), new_package_access)
			}

		}
		return
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
		_, err = optSvc.GatewaySvc.FindGatewayByGatewayID(context.Background(), gwId.String())
		if err != nil {
			return
		}
		logger.LogfWithFields(logger.MQTT, logger.DebugLevel, logger.LoggerFields{
			"payload": payloadStr,
		}, "Gateway bootup with ID %s and uhfs %s", gwId.String(), uhfs.String())
		for _, uhf := range uhf_list {
			_, err := optSvc.UHFSvc.FindUHFByAddress(context.Background(), uhf["address"], gwId.String())
			if err != nil {
				newUHF := &models.UHF{}
				newUHF.GatewayID = gwId.String()
				newUHF.ConnectState = "connect"
				newUHF.ActiveState = "inactive"
				newUHF.UHFAddress = uhf["address"]
				newUHF.UHFSerialNumber = uuid.New().String()
				optSvc.UHFSvc.CreateUHF(context.Background(), newUHF)
			}
			//else {
			//	//existing_uhf.ConnectState = uhf["connect_state"]
			//	//existing_uhf.ActiveState = uhf["state"]
			//	existing_uhf.UHFAddress = uhf["uhf_address"]
			//	optSvc.UHFSvc.UpdateUHF(context.Background(), existing_uhf)
			//}
		}
		checkGw, _ := optSvc.GatewaySvc.FindGatewayByGatewayID(context.Background(), gwId.String())
		uhf_in_db_list := checkGw.UHFs
		for _, uhf_in_db := range uhf_in_db_list {
			if check_exist(uhf_in_db, uhf_list) == false {
				optSvc.UHFSvc.DeleteUHF(context.Background(), strconv.FormatUint(uint64(uhf_in_db.ID), 10))
			}
		}
		checkGw_again, _ := optSvc.GatewaySvc.FindGatewayByGatewayID(context.Background(), gwId.String())
		t := client.Publish(TOPIC_SV_SYNC, 1, false, ServerBootupSystemPayload(gwId.String(), checkGw_again.UHFs))
		HandleMqttErr(t)
		return
	}
}

func check_exist(uhf_to_check models.UHF, list_to_check []map[string]string) bool {
	for _, uhf := range list_to_check {
		if uhf_to_check.UHFAddress == uhf["address"] {
			return true
		}
	}
	return false
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

		checkGw, _ := optSvc.GatewaySvc.FindGatewayByGatewayID(context.Background(), gwId.String())

		if checkGw == nil {
			newGw := &models.Gateway{}
			newGw.GatewayID = gwId.String()
			newGw.ConnectState = "connect"
			newGw.SoftwareVersion = gjson.Get(payloadStr, "message.version").String()
			optSvc.GatewaySvc.CreateGateway(context.Background(), newGw)
			uhfs := []models.UHF{}
			new_gateway_log := &models.GatewayLog{}
			new_gateway_log.GatewayID = gwId.String()
			new_gateway_log.StateType = "Connect State"
			new_gateway_log.StateValue = "connect"
			new_gateway_log.LogTime = time.Now()
			optSvc.LogSvc.CreateGatewayLog(context.Background(), new_gateway_log)
			t := client.Publish(TOPIC_SV_SYNC, 1, false, ServerBootupSystemPayload(gwId.String(), uhfs))
			HandleMqttErr(t)
			return
		}
		checkGw.SoftwareVersion = gjson.Get(payloadStr, "message.version").String()
		optSvc.GatewaySvc.UpdateGateway(context.Background(), checkGw)
		uhfs := checkGw.UHFs
		t := client.Publish(TOPIC_SV_SYNC, 1, false, ServerBootupSystemPayload(gwId.String(), uhfs))
		HandleMqttErr(t)
		new_gateway_log := &models.GatewayLog{}
		new_gateway_log.GatewayID = gwId.String()
		new_gateway_log.StateType = "Connect State"
		new_gateway_log.StateValue = "connect"
		new_gateway_log.LogTime = time.Now()
		optSvc.LogSvc.CreateGatewayLog(context.Background(), new_gateway_log)
		return
	}
}
