package stableswap

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/furya-labs/furya/osmomath"
	"github.com/furya-labs/furya/v20/x/gamm/types"
)

func createTestPool(t *testing.T, poolLiquidity sdk.Coins, spreadFactor, exitFee osmomath.Dec, scalingFactors []uint64) types.CFMMPoolI {
	t.Helper()
	scalingFactors, _ = applyScalingFactorMultiplier(scalingFactors)

	pool, err := NewStableswapPool(1, PoolParams{
		SwapFee: spreadFactor,
		ExitFee: exitFee,
	}, poolLiquidity, scalingFactors, "", "")

	require.NoError(t, err)

	return &pool
}
