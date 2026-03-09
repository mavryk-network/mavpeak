package constants

const (
	MAVPEAK_VERSION  = "<VERSION>"
	MAVPEAK_CODENAME = "<CODENAME>"

	DEFAULT_LISTEN_ADDRESS       = "localhost:8733"
	DEFAULT_HTTP_TIMEOUT_SECONDS = 30

	// mavbake
	MAVBAKE_MODULE_ID             = "mavbake"
	ENV_MAVPEAK_CONFIG_FILE       = "MAVPEAK_CONFIG_FILE"
	MAX_SERVICES_REFRESH_INTERVAL = 300 // 5 minutes
	MIN_SERVICES_REFRESH_INTERVAL = 5   // 5 seconds
	DEFAULT_NODE_APP_PATH         = "node"
	DEFAULT_SIGNER_APP_PATH       = "signer"
	DEFAULT_BAKER_NODE_URL        = "http://localhost:8732"
	DEFAULT_BAKER_SIGNER_URL      = "http://localhost:20090"
	DEFAULT_RIGHTS_BLOCK_WINDOW   = 50
	DEFAULT_MONITOR_LEDGER_STATUS = true
	DEFAULT_ARC_BINARY_PATH       = ""

	// mavpay
	MAVPAY_MODULE_ID        = "mavpay"
	DEFAULT_MAVPAY_APP_PATH = "pay"

	// tx constants
	MAX_OPERATION_TTL         = 12
	MAX_WAIT_FOR_CONFIRMATION = 120
)

var (
	PRIVATE_NETWORK_HOSTS = []string{
		"localhost",
		"127.0.0.1",
		"::1",
	}
)
