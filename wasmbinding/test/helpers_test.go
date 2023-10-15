package wasmbinding

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/fury-labs/furya/osmomath"
	"github.com/fury-labs/furya/v20/app"
)

func CreateTestInput() (*app.FuryaApp, sdk.Context) {
	furya := app.Setup(false)
	ctx := furya.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "furya-1", Time: time.Now().UTC()})
	return furya, ctx
}

func FundAccount(t *testing.T, ctx sdk.Context, furya *app.FuryaApp, acct sdk.AccAddress) {
	t.Helper()
	err := simapp.FundAccount(furya.BankKeeper, ctx, acct, sdk.NewCoins(
		sdk.NewCoin("ufury", osmomath.NewInt(10000000000)),
	))
	require.NoError(t, err)
}

// we need to make this deterministic (same every test run), as content might affect gas costs
func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := ed25519.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}

func RandomAccountAddress() sdk.AccAddress {
	_, _, addr := keyPubAddr()
	return addr
}

func RandomBech32AccountAddress() string {
	return RandomAccountAddress().String()
}
