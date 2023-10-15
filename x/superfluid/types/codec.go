package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	authzcodec "github.com/cosmos/cosmos-sdk/x/authz/codec"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgSuperfluidDelegate{}, "furya/superfluid-delegate", nil)
	cdc.RegisterConcrete(&MsgSuperfluidUndelegate{}, "furya/superfluid-undelegate", nil)
	cdc.RegisterConcrete(&MsgLockAndSuperfluidDelegate{}, "furya/lock-and-superfluid-delegate", nil)
	cdc.RegisterConcrete(&MsgSuperfluidUnbondLock{}, "furya/superfluid-unbond-lock", nil)
	cdc.RegisterConcrete(&MsgSuperfluidUndelegateAndUnbondLock{}, "furya/sf-undelegate-and-unbond-lock", nil)
	cdc.RegisterConcrete(&SetSuperfluidAssetsProposal{}, "furya/set-superfluid-assets-proposal", nil)
	cdc.RegisterConcrete(&UpdateUnpoolWhiteListProposal{}, "furya/update-unpool-whitelist", nil)
	cdc.RegisterConcrete(&RemoveSuperfluidAssetsProposal{}, "furya/del-superfluid-assets-proposal", nil)
	cdc.RegisterConcrete(&MsgUnPoolWhitelistedPool{}, "furya/unpool-whitelisted-pool", nil)
	cdc.RegisterConcrete(&MsgUnlockAndMigrateSharesToFullRangeConcentratedPosition{}, "furya/unlock-and-migrate", nil)
	cdc.RegisterConcrete(&MsgCreateFullRangePositionAndSuperfluidDelegate{}, "furya/full-range-and-sf-delegate", nil)
	cdc.RegisterConcrete(&MsgAddToConcentratedLiquiditySuperfluidPosition{}, "furya/add-to-cl-superfluid-position", nil)
	cdc.RegisterConcrete(&MsgUnbondConvertAndStake{}, "furya/unbond-convert-and-stake", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgSuperfluidDelegate{},
		&MsgSuperfluidUndelegate{},
		&MsgLockAndSuperfluidDelegate{},
		&MsgSuperfluidUnbondLock{},
		&MsgSuperfluidUndelegateAndUnbondLock{},
		&MsgUnPoolWhitelistedPool{},
		&MsgUnlockAndMigrateSharesToFullRangeConcentratedPosition{},
		&MsgCreateFullRangePositionAndSuperfluidDelegate{},
		&MsgAddToConcentratedLiquiditySuperfluidPosition{},
		&MsgUnbondConvertAndStake{},
	)

	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&SetSuperfluidAssetsProposal{},
		&RemoveSuperfluidAssetsProposal{},
		&UpdateUnpoolWhiteListProposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterCodec(amino)
	sdk.RegisterLegacyAminoCodec(amino)

	// Register all Amino interfaces and concrete types on the authz Amino codec so that this can later be
	// used to properly serialize MsgGrant and MsgExec instances
	RegisterCodec(authzcodec.Amino)
	amino.Seal()
}
