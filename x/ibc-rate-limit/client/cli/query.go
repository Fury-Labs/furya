package cli

import (
	"github.com/spf13/cobra"

	"github.com/furya-labs/furya/osmoutils/osmocli"
	"github.com/fury-labs/furya/v20/x/ibc-rate-limit/client/queryproto"
	"github.com/fury-labs/furya/v20/x/ibc-rate-limit/types"
)

// GetQueryCmd returns the cli query commands for this module.
func GetQueryCmd() *cobra.Command {
	cmd := osmocli.QueryIndexCmd(types.ModuleName)

	cmd.AddCommand(
		osmocli.GetParams[*queryproto.ParamsRequest](
			types.ModuleName, queryproto.NewQueryClient),
	)

	return cmd
}
