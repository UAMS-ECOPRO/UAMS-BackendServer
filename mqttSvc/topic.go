package mqttSvc

const (
	TOPIC_GW_LOG_C           string = "uams/gateway/log/create"
	TOPIC_GW_DOORLOCK_STATUS string = "uams/gateway/doorlock/status"
	TOPIC_GW_DOORLOCK_C      string = "uams/gateway/doorlock/create"
	TOPIC_GW_DOORLOCK_U      string = "uams/gateway/doorlock/update"
	TOPIC_GW_DOORLOCK_D      string = "uams/gateway/doorlock/delete"

	TOPIC_GW_BOOTUP            string = "uams/gateway/bootup"
	TOPIC_GW_UHF_CONNECT_STATE string = "uams/gateway/uhf/update"
	TOPIC_GW_UHF_SCAN          string = "uams/gateway/uhf/scan"
	TOPIC_GW_SHUTDOWN          string = "uams/gateway/shutdown"
	TOPIC_GW_LASTWILL          string = "uams/gateway/lastwill"
	TOPIC_GW_TAG               string = "uams/gateway/uhf/tag"
	TOPIC_GW_LOG               string = "uams/gateway/log"
	TOPIC_GW_GW_CONNECT_STATE  string = "uams/gateway/gateway/update"

	TOPIC_SV_DOORLOCK_C   string = "uams/server/uhf/create"
	TOPIC_SV_UHF_U        string = "uams/server/uhf/update"
	TOPIC_SV_UHF_D        string = "uams/server/uhf/delete"
	TOPIC_SV_DOORLOCK_CMD string = "server/doorlock/command"

	TOPIC_SV_GATEWAY_U string = "uams/server/gateway/update"
	TOPIC_SV_GATEWAY_D string = "uams/server/gateway/delete"

	TOPIC_SV_SCHEDULER_C      string = "uams/server/register/create"
	TOPIC_SV_SCHEDULER_U      string = "uams/server/register/update"
	TOPIC_SV_SCHEDULER_D      string = "uams/server/register/delete"
	TOPIC_SV_SCHEDULER_BOOTUP string = "uams/server/register/bootup"

	TOPIC_SV_HP_BOOTUP string = "uams/server/hp/bootup"
	TOPIC_SV_HP_C      string = "uams/server/hp/create"
	TOPIC_SV_HP_U      string = "uams/server/hp/update"
	TOPIC_SV_HP_D      string = "uams/server/hp/delete"

	TOPIC_SV_USER_U string = "uams/server/user/update"
	TOPIC_SV_USER_D string = "uams/server/user/delete"

	TOPIC_SV_SYSTEM_U string = "uams/server/system/update"
	TOPIC_SV_SYNC     string = "uams/server/sync"
)
