package configuration

import (
	"github.com/mavryk-network/mavpeak/constants"
	"github.com/mavryk-network/gomavryk/mavryk"
)

type PayoutWalletPreferences struct {
	BalanceWarningThreshold int64 `json:"balance_warning_threshold,omitempty"`
	BalanceErrorThreshold   int64 `json:"balance_error_threshold,omitempty"`
}

type MavpayModuleConfiguration struct {
	moduleConfigurationbase

	PayoutWallet            string                  `json:"payout_wallet,omitempty"`
	PayoutWalletPreferences PayoutWalletPreferences `json:"payout_wallet_preferences,omitempty"`
	ForceDryRun             bool                    `json:"force_dry_run,omitempty"`
}

func getDefaultMavpayModuleConfiguration() *MavpayModuleConfiguration {
	return &MavpayModuleConfiguration{
		moduleConfigurationbase: moduleConfigurationbase{
			Applications: map[string]string{
				"mavpay": constants.DEFAULT_MAVPAY_APP_PATH,
			},
		},
	}
}

func (c *MavpayModuleConfiguration) Hydrate() {

}

func (c *MavpayModuleConfiguration) Validate() error {
	if _, err := mavryk.ParseAddress(c.PayoutWallet); err != nil {
		return constants.ErrInvalidPayoutWallet
	}

	if mavpayAppPath, ok := c.Applications["mavpay"]; !ok || mavpayAppPath == "" {
		return constants.ErrNoMavpayAppPath
	}

	return nil
}
