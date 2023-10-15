package client

import (
	"github.com/furya-labs/furya/v20/x/txfees/client/cli"
	"github.com/furya-labs/furya/v20/x/txfees/client/rest"

	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

var (
	SubmitUpdateFeeTokenProposalHandler = govclient.NewProposalHandler(cli.NewCmdSubmitUpdateFeeTokenProposal, rest.ProposalUpdateFeeTokenProposalRESTHandler)
)
