package simulation

import (
	"encoding/json"
	"fmt"

	"github.com/osmosis-labs/osmosis/osmomath"
	"github.com/fury-labs/furya/v20/x/superfluid/types"

	"github.com/cosmos/cosmos-sdk/types/module"
)

// RandomizedGenState generates a random GenesisState for staking.
func RandomizedGenState(simState *module.SimulationState) {
	superfluidGenesis := &types.GenesisState{
		Params: types.Params{
			MinimumRiskFactor: osmomath.NewDecWithPrec(5, 2), // 5%
		},
		SuperfluidAssets:          []types.SuperfluidAsset{},
		FuryEquivalentMultipliers: []types.FuryEquivalentMultiplierRecord{},
	}

	bz, err := json.MarshalIndent(&superfluidGenesis.Params, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated superfluid parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(superfluidGenesis)
}
