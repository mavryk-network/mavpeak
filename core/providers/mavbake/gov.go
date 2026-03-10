package mavbake

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/mavryk-network/mavpeak/configuration"
	"github.com/mavryk-network/mavpeak/constants"
	"github.com/mavryk-network/mavpeak/core/common"
	"github.com/mavryk-network/mavpeak/util"
	"github.com/mavryk-network/gomavryk/codec"
	"github.com/mavryk-network/gomavryk/rpc"
	"github.com/mavryk-network/gomavryk/signer/remote"
	"github.com/mavryk-network/gomavryk/mavryk"
)

type GovernanceStatus struct {
}

type GovernanceProvider struct {
	configuration *configuration.MavbakeModuleConfiguration

	signerUrl string
}

type VoteList map[string][]string

type GovernancePeriodDetail struct {
	Info      *rpc.VotingPeriodInfo `json:"info"`
	Voters    rpc.VoterList         `json:"voters"`
	Summary   *rpc.BallotSummary    `json:"summary"`
	Proposal  mavryk.ProtocolHash    `json:"proposal"`
	Proposals rpc.ProposalList      `json:"proposals"`
	Quorum    int                   `json:"quorum"`
	Ballots   rpc.BallotList        `json:"ballots"`
	Votes     VoteList              `json:"votes"`
}

type UpvoteParams struct {
	Source    mavryk.Address        `json:"source"`
	Proposals []mavryk.ProtocolHash `json:"proposals"`
	Period    int32                `json:"period"`
}

func (p *UpvoteParams) ToContents() codec.Operation {
	return &codec.Proposals{
		Source:    p.Source,
		Proposals: p.Proposals,
		Period:    p.Period,
	}
}

type VoteParams struct {
	Source   mavryk.Address      `json:"source"`
	Proposal mavryk.ProtocolHash `json:"proposal"`
	Period   int32              `json:"period"`
	Ballot   string             `json:"ballot"`
}

func (p *VoteParams) ToContents() codec.Operation {
	return &codec.Ballot{
		Source:   p.Source,
		Proposal: p.Proposal,
		Period:   p.Period,
		Ballot:   mavryk.ParseBallotVote(p.Ballot),
	}
}

func (governanceProvider *GovernanceProvider) CanVote() bool {
	return governanceProvider.configuration.Mode == configuration.PrivatePeakMode
}

func attemptWithGovernanceRpcClients[T any](ctx context.Context, f func(client *common.ActiveRpcNode) (T, error)) (T, error) {
	return common.AttemptWithRpcClients(ctx, func(client *common.ActiveRpcNode) (T, error) {
		var result T
		if !client.IsGovernanceProvider {
			return result, errors.New("not a governance provider")
		}
		return f(client)
	})
}

func wrapInWaithGroup(wg *sync.WaitGroup, f func()) {
	wg.Add(1)
	go func() {
		f()
		wg.Done()
	}()
}

func (governanceProvider *GovernanceProvider) startVotersCollector(ctx context.Context, detail *GovernancePeriodDetail, wg *sync.WaitGroup) {
	wrapInWaithGroup(wg, func() {
		voters, _ := attemptWithGovernanceRpcClients(ctx, func(client *common.ActiveRpcNode) (rpc.VoterList, error) {
			voters, err := client.ListVoters(ctx, rpc.Head)
			if err != nil {
				return nil, err
			}
			// voters[0].Power <- power
			return voters, err
		})

		detail.Voters = voters
	})
}

func (governanceProvider *GovernanceProvider) startProposalsCollector(ctx context.Context, detail *GovernancePeriodDetail, wg *sync.WaitGroup) {
	wrapInWaithGroup(wg, func() {
		proposals, _ := attemptWithGovernanceRpcClients(ctx, func(client *common.ActiveRpcNode) (rpc.ProposalList, error) {
			proposals, err := client.ListProposals(ctx, rpc.Head)
			if err != nil {
				return nil, err
			}
			return proposals, err
		})
		detail.Proposals = proposals
	})
}

