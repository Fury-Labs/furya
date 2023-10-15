package v12

import (
	"github.com/furya-labs/furya/v20/app/upgrades"
	twaptypes "github.com/furya-labs/furya/v20/x/twap/types"

	store "github.com/cosmos/cosmos-sdk/store/types"
)

// UpgradeName defines the on-chain upgrade name for the Furya v12 upgrade.
const UpgradeName = "v12"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added:   []string{twaptypes.StoreKey},
		Deleted: []string{}, // double check bech32ibc
	},
}
