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
	cdc.RegisterInterface((*ConcentratedPoolExtension)(nil), nil)

	// msgs
	cdc.RegisterConcrete(&MsgCreatePosition{}, "furya/cl-create-position", nil)
	cdc.RegisterConcrete(&MsgAddToPosition{}, "furya/cl-add-to-position", nil)
	cdc.RegisterConcrete(&MsgWithdrawPosition{}, "furya/cl-withdraw-position", nil)
	cdc.RegisterConcrete(&MsgCollectSpreadRewards{}, "furya/cl-col-sp-rewards", nil)
	cdc.RegisterConcrete(&MsgCollectIncentives{}, "furya/cl-collect-incentives", nil)
	cdc.RegisterConcrete(&MsgFungifyChargedPositions{}, "furya/cl-fungify-charged-positions", nil)

	// gov proposals
	cdc.RegisterConcrete(&CreateConcentratedLiquidityPoolsProposal{}, "furya/create-cl-pools-proposal", nil)
	cdc.RegisterConcrete(&TickSpacingDecreaseProposal{}, "furya/cl-tick-spacing-dec-prop", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterInterface(
		"furya.concentratedliquidity.v1beta1.ConcentratedPoolExtension",
		(*ConcentratedPoolExtension)(nil),
	)

	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreatePosition{},
		&MsgAddToPosition{},
		&MsgWithdrawPosition{},
		&MsgCollectSpreadRewards{},
		&MsgCollectIncentives{},
		&MsgFungifyChargedPositions{},
	)

	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&CreateConcentratedLiquidityPoolsProposal{},
		&TickSpacingDecreaseProposal{},
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