// curl 127.0.0.1:8732/chains/main/blocks/head/context/raw/json/votes/proposals?depth=1
// [["Pt1JoinAscentToMountVinsonAGNqxgMLDAB8TqZpDwMTU5eCx",["tz1P6WKJu2rcbxKiKRZHKQKmKrpC9TfW1AwM","tz1LjZjdF1wFgUtVyNsrr8P1uYaoBJGTTPyr"]]]

func (governanceProvider *GovernanceProvider) startVotesCollector(ctx context.Context, detail *GovernancePeriodDetail, wg *sync.WaitGroup) {
	wrapInWaithGroup(wg, func() {
		votes, _ := attemptWithGovernanceRpcClients(ctx, func(client *common.ActiveRpcNode) (VoteList, error) {
			var rawVotes [][]any

			err := client.Get(ctx, "chains/main/blocks/head/context/raw/json/votes/proposals?depth=1", &rawVotes)
			if err != nil {
				return nil, err
			}

			result := make(VoteList)
			for _, rawVote := range rawVotes {
				if len(rawVote) != 2 {
					continue
				}
				if proposal, ok := rawVote[0].(string); ok {
					voters := rawVote[1].([]any)
					var votersList []string
					if len(voters) > 0 {
						votersList = lo.Map(voters, func(voter any, _ int) string {
							if v, ok := voter.(string); ok {
								return v
							}
							return ""
						})
					}
					result[proposal] = votersList
				}
			}

			return result, err
		})
		detail.Votes = votes
	})
}

func (governanceProvider *GovernanceProvider) startQuorumCollector(ctx context.Context, detail *GovernancePeriodDetail, wg *sync.WaitGroup) {
	wrapInWaithGroup(wg, func() {
		quorum, _ := attemptWithGovernanceRpcClients(ctx, func(client *common.ActiveRpcNode) (int, error) {
			quorum, err := client.GetVoteQuorum(ctx, rpc.Head)
			if err != nil {
				return 0, err
			}
			return quorum, err
		})
		detail.Quorum = quorum
	})
}

func (governanceProvider *GovernanceProvider) startBallotsCollector(ctx context.Context, detail *GovernancePeriodDetail, wg *sync.WaitGroup) {
	wrapInWaithGroup(wg, func() {
		ballotList, _ := attemptWithGovernanceRpcClients(ctx, func(client *common.ActiveRpcNode) (rpc.BallotList, error) {
			ballots, err := client.ListBallots(ctx, rpc.Head)
			if err != nil {
				return nil, err
			}
			return ballots, err
		})
		detail.Ballots = ballotList
	})
}

func (governanceProvider *GovernanceProvider) startSummaryCollector(ctx context.Context, detail *GovernancePeriodDetail, wg *sync.WaitGroup) {
	wrapInWaithGroup(wg, func() {
		summary, _ := attemptWithGovernanceRpcClients(ctx, func(client *common.ActiveRpcNode) (*rpc.BallotSummary, error) {
			summary, err := client.GetVoteResult(ctx, rpc.Head)
			if err != nil {
				return nil, err
			}
			return &summary, err
		})
		detail.Summary = summary
	})
}

func (governanceProvider *GovernanceProvider) startProtocolCollector(ctx context.Context, detail *GovernancePeriodDetail, wg *sync.WaitGroup) {
	wrapInWaithGroup(wg, func() {
		currentProposal, _ := attemptWithGovernanceRpcClients(ctx, func(client *common.ActiveRpcNode) (mavryk.ProtocolHash, error) {
			proposal, err := client.GetVoteProposal(ctx, rpc.Head)
			if err != nil {
				return mavryk.ProtocolHash{}, err
			}
			return proposal, err
		})
		detail.Proposal = currentProposal
	})
}

func (governanceProvider *GovernanceProvider) GetAvailablePkhs(ctx context.Context) ([]string, error) {
	return governanceProvider.configuration.Bakers, nil
}

