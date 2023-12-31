package v8

import (
	"github.com/fury-labs/furya/v20/app/upgrades"
	v8constants "github.com/fury-labs/furya/v20/app/upgrades/v8/constants"
)

const (
	// UpgradeName defines the on-chain upgrade name for the Furya v8 upgrade.
	UpgradeName = v8constants.UpgradeName

	// UpgradeHeight defines the block height at which the Furya v8 upgrade is
	// triggered.
	UpgradeHeight = v8constants.UpgradeHeight
)

var Fork = upgrades.Fork{
	UpgradeName:    UpgradeName,
	UpgradeHeight:  UpgradeHeight,
	BeginForkLogic: RunForkLogic,
}
