package wasmbinding

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/fury-labs/furya/v20/app"
	"github.com/fury-labs/furya/v20/app/apptesting"
	"github.com/fury-labs/furya/v20/wasmbinding/bindings"
)

func SetupCustomApp(t *testing.T, addr sdk.AccAddress) (*app.FuryaApp, sdk.Context) {
	t.Helper()

	furya, ctx := CreateTestInput()
	wasmKeeper := furya.WasmKeeper

	storeReflectCode(t, ctx, furya, addr)

	cInfo := wasmKeeper.GetCodeInfo(ctx, 1)
	require.NotNil(t, cInfo)

	return furya, ctx
}

func TestQueryFullDenom(t *testing.T) {
	apptesting.SkipIfWSL(t)
	actor := RandomAccountAddress()
	furya, ctx := SetupCustomApp(t, actor)

	reflect := instantiateReflectContract(t, ctx, furya, actor)
	require.NotEmpty(t, reflect)

	// query full denom
	query := bindings.FuryaQuery{
		FullDenom: &bindings.FullDenom{
			CreatorAddr: reflect.String(),
			Subdenom:    "ustart",
		},
	}
	resp := bindings.FullDenomResponse{}
	queryCustom(t, ctx, furya, reflect, query, &resp)

	expected := fmt.Sprintf("factory/%s/ustart", reflect.String())
	require.EqualValues(t, expected, resp.Denom)
}

type ReflectQuery struct {
	Chain *ChainRequest `json:"chain,omitempty"`
}

type ChainRequest struct {
	Request wasmvmtypes.QueryRequest `json:"request"`
}

type ChainResponse struct {
	Data []byte `json:"data"`
}

func queryCustom(t *testing.T, ctx sdk.Context, furya *app.FuryaApp, contract sdk.AccAddress, request bindings.FuryaQuery, response interface{}) {
	t.Helper()

	msgBz, err := json.Marshal(request)
	require.NoError(t, err)

	query := ReflectQuery{
		Chain: &ChainRequest{
			Request: wasmvmtypes.QueryRequest{Custom: msgBz},
		},
	}
	queryBz, err := json.Marshal(query)
	require.NoError(t, err)

	resBz, err := furya.WasmKeeper.QuerySmart(ctx, contract, queryBz)
	require.NoError(t, err)
	var resp ChainResponse
	err = json.Unmarshal(resBz, &resp)
	require.NoError(t, err)
	err = json.Unmarshal(resp.Data, response)
	require.NoError(t, err)
}

func storeReflectCode(t *testing.T, ctx sdk.Context, furya *app.FuryaApp, addr sdk.AccAddress) {
	t.Helper()

	govKeeper := furya.GovKeeper
	wasmCode, err := os.ReadFile("../testdata/osmo_reflect.wasm")
	require.NoError(t, err)

	src := wasmtypes.StoreCodeProposalFixture(func(p *wasmtypes.StoreCodeProposal) {
		p.RunAs = addr.String()
		p.WASMByteCode = wasmCode
		checksum := sha256.Sum256(wasmCode)
		p.CodeHash = checksum[:]
	})

	// when stored
	storedProposal, err := govKeeper.SubmitProposal(ctx, src, false)
	require.NoError(t, err)

	// and proposal execute
	handler := govKeeper.Router().GetRoute(storedProposal.ProposalRoute())
	err = handler(ctx, storedProposal.GetContent())
	require.NoError(t, err)
}

func instantiateReflectContract(t *testing.T, ctx sdk.Context, furya *app.FuryaApp, funder sdk.AccAddress) sdk.AccAddress {
	t.Helper()

	initMsgBz := []byte("{}")
	contractKeeper := keeper.NewDefaultPermissionKeeper(furya.WasmKeeper)
	codeID := uint64(1)
	addr, _, err := contractKeeper.Instantiate(ctx, codeID, funder, funder, initMsgBz, "demo contract", nil)
	require.NoError(t, err)

	return addr
}

func fundAccount(t *testing.T, ctx sdk.Context, furya *app.FuryaApp, addr sdk.AccAddress, coins sdk.Coins) {
	t.Helper()
	err := simapp.FundAccount(
		furya.BankKeeper,
		ctx,
		addr,
		coins,
	)
	require.NoError(t, err)
}
