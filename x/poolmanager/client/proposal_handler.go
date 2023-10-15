package client

import (
	"github.com/furya-labs/furya/v20/x/poolmanager/client/cli"
	"github.com/furya-labs/furya/v20/x/poolmanager/client/rest"

	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

var (
	DenomPairTakerFeeProposalHandler = govclient.NewProposalHandler(cli.NewCmdHandleDenomPairTakerFeeProposal, rest.ProposalDenomPairTakerFeeRESTHandler)
)
