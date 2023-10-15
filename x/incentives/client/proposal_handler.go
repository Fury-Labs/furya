package client

import (
	"github.com/furya-labs/furya/v20/x/incentives/client/cli"
	"github.com/furya-labs/furya/v20/x/incentives/client/rest"

	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

var (
	HandleCreateGroupsProposal = govclient.NewProposalHandler(cli.NewCmdHandleCreateGroupsProposal, rest.ProposalCreateGroupsRESTHandler)
)
