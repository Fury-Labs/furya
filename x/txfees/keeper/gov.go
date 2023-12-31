package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/fury-labs/furya/v20/x/txfees/types"
)

func (k Keeper) HandleUpdateFeeTokenProposal(ctx sdk.Context, p *types.UpdateFeeTokenProposal) error {
	// setFeeToken internally calls ValidateFeeToken
	for _, feeToken := range p.Feetokens {
		if err := k.setFeeToken(ctx, feeToken); err != nil {
			return err
		}
	}
	return nil
}
