package client

import (
	"github.com/fury-labs/furya/v20/x/concentrated-liquidity/client/cli"
	"github.com/fury-labs/furya/v20/x/concentrated-liquidity/client/rest"

	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

var (
	TickSpacingDecreaseProposalHandler             = govclient.NewProposalHandler(cli.NewTickSpacingDecreaseProposal, rest.ProposalTickSpacingDecreaseRESTHandler)
	CreateConcentratedLiquidityPoolProposalHandler = govclient.NewProposalHandler(cli.NewCmdCreateConcentratedLiquidityPoolsProposal, rest.ProposalCreateConcentratedLiquidityPoolHandler)
)
