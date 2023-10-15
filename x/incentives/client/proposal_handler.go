package client

import (
	"github.com/fury-labs/furya/v20/x/incentives/client/cli"
	"github.com/fury-labs/furya/v20/x/incentives/client/rest"

	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

var (
	HandleCreateGroupsProposal = govclient.NewProposalHandler(cli.NewCmdHandleCreateGroupsProposal, rest.ProposalCreateGroupsRESTHandler)
)
