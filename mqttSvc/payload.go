package mqttSvc

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ecoprohcm/DMS_BackendServer/models"
)

type UserIDPassword struct {
	UserId     string `json:"user_id"`
	RfidPass   string `json:"rfid_pw"`
	KeypadPass string `json:"keypad_pw"`
}
type UHFBootUp struct {
	UHFAddress  string `json:"uhf_address"`
	ActiveState string `json:"uhf_active_state"`
}

type SchedulerBootUp struct {
	SchedulerId     string `json:"register_id"`
	UserId          string `json:"user_id"`
	RfidPass        string `json:"rfid_pw"`
	KeypadPass      string `json:"keypad_pw"`
	DoorlockAddress string `json:"doorlock_address"`
	StartDate       string `json:"start_date"`
	EndDate         string `json:"end_date"`
	WeekDay         string `json:"week_day"`
	StartClass      string `json:"start_class"`
	EndClass        string `json:"end_class"`
}

type UHFSyncPayload struct {
	UHFAddress   string `json:"uhf_address"`
	ConnectState string `json:"connect_state"`
	State        string `json:"state"`
}

type SyncPayload struct {
	UHFs  []UHFSyncPayload `json:"uhfs"`
	State string           `json:"state"`
}

func ServerUpdateUHFPayload(uhf *models.UHF) string {
	msg := fmt.Sprintf(`{"address":"%s","state":"%s"}`,
		uhf.UHFAddress, uhf.ActiveState)
	return PayloadWithGatewayId(uhf.GatewayID, msg)
}

func ServerDeleteUHFPayload(uhf *models.UHF) string {
	msg := fmt.Sprintf(`{"address":"%s"}`, uhf.UHFAddress)
	return PayloadWithGatewayId(uhf.GatewayID, msg)
}

func ServerUpdateGatewayPayload(gw *models.Gateway) string {
	msg := fmt.Sprintf(`{"state":"%s"}`, gw.ConnectState)
	return PayloadWithGatewayId(gw.GatewayID, msg)
}

func ServerDeleteGatewayPayload(gwID string) string {
	msg := `{}`
	return PayloadWithGatewayId(gwID, msg)
}

func PayloadWithGatewayId(gwId string, msg string) string {
	return fmt.Sprintf(`{"gateway_id":"%s","message":%s}`, gwId, msg)
}

func getDayMonthYearSlice(str string) []int {
	strs := strings.Split(str, "/")
	var dmySlice = []int{}
	for _, s := range strs {
		number, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			panic(err)
		}
		dmySlice = append(dmySlice, int(number))
	}
	return dmySlice
}

func ServerUpdateSecretKeyPayload(gwId string, secretKey string) string {
	msg := fmt.Sprintf(`{"secret_key":"%s"}`, secretKey)
	return PayloadWithGatewayId(gwId, msg)
}

func ServerUpdateGatewayCmd(gwId string, action string) string {
	msg := fmt.Sprintf(`{"action":"%s"}`, action)
	return PayloadWithGatewayId(gwId, msg)
}

func ServerBootupUHFsPayload(gwId string, dls []models.UHF) string {
	bootupDls := []UHFBootUp{}
	for _, dl := range dls {
		buDl := UHFBootUp{
			UHFAddress:  dl.UHFAddress,
			ActiveState: dl.ActiveState,
		}
		bootupDls = append(bootupDls, buDl)
	}
	bootupDlsJson, _ := json.Marshal(bootupDls)
	return PayloadWithGatewayId(gwId, string(bootupDlsJson))
}

func ServerBootupRegisterPayload(
	gwId string,
	scheBoUpListPointer []*SchedulerBootUp,
) string {
	scheBoUpList := []SchedulerBootUp{}
	for _, sche := range scheBoUpListPointer {
		loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
		startDmySlice := getDayMonthYearSlice(sche.StartDate)
		start := time.Date(startDmySlice[2], time.Month(startDmySlice[1]), startDmySlice[0], 0, 0, 0, 0, loc).Unix()
		sche.StartDate = strconv.FormatInt(start, 10)
		endDmySlice := getDayMonthYearSlice(sche.EndDate)
		end := time.Date(endDmySlice[2], time.Month(endDmySlice[1]), endDmySlice[0], 23, 59, 59, 0, loc).Unix()
		sche.EndDate = strconv.FormatInt(end, 10)

		if !isPastTime(end) {
			scheBoUpList = append(scheBoUpList, *sche)
		}
	}
	bootupScheJson, _ := json.Marshal(scheBoUpList)
	return PayloadWithGatewayId(gwId, string(bootupScheJson))
}

func isPastTime(t_compared int64) bool {

	t_now := time.Now().Unix()
	if t_compared >= t_now {
		return false
	}

	return true
}

func ServerBootupSystemPayload(gwId string, uhfs []models.UHF) string {
	uhf_important_info := []UHFSyncPayload{}
	for _, item := range uhfs {
		new_uhf_important_info := UHFSyncPayload{}
		new_uhf_important_info.UHFAddress = item.UHFAddress
		new_uhf_important_info.ConnectState = item.ConnectState
		new_uhf_important_info.State = item.ActiveState
		uhf_important_info = append(uhf_important_info, new_uhf_important_info)
	}
	sync_payload := SyncPayload{UHFs: uhf_important_info, State: "active"}
	sync_payload_converted, _ := json.Marshal(sync_payload)
	return PayloadWithGatewayId(gwId, string(sync_payload_converted))
}
