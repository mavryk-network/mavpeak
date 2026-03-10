package mavpay

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mavryk-network/mavpeak/configuration"
	"github.com/mavryk-network/mavpeak/core/common"
	"golang.org/x/exp/maps"
)

type Status struct {
	Services common.AplicationServicesStatus `json:"services,omitempty"`
	Wallet   WalletStatus                    `json:"wallet,omitempty"`
}

func (status *Status) Clone() *Status {
	return &Status{
		Services: common.AplicationServicesStatus{
			Applications: maps.Clone(status.Services.Applications),
			Timestamp:    status.Services.Timestamp,
		},
		Wallet: status.Wallet,
	}
}

type StatusUpdate struct {
	Status *Status
}

func (statusUpdate *StatusUpdate) GetId() string {
	return "mavbake"
}

func (statusUpdate *StatusUpdate) GetData() any {
	return statusUpdate.Status
}

func GetEmptyStatus() *Status {
	return &Status{
		Services: common.AplicationServicesStatus{
			Applications: make(map[string]common.ApplicationServices),
			Timestamp:    time.Now().Unix(),
		},
	}
}

func SetupModule(ctx context.Context, configuration *configuration.MavpayModuleConfiguration, app *fiber.Group, statusChannel chan<- common.StatusUpdate) error {
	err := setupMavpayProvider(configuration, app)
	if err != nil {
		return err
	}

	mavpayStatus := GetEmptyStatus()
	mavpayStatusChannel := make(chan common.StatusUpdate, 100)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case statusUpdate := <-mavpayStatusChannel:
				switch statusUpdate := statusUpdate.(type) {
				case *common.ServicesStatusUpdate:
					application := statusUpdate.Application
					mavpayStatus.Services.Applications[application] = statusUpdate.Status
					mavpayStatus.Services.Timestamp = time.Now().Unix()
				case *WalletBalanceUpdate:
					mavpayStatus.Wallet = statusUpdate.Status
				}

				statusChannel <- &StatusUpdate{
					Status: mavpayStatus.Clone(),
				}
			}
		}
	}()

	common.StartServiceStatusProviders(ctx, configuration.Applications, mavpayStatusChannel)
	startWalletStatusProviders(ctx, configuration.PayoutWallet, configuration.PayoutWalletPreferences, mavpayStatusChannel)

	return nil
}
