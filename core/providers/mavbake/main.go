package mavbake

import (
	"context"
	"log/slog"
	"maps"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mavryk-network/mavbake/apps/base"
	"github.com/mavryk-network/mavpeak/configuration"
	"github.com/mavryk-network/mavpeak/core/common"
)

type Status struct {
	Rights   RightsStatus                    `json:"rights,omitempty"`
	Services common.AplicationServicesStatus `json:"services,omitempty"`
	Bakers   BakersStatus                    `json:"bakers,omitempty"`
	Wallets  map[string]base.AmiWalletInfo   `json:"wallets,omitempty"`
}

func (status *Status) Clone() *Status {
	return &Status{
		// no need to clone RightsStatus
		status.Rights,
		common.AplicationServicesStatus{
			Applications: maps.Clone(status.Services.Applications),
			Timestamp:    status.Services.Timestamp,
		},
		status.Bakers,  // no need to clone BakersStatus
		status.Wallets, // no need to clone LedgerStatus
	}
}

func GetEmptyStatus() *Status {
	return &Status{
		Rights: RightsStatus{
			Level:  0,
			Rights: []BlockRights{},
		},
		Services: common.AplicationServicesStatus{
			Applications: make(map[string]common.ApplicationServices),
			Timestamp:    time.Now().Unix(),
		},
		Bakers: BakersStatus{
			Level:  0,
			Bakers: map[string]*BakerStakingStatus{},
		},
		Wallets: make(map[string]base.AmiWalletInfo),
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

func SetupModule(ctx context.Context, configuration *configuration.MavbakeModuleConfiguration, app *fiber.Group, statusChannel chan<- common.StatusUpdate) error {
	err := setupGovernanceProvider(configuration, app)
	if err != nil {
		return err
	}

	mavbakeStatus := GetEmptyStatus()
	mavbakeStatusChannel := make(chan common.StatusUpdate, 100)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case statusUpdate := <-mavbakeStatusChannel:
				switch statusUpdate := statusUpdate.(type) {
				case *common.ServicesStatusUpdate:
					application := statusUpdate.Application
					mavbakeStatus.Services.Applications[application] = statusUpdate.Status
					mavbakeStatus.Services.Timestamp = time.Now().Unix()
				case *RightsStatusUpdate:
					mavbakeStatus.Rights = statusUpdate.RightsStatus
				case *BakersStatusUpdate:
					mavbakeStatus.Bakers = statusUpdate.BakersStatus
				case *WalletsStatusUpdate:
					mavbakeStatus.Wallets = statusUpdate.WalletsStatus
					slog.Info("Ledger status updated", "wallets", len(mavbakeStatus.Wallets), "status", mavbakeStatus.Wallets)
					// case *LedgerStatusUpdate:
					// TODO: LedgerStatusUpdate
				}

				statusChannel <- &StatusUpdate{
					Status: mavbakeStatus.Clone(),
				}
			}
		}
	}()

	if configuration.RightsBlockWindow > 1 {
		startRightsStatusProviders(ctx, configuration.Bakers, configuration.RightsBlockWindow, mavbakeStatusChannel)
	}
	setupBakerStatusProviders(ctx, configuration.Bakers, mavbakeStatusChannel)
	if configuration.ArcBinaryPath != "" {
		go startWalletsStatusProvider(ctx, configuration.Applications["signer"], configuration.ArcBinaryPath, configuration.LedgerWallets, mavbakeStatusChannel)
	}
	common.StartServiceStatusProviders(ctx, configuration.Applications, mavbakeStatusChannel)

	return nil
}