func (governanceProvider *GovernanceProvider) GetGovernancePeriodDetail(ctx context.Context) (*GovernancePeriodDetail, error) {
	periodInfo, err := attemptWithGovernanceRpcClients(ctx, func(client *common.ActiveRpcNode) (*rpc.VotingPeriodInfo, error) {
		meta, err := client.GetBlockMetadata(ctx, rpc.Head)
		if err != nil {
			return nil, err
		}
		return meta.VotingPeriodInfo, err
	})
	if err != nil {
		slog.Warn("failed to get voting period kind", "error", err.Error())
		return nil, err
	}

	detail := &GovernancePeriodDetail{
		Info: periodInfo,
	}

	var wg sync.WaitGroup

	governanceProvider.startVotersCollector(ctx, detail, &wg)

	if periodInfo.VotingPeriod.Kind == mavryk.VotingPeriodProposal {
		governanceProvider.startProposalsCollector(ctx, detail, &wg)
		governanceProvider.startVotesCollector(ctx, detail, &wg)
	}

	if periodInfo.VotingPeriod.Kind == mavryk.VotingPeriodExploration || periodInfo.VotingPeriod.Kind == mavryk.VotingPeriodPromotion {
		governanceProvider.startQuorumCollector(ctx, detail, &wg)
		governanceProvider.startBallotsCollector(ctx, detail, &wg)
		governanceProvider.startSummaryCollector(ctx, detail, &wg)
		governanceProvider.startProtocolCollector(ctx, detail, &wg)
	}

	wg.Wait()
	return detail, nil
}

func (governanceProvider *GovernanceProvider) buildAndBroadcastGovernanceOperation(ctx context.Context, pkh mavryk.Address, contents codec.Operation) (mavryk.OpHash, error) {
	rs, err := remote.New(governanceProvider.signerUrl, nil)
	if err != nil {
		err = util.TryUnwrapRPCError(err)
		slog.Error("failed to create remote signer", "error", err.Error())
		return mavryk.OpHash{}, errors.Join(constants.ErrFailedToCreateRemoteSigner, err)
	}

	key, err := rs.GetKey(ctx, pkh)
	if err != nil {
		err = util.TryUnwrapRPCError(err)
		slog.Error("failed to get key", "error", err.Error())
		return mavryk.OpHash{}, errors.Join(constants.ErrFailedToGetPublicKey, err)
	}

	// complete the operation
	op, err := attemptWithGovernanceRpcClients(ctx, func(client *common.ActiveRpcNode) (*codec.Op, error) {
		params, err := client.GetParams(ctx, rpc.Head)
		if err != nil {
			return nil, err
		}
		op := codec.NewOp().WithContents(contents).WithSource(pkh)
		op.WithTTL(constants.MAX_OPERATION_TTL)
		op.WithContents(contents)

		op = op.WithParams(params)
		err = client.Complete(ctx, op, key)
		if err != nil {
			return op, err
		}

		return op, err
	})
	if err != nil {
		err = util.TryUnwrapRPCError(err)
		slog.Error("failed to complete operation", "error", err.Error())
		return mavryk.OpHash{}, errors.Join(constants.ErrFailedToCompleteOperation, err)
	}

	signature, err := rs.SignOperation(ctx, pkh, op)
	if err != nil {
		err = util.TryUnwrapRPCError(err)
		slog.Error("failed to sign operation", "error", err.Error())
		return mavryk.OpHash{}, errors.Join(constants.ErrFailedToSignOperation, err)
	}
	op = op.WithSignature(signature)

	opHash, err := attemptWithGovernanceRpcClients(ctx, func(client *common.ActiveRpcNode) (mavryk.OpHash, error) {
		opHash, err := client.Broadcast(ctx, op)
		if err != nil {
			slog.Error("failed to broadcast operation", "error", err.Error())
			return mavryk.OpHash{}, err
		}
		return opHash, err
	})
	if err != nil {
		err = util.TryUnwrapRPCError(err)
		slog.Error("failed to broadcast operation", "error", err.Error())
		return mavryk.OpHash{}, errors.Join(constants.ErrFailedToBroadcastOperation, err)
	}

	return opHash, nil
}

func (governanceProvider *GovernanceProvider) Upvote(ctx context.Context, params *UpvoteParams) (mavryk.OpHash, error) {
	return governanceProvider.buildAndBroadcastGovernanceOperation(ctx, params.Source, params.ToContents())
}

