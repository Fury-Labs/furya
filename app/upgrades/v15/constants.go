package v15

import (
	store "github.com/cosmos/cosmos-sdk/store/types"
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v4/router/types"
	icqtypes "github.com/cosmos/ibc-apps/modules/async-icq/v4/types"

	"github.com/fury-labs/furya/v20/app/upgrades"
	poolmanagertypes "github.com/fury-labs/furya/v20/x/poolmanager/types"
	protorevtypes "github.com/fury-labs/furya/v20/x/protorev/types"
	valsetpreftypes "github.com/fury-labs/furya/v20/x/valset-pref/types"
)

// UpgradeName defines the on-chain upgrade name for the Furya v15 upgrade.
const UpgradeName = "v15"

// pool ids to migrate
const (
	stFURY_FURYPoolId   = 833
	stJUNO_JUNOPoolId   = 817
	stSTARS_STARSPoolId = 810
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added:   []string{poolmanagertypes.StoreKey, valsetpreftypes.StoreKey, protorevtypes.StoreKey, icqtypes.StoreKey, packetforwardtypes.StoreKey},
		Deleted: []string{},
	},
}
