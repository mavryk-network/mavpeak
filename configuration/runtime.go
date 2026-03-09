package configuration

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net"
	"os"
	"path/filepath"

	"github.com/hjson/hjson-go/v4"
	"github.com/mavryk-network/mavpeak/constants"
)

type moduleConfigurationbase struct {
	Applications map[string]string `json:"applications,omitempty"`

	Mode PeakMode `json:"mode,omitempty"`
}

type PeakMode string

const (
	PrivatePeakMode PeakMode = "private"
	PublicPeakMode  PeakMode = "public"
	AutoPeakMode    PeakMode = "auto"
)

type versionedConfig interface {
	ToRuntime() *Runtime
}

type deserializedConfigVersion struct {
	Version int `json:"version,omitempty"`
}

type MavrykNode struct {
	Address               string `json:"address"`
	IsRightsProvider      bool   `json:"is_rights_provider,omitempty"`
	IsBlockProvider       bool   `json:"is_block_provider,omitempty"`
	IsGovernanceProvider  bool   `json:"is_governance_provider,omitempty"`
	IsNetworkInfoProvider bool   `json:"is_network_info_provider,omitempty"`
	IsEssential           bool   `json:"is_essential,omitempty"`
	Priority              int    `json:"priority,omitempty"`
}

var (
	MVD_RPC = MavrykNode{
		Address:              "https://rpc.mavryk.network/",
		IsGovernanceProvider: true,
	}
	MVKT_RPC = MavrykNode{
		Address:              "https://mainnet.rpc.mavryk.network/",
		IsBlockProvider:      true,
		IsRightsProvider:     true,
		IsGovernanceProvider: true,
	}
	BAKER_NODE = MavrykNode{
		Address:               "http://127.0.0.1:8732/",
		IsRightsProvider:      true,
		IsBlockProvider:       true,
		IsGovernanceProvider:  true,
		IsNetworkInfoProvider: true,
		IsEssential:           true,
		Priority:              100,
	}
)

type Runtime struct {
	Id     string
	Listen string
	Mode   PeakMode
	// path to the root where are the apps located e.g. /mavbake
	AppRoot string

	Modules map[string]json.RawMessage `json:"modules,omitempty"`

	Nodes map[string]MavrykNode
}

func gerDefaultRuntime() *Runtime {
	return &Runtime{
		Id:      "",
		Listen:  constants.DEFAULT_LISTEN_ADDRESS,
		Mode:    AutoPeakMode,
		Modules: map[string]json.RawMessage{},
	}
}

func (v *Runtime) GetMavbakeModuleConfiguration() (bool, *MavbakeModuleConfiguration) {
	rawConfiguration, ok := v.Modules[constants.MAVBAKE_MODULE_ID]
	if !ok {
		return false, nil
	}

	configuration := getDefaultMavbakeModuleConfiguration()
	err := hjson.Unmarshal(rawConfiguration, configuration)
	if err != nil {
		slog.Error("failed to parse mavbake module configuration", "error", err.Error())
		return false, nil
	}

	for key, value := range configuration.Applications {
		if filepath.IsAbs(value) {
			continue // skip absolute paths
		}
		if value == "" {
			continue // skip empty paths
		}
		configuration.Applications[key] = filepath.Join(v.AppRoot, value)
	}

	if configuration.Mode == "" {
		configuration.Mode = v.Mode
	}
	configuration.Hydrate()

	if err := configuration.Validate(); err != nil {
		slog.Error("failed to validate mavbake module configuration", "error", err.Error())
		return false, nil
	}

	return true, configuration
}

func (v *Runtime) GetMavpayModuleConfiguration() (bool, *MavpayModuleConfiguration) {
	rawConfiguration, ok := v.Modules[constants.MAVPAY_MODULE_ID]
	if !ok {
		return false, nil
	}

	configuration := getDefaultMavpayModuleConfiguration()
	err := hjson.Unmarshal(rawConfiguration, configuration)
	if err != nil {
		slog.Error("failed to parse mavpay module configuration", "error", err.Error())
		return false, nil
	}

	for key, value := range configuration.Applications {
		if filepath.IsAbs(value) {
			continue // skip absolute paths
		}
		configuration.Applications[key] = filepath.Join(v.AppRoot, value)
	}

	if configuration.Mode == "" {
		configuration.Mode = v.Mode
	}
	configuration.Hydrate()

	if err := configuration.Validate(); err != nil {
		slog.Error("failed to validate mavpay module configuration", "error", err.Error())
		return false, nil
	}

	return true, configuration
}

func (r *Runtime) Validate() (*Runtime, error) {
	if r.Listen != "" {
		_, _, err := net.SplitHostPort(r.Listen)
		if err != nil {
			return nil, constants.ErrInvalidListenAddress
		}
	}

	if r.AppRoot == "" {
		return nil, constants.ErrInvalidWorkingDirectory
	}

	// NOTE: should we validate child configurations and exit early if invalid?

	return r, nil
}

func (r *Runtime) Hydrate() *Runtime {
	if r.AppRoot == "" {
		r.AppRoot, _ = os.Getwd()
	}

	if len(r.Nodes) == 0 {
		r.Nodes = map[string]MavrykNode{
			"baker":  BAKER_NODE,
			"MVD":     MVD_RPC,
			"MVKT":   MVKT_RPC,
		}
	}

	return r
}

func Load() (*Runtime, error) {
	var err error
	configFilePath := os.Getenv(constants.ENV_MAVPEAK_CONFIG_FILE)
	if configFilePath == "" {
		configFilePath = "config.hjson"
	}

	configBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		slog.Debug("failed to read config file", "error", err.Error())

		// return config loaded from environment variables
		return gerDefaultRuntime().Hydrate().Validate()
	}

	var configVersion deserializedConfigVersion
	err = hjson.Unmarshal(configBytes, &configVersion)
	if err != nil {
		return nil, errors.Join(constants.ErrInvalidConfigVersion, err)
	}

	var configuration versionedConfig
	switch configVersion.Version {
	case 0:
		configuration, err = load_v0(configBytes)
	default:
		return nil, constants.ErrInvalidConfigVersion
	}

	if err != nil {
		return nil, errors.Join(constants.ErrInvalidConfig, err)
	}

	return configuration.ToRuntime().Hydrate().Validate()
}
