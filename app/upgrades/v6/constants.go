package v6

import "github.com/fury-labs/furya/v20/app/upgrades"

const (
	// UpgradeName defines the on-chain upgrade name for the Furya v6 upgrade.
	UpgradeName = "v6"

	// UpgradeHeight defines the block height at which the Furya v6 upgrade is
	// triggered.
	UpgradeHeight = 2_464_000
)

var Fork = upgrades.Fork{
	UpgradeName:    UpgradeName,
	UpgradeHeight:  UpgradeHeight,
	BeginForkLogic: RunForkLogic,
}
