package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)

const (
	// msgs
	setHotRoutes             = "furya/MsgSetHotRoutes"
	setDeveloperAccount      = "furya/MsgSetDeveloperAccount"
	setMaxPoolPointsPerTx    = "furya/MsgSetMaxPoolPointsPerTx"
	setMaxPoolPointsPerBlock = "furya/MsgSetMaxPoolPointsPerBlock"
	setInfoByPoolType        = "furya/MsgSetInfoByPoolType"
	setBaseDenoms            = "furya/MsgSetBaseDenoms"

	// proposals
	setProtoRevEnabledProposal      = "furya/SetProtoRevEnabledProposal"
	setProtoRevAdminAccountProposal = "furya/SetProtoRevAdminAccountProposal"
)

func init() {
	RegisterCodec(amino)
	sdk.RegisterLegacyAminoCodec(amino)
	amino.Seal()
}

func RegisterCodec(cdc *codec.LegacyAmino) {
	// msgs
	cdc.RegisterConcrete(&MsgSetHotRoutes{}, setHotRoutes, nil)
	cdc.RegisterConcrete(&MsgSetDeveloperAccount{}, setDeveloperAccount, nil)
	cdc.RegisterConcrete(&MsgSetMaxPoolPointsPerTx{}, setMaxPoolPointsPerTx, nil)
	cdc.RegisterConcrete(&MsgSetMaxPoolPointsPerBlock{}, setMaxPoolPointsPerBlock, nil)
	cdc.RegisterConcrete(&MsgSetInfoByPoolType{}, setInfoByPoolType, nil)
	cdc.RegisterConcrete(&MsgSetBaseDenoms{}, setBaseDenoms, nil)

	// proposals
	cdc.RegisterConcrete(&SetProtoRevEnabledProposal{}, setProtoRevEnabledProposal, nil)
	cdc.RegisterConcrete(&SetProtoRevAdminAccountProposal{}, setProtoRevAdminAccountProposal, nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// msgs
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSetHotRoutes{},
		&MsgSetDeveloperAccount{},
		&MsgSetMaxPoolPointsPerTx{},
		&MsgSetMaxPoolPointsPerBlock{},
		&MsgSetInfoByPoolType{},
		&MsgSetBaseDenoms{},
	)

	// proposals
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&SetProtoRevEnabledProposal{},
		&SetProtoRevAdminAccountProposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
