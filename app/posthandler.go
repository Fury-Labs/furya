package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	protorevkeeper "github.com/fury-labs/furya/v20/x/protorev/keeper"
)

func NewPostHandler(protoRevKeeper *protorevkeeper.Keeper) sdk.AnteHandler {
	protoRevDecorator := protorevkeeper.NewProtoRevDecorator(*protoRevKeeper)
	return sdk.ChainAnteDecorators(protoRevDecorator)
}
