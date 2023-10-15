package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	authzcodec "github.com/cosmos/cosmos-sdk/x/authz/codec"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// RegisterLegacyAminoCodec registers the necessary x/gamm interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*CFMMPoolI)(nil), nil)
	cdc.RegisterConcrete(&MsgJoinPool{}, "furya/gamm/join-pool", nil)
	cdc.RegisterConcrete(&MsgExitPool{}, "furya/gamm/exit-pool", nil)
	cdc.RegisterConcrete(&MsgSwapExactAmountIn{}, "furya/gamm/swap-exact-amount-in", nil)
	cdc.RegisterConcrete(&MsgSwapExactAmountOut{}, "furya/gamm/swap-exact-amount-out", nil)
	cdc.RegisterConcrete(&MsgJoinSwapExternAmountIn{}, "furya/gamm/join-swap-extern-amount-in", nil)
	cdc.RegisterConcrete(&MsgJoinSwapShareAmountOut{}, "furya/gamm/join-swap-share-amount-out", nil)
	cdc.RegisterConcrete(&MsgExitSwapExternAmountOut{}, "furya/gamm/exit-swap-extern-amount-out", nil)
	cdc.RegisterConcrete(&MsgExitSwapShareAmountIn{}, "furya/gamm/exit-swap-share-amount-in", nil)
	cdc.RegisterConcrete(&UpdateMigrationRecordsProposal{}, "furya/gamm/update-migration-records-proposal", nil)
	cdc.RegisterConcrete(&ReplaceMigrationRecordsProposal{}, "furya/gamm/replace-migration-records-proposal", nil)
	cdc.RegisterConcrete(&CreateConcentratedLiquidityPoolsAndLinktoCFMMProposal{}, "furya/gamm/create-cl-pool-and-cfmm-link", nil)
	cdc.RegisterConcrete(&SetScalingFactorControllerProposal{}, "furya/gamm/scaling-factor-controller", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterInterface(
		"furya.gamm.v1beta1.PoolI", // N.B.: the old proto-path is preserved for backwards-compatibility.
		(*CFMMPoolI)(nil),
	)

	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgJoinPool{},
		&MsgExitPool{},
		&MsgSwapExactAmountIn{},
		&MsgSwapExactAmountOut{},
		&MsgJoinSwapExternAmountIn{},
		&MsgJoinSwapShareAmountOut{},
		&MsgExitSwapExternAmountOut{},
		&MsgExitSwapShareAmountIn{},
	)

	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&UpdateMigrationRecordsProposal{},
		&ReplaceMigrationRecordsProposal{},
		&CreateConcentratedLiquidityPoolsAndLinktoCFMMProposal{},
		&SetScalingFactorControllerProposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/bank module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/staking and
	// defined at the application level.
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	// Register all Amino interfaces and concrete types on the authz Amino codec so that this can later be
	// used to properly serialize MsgGrant and MsgExec instances
	sdk.RegisterLegacyAminoCodec(amino)
	RegisterLegacyAminoCodec(authzcodec.Amino)

	amino.Seal()
}
