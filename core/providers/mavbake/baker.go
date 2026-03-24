package mavbake

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mavryk-network/mavpeak/core/common"
	"github.com/mavryk-network/gomavryk/rpc"
	"github.com/mavryk-network/gomavryk/mavryk"
)

type BakersStatus struct {
	Level  int64                          `json:"level"`
	Bakers map[string]*BakerStakingStatus `json:"bakers"`
}

type BakersStatusUpdate struct {
	BakersStatus
}

func (s *BakersStatusUpdate) GetId() string {
	return "bakers"
}

func (s *BakersStatusUpdate) GetData() any {
	return s
}

type BakerStakingStatus struct {
	FullBalance    string `json:"full_balance"`
	StakedBalance  string `json:"staked_balance"`
	LiquidBalance  string `json:"liquid_balance"`
}

func getBakerBalances(ctx context.Context, client *common.ActiveRpcNode, addr mavryk.Address, id rpc.BlockID) (*BakerStakingStatus, error) {
	base := fmt.Sprintf("chains/main/blocks/%s/context/contracts/%s", id, addr)

	var fullBalance string
	if err := client.Get(ctx, base+"/full_balance", &fullBalance); err != nil {
		return nil, fmt.Errorf("full_balance: %w", err)
	}

	var stakedBalance string
	if err := client.Get(ctx, base+"/staked_balance", &stakedBalance); err != nil {
		return nil, fmt.Errorf("staked_balance: %w", err)
	}

	var liquidBalance string
	if err := client.Get(ctx, base+"/balance", &liquidBalance); err != nil {
		return nil, fmt.Errorf("balance: %w", err)
	}

	return &BakerStakingStatus{
		FullBalance:   fullBalance,
		StakedBalance: stakedBalance,
		LiquidBalance: liquidBalance,
	}, nil
}

func getBakerStatusFor(ctx context.Context, baker string) (*BakerStakingStatus, error) {
	addr, err := mavryk.ParseAddress(baker)
	if err != nil {
		return nil, err
	}
	return common.AttemptWithRpcClients(ctx, func(client *common.ActiveRpcNode) (*BakerStakingStatus, error) {
		return getBakerBalances(ctx, client, addr, rpc.Head)
	})
}

func setupBakerStatusProviders(ctx context.Context, bakers []string, statusChannel chan<- common.StatusUpdate) {
	blockChannelId, blockChannel, err := common.SubscribeToBlockHeaderEvents()
	if err != nil {
		slog.Error("failed to subscribe to block events", "error", err.Error())
		return
	}

	go func() {
		defer func() {
			common.UnsubscribeFromBlockHeaderEvents(blockChannelId)
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case block, ok := <-blockChannel:
				if !ok {
					return
				}

				if ctx.Done() != nil {
					return
				}

				bakersStatus := map[string]*BakerStakingStatus{}
				for _, baker := range bakers {
					status, err := getBakerStatusFor(ctx, baker)
					if err != nil {
						slog.Debug("failed to get baker status", "baker", baker, "error", err.Error())
						continue
					}
					bakersStatus[baker] = status
				}
				statusChannel <- &BakersStatusUpdate{
					BakersStatus: BakersStatus{
						Level:  block.Level,
						Bakers: bakersStatus,
					},
				}
			}
		}
	}()
}
