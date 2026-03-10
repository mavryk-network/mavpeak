package mavpay

import (
	"context"
	"log/slog"

	"github.com/mavryk-network/mavpeak/configuration"
	"github.com/mavryk-network/mavpeak/core/common"
	"github.com/mavryk-network/gomavryk/rpc"
	"github.com/mavryk-network/gomavryk/mavryk"
)

type WalletStatus struct {
	Address string `json:"address"`
	Balance int64  `json:"balance"`
	Level   string `json:"level"`
}

type WalletBalanceUpdate struct {
	Status WalletStatus
}

func (s *WalletBalanceUpdate) GetId() string {
	return "wallet"
}

func (s *WalletBalanceUpdate) GetData() any {
	return s.Status
}

func startWalletStatusProviders(ctx context.Context, wallet string, preferences configuration.PayoutWalletPreferences, statusChannel chan<- common.StatusUpdate) {
	blockChannelId, blockChannel, err := common.SubscribeToBlockHeaderEvents()
	if err != nil {
		slog.Error("failed to subscribe to block events", "error", err.Error())
		return
	}

	go func() {
		defer func() {
			common.UnsubscribeFromBlockHeaderEvents(blockChannelId)
		}()

		status := WalletStatus{
			Address: wallet,
			Level:   "error",
			Balance: 0,
		}

		for {
			select {
			case <-ctx.Done():
				return
			case _, ok := <-blockChannel:
				if !ok {
					// levelChannel is closed, exit the loop
					return
				}

				if ctx.Done() != nil {
					return
				}

				balance, err := common.AttemptWithRpcClients(ctx, func(client *common.ActiveRpcNode) (int64, error) {
					balance, err := client.GetContractBalance(ctx, mavryk.MustParseAddress(wallet), rpc.Head)
					if err != nil {
						return 0, err
					}
					return balance.Int64(), nil
				})

				if err != nil {
					slog.Error("failed to get contract balance", "error", err.Error())
					continue
				}

				if balance == status.Balance {
					continue
				}

				status.Balance = balance
				switch {
				case status.Balance < (preferences.BalanceErrorThreshold * 1000000) /* mumav */ :
					status.Level = "error"
				case status.Balance < (preferences.BalanceWarningThreshold * 1000000) /* mumav */ :
					status.Level = "warning"
				default:
					status.Level = "ok"
				}

				statusChannel <- &WalletBalanceUpdate{
					Status: status,
				}
			}
		}
	}()

}