func (governanceProvider *GovernanceProvider) Vote(ctx context.Context, params *VoteParams) (mavryk.OpHash, error) {
	return governanceProvider.buildAndBroadcastGovernanceOperation(ctx, params.Source, params.ToContents())
}

func (governanceProvider *GovernanceProvider) WaitConfirmation(ctx context.Context, opHash string) (bool, error) {
	op, err := mavryk.ParseOpHash(opHash)
	if err != nil {
		return false, err
	}

	ctx, cancel := context.WithTimeout(ctx, constants.MAX_WAIT_FOR_CONFIRMATION*time.Second)
	defer cancel()
	result, err := attemptWithGovernanceRpcClients(ctx, func(client *common.ActiveRpcNode) (bool, error) {
		result := rpc.NewResult(op)
		client.Listen()
		defer client.Close()
		result.Listen(client.BlockObserver)
		result.WaitContext(ctx)

		if err := result.Err(); err != nil {
			return false, err
		}

		return true, err
	})
	return result, err
}

func (governanceProvider *GovernanceProvider) RegisterApi(app *fiber.Group) error {
	app.Get("/governance/can-vote", func(c *fiber.Ctx) error {
		return c.JSON(governanceProvider.CanVote())
	})

	app.Get("/governance/period-detail", func(c *fiber.Ctx) error {
		if !governanceProvider.CanVote() {
			return c.Status(403).SendString("not allowed")
		}

		periodInfo, err := governanceProvider.GetGovernancePeriodDetail(c.Context())
		if err != nil {
			return c.Status(500).SendString("failed to get governance period detail")
		}

		return c.JSON(periodInfo)
	})

	app.Get("/governance/available-pkhs", func(c *fiber.Ctx) error {
		if !governanceProvider.CanVote() {
			return c.Status(403).SendString("not allowed")
		}

		pkhs, err := governanceProvider.GetAvailablePkhs(c.Context())
		if err != nil {
			return c.Status(500).SendString("failed to get available pkhs")
		}
		return c.JSON(pkhs)
	})

	app.Post("/governance/vote", func(c *fiber.Ctx) error {
		if !governanceProvider.CanVote() {
			return c.Status(403).SendString("not allowed")
		}

		var params VoteParams
		if err := c.BodyParser(&params); err != nil {
			return c.Status(400).SendString("invalid request")
		}

		opHash, err := governanceProvider.Vote(c.Context(), &params)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.JSON(opHash)
	})

	app.Post("/governance/upvote", func(c *fiber.Ctx) error {
		if !governanceProvider.CanVote() {
			return c.Status(403).SendString("not allowed")
		}

		var params UpvoteParams
		if err := c.BodyParser(&params); err != nil {
			slog.Error("failed to parse upvote params", "error", err.Error())
			return c.Status(400).SendString("invalid request")
		}

		opHash, err := governanceProvider.Upvote(c.Context(), &params)
		if err != nil {
			slog.Error("failed to upvote", "error", err.Error())
			return c.Status(500).SendString(err.Error())
		}

		return c.JSON(opHash)
	})

	app.Post("/governance/wait-for-apply", func(c *fiber.Ctx) error {
		if !governanceProvider.CanVote() {
			return c.Status(403).SendString("not allowed")
		}

		var params string
		if err := c.BodyParser(&params); err != nil {
			slog.Error("failed to parse upvote params", "error", err.Error())
			return c.Status(400).SendString("invalid request")
		}

		applied, err := governanceProvider.WaitConfirmation(c.Context(), params)
		if err != nil {
			slog.Error("failed to upvote", "error", err.Error())
			return c.Status(500).SendString("failed to vote")
		}

		return c.JSON(applied)
	})

	return nil
}

func setupGovernanceProvider(configuration *configuration.MavbakeModuleConfiguration, app *fiber.Group) error {
	provider := &GovernanceProvider{
		configuration: configuration,
		signerUrl:     configuration.SignerUrl,
	}

	return provider.RegisterApi(app)
}
