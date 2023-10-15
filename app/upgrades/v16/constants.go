package v16

import (
	"github.com/fury-labs/furya/v20/app/upgrades"

	store "github.com/cosmos/cosmos-sdk/store/types"

	cltypes "github.com/fury-labs/furya/v20/x/concentrated-liquidity/types"
	cosmwasmpooltypes "github.com/fury-labs/furya/v20/x/cosmwasmpool/types"
)

// UpgradeName defines the on-chain upgrade name for the Furya v16 upgrade.
const UpgradeName = "v16"

// new token factory parameters
//
// at the current gas price of 0.0025ufury, this corresponds to 0.1 FURY per
// denom creation.
const NewDenomCreationGasConsume uint64 = 40_000_000

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added:   []string{cltypes.StoreKey, cosmwasmpooltypes.StoreKey},
		Deleted: []string{},
	},
}
