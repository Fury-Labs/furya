package v11

import (
	store "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/fury-labs/furya/v20/app/upgrades"
)

// UpgradeName defines the on-chain upgrade name for the Furya v11 upgrade.
const UpgradeName = "v11"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades:        store.StoreUpgrades{},
}
