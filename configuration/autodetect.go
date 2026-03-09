package configuration

import (
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"slices"

	"github.com/hjson/hjson-go/v4"
	"github.com/mavryk-network/mavbake/ami"
	"github.com/mavryk-network/mavpay/configuration"
	"github.com/mavryk-network/mavpeak/constants"

	signer_engines "github.com/mavryk-network/mavpay/engines/signer"
	"github.com/mavryk-network/mavpay/state"
)

/*
	mavpay: {
		applications: {
				mavpay: mavpay
		}
		payout_wallet: tz1X7U9XxVz6NDxL4DSZhijME61PW45bYUJE
		payout_wallet_preferences: {
			balance_warning_threshold: 100
			balance_error_threshold: 50
		}
	}
*/
func autoDetectMavpayConfiguration(rootDir string) (json.RawMessage, error) {
	if rootDir == "" {
		return nil, errors.New("rootDir is empty")
	}

	workingDir := path.Join(rootDir, constants.DEFAULT_MAVPAY_APP_PATH)
	if !ami.IsAppInstalled(workingDir) {
		return nil, errors.New("mavpay not found, skipping")
	}

	state.Init(workingDir, state.StateInitOptions{})

	config, err := configuration.Load()
	if err != nil {
		return nil, errors.Join(errors.New("failed to load mavpay configuration"), err)
	}

	signerEngine := state.Global.SignerOverride
	if signerEngine == nil {
		signerEngine, err = signer_engines.Load(string(config.PayoutConfiguration.WalletMode))
		if err != nil {
			return nil, errors.Join(errors.New("failed to to load mavpay signer"), err)
		}
	}

	slog.Info("Checking PKH used for mavpay payouts...")
	pkh := signerEngine.GetPKH()

	mavpayModuleConfiguration := &MavpayModuleConfiguration{
		moduleConfigurationbase: moduleConfigurationbase{
			Applications: map[string]string{
				"mavpay": constants.DEFAULT_MAVPAY_APP_PATH,
			},
		},
		PayoutWallet: pkh.String(),
		PayoutWalletPreferences: PayoutWalletPreferences{
			BalanceWarningThreshold: 100,
			BalanceErrorThreshold:   50,
		},
	}
	return hjson.MarshalWithOptions(mavpayModuleConfiguration, hjson.DefaultOptions())
}

/*
	{
		"configuration": {
				...
				"additional_key_aliases": [ "key" ]
		},
	}
*/
type nodeAppJsonPartialConfiguration struct {
	RemoteSignerUrl      string   `json:"REMOTE_SIGNER_ADDR,omitempty"`
	AdditionalKeyAliases []string `json:"additional_key_aliases,omitempty"`
}

type nodeAppJsonPartial struct {
	Configuration nodeAppJsonPartialConfiguration
}

/*
	mavbake: {
		bakers: [
			tz1P6WKJu2rcbxKiKRZHKQKmKrpC9TfW1AwM
			tz1hZvgjekGo7DmQjWh7XnY5eLQD8wNYPczE
		]
	}
*/
func autoDetectMavbakeConfiguration(rootDir string) (json.RawMessage, error) {
	nodeAppPath := path.Join(rootDir, constants.DEFAULT_NODE_APP_PATH)
	signerAppPath := path.Join(rootDir, constants.DEFAULT_SIGNER_APP_PATH)

	applications := map[string]string{}
	bakers := []string{}
	remoteSignerUrl := constants.DEFAULT_BAKER_SIGNER_URL
	if ami.IsAppInstalled(nodeAppPath) {
		applications["node"] = constants.DEFAULT_NODE_APP_PATH

		keys := []string{"baker"}

		nodeAppConfig := nodeAppJsonPartial{
			Configuration: nodeAppJsonPartialConfiguration{
				AdditionalKeyAliases: []string{},
			},
		}

		nodeAppPath := path.Join(rootDir, constants.DEFAULT_NODE_APP_PATH, "app.json")
		fileContent := []byte("{}")
		if data, err := os.ReadFile(nodeAppPath); err == nil {
			fileContent = data
		} else {
			nodeAppPath = path.Join(rootDir, constants.DEFAULT_NODE_APP_PATH, "app.hjson")
			if data, err := os.ReadFile(nodeAppPath); err == nil {
				fileContent = data
			}
		}
		if err := hjson.Unmarshal(fileContent, &nodeAppConfig); err == nil {
			keys = append(keys, nodeAppConfig.Configuration.AdditionalKeyAliases...)
		}

		if nodeAppConfig.Configuration.RemoteSignerUrl != "" {
			remoteSignerUrl = nodeAppConfig.Configuration.RemoteSignerUrl
		}

		// read pkhs based on aliases
		pathToPkhs := path.Join(rootDir, constants.DEFAULT_NODE_APP_PATH, "data/.mavryk-client/public_key_hashs")
		pkhs := nodePublicKeys{}
		fileContent = []byte{}
		if data, err := os.ReadFile(pathToPkhs); err == nil {
			fileContent = data
		}
		if err := hjson.Unmarshal(fileContent, &pkhs); err != nil {
			return nil, errors.New("failed to read public key hashes")
		}
		for _, pkh := range pkhs {
			if slices.Contains(keys, pkh.Name) {
				bakers = append(bakers, pkh.Hash)
			}
		}
	} else {
		slog.Warn("Node app is not found, skipping")
	}

	if ami.IsAppInstalled(signerAppPath) {
		applications["signer"] = constants.DEFAULT_SIGNER_APP_PATH
	} else {
		slog.Warn("Signer app is not found, skipping")
	}

	mavbakeModuleConfiguration := &MavbakeModuleConfiguration{
		moduleConfigurationbase: moduleConfigurationbase{
			Applications: applications,
		},
		SignerUrl:         remoteSignerUrl,
		RightsBlockWindow: constants.DEFAULT_RIGHTS_BLOCK_WINDOW,
		Bakers:            bakers,
	}

	return hjson.MarshalWithOptions(mavbakeModuleConfiguration, hjson.DefaultOptions())
}

func AutoDetect(rootDir string, destinationFile string) {
	modules := map[string]json.RawMessage{}
	absRootDir, err := filepath.Abs(rootDir)
	if err != nil {
		slog.Warn("Failed to get absolute path to root dir", "error", err.Error())
	} else {
		rootDir = absRootDir
	}

	mavpayConfig, err := autoDetectMavpayConfiguration(rootDir)
	if err != nil {
		slog.Warn("Failed to auto-detect mavpay configuration", "error", err.Error())
	} else {
		modules[constants.MAVPAY_MODULE_ID] = mavpayConfig
	}

	mavbakeConfig, err := autoDetectMavbakeConfiguration(rootDir)
	if err != nil {
		slog.Warn("Failed to auto-detect mavbake configuration", "error", err.Error())
	} else {
		modules[constants.MAVBAKE_MODULE_ID] = mavbakeConfig
	}

	config := Runtime{
		Id:      "",
		AppRoot: rootDir,
		Listen:  constants.DEFAULT_LISTEN_ADDRESS,
		Modules: modules,
		Mode:    AutoPeakMode,
	}

	if data, err := hjson.MarshalWithOptions(config, hjson.DefaultOptions()); err == nil {
		if err := os.WriteFile(destinationFile, data, 0644); err != nil {
			slog.Error("Failed to write autodetected configuration", "error", err.Error())
		}
	} else {
		slog.Error("Failed to marshal autodetected configuration", "error", err.Error())
	}
}
