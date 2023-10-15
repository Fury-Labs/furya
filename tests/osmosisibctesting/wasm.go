package furyaibctesting

import (
	"crypto/sha256"
	"fmt"
	"os"

	"github.com/tidwall/gjson"

	"github.com/stretchr/testify/require"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	transfertypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
	"github.com/furya-labs/furya/v20/x/ibc-rate-limit/types"
	"github.com/stretchr/testify/suite"
)

func (chain *TestChain) StoreContractCode(suite *suite.Suite, path string) {
	furyaApp := chain.GetFuryaApp()

	govKeeper := furyaApp.GovKeeper
	wasmCode, err := os.ReadFile(path)
	suite.Require().NoError(err)

	addr := furyaApp.AccountKeeper.GetModuleAddress(govtypes.ModuleName)
	src := wasmtypes.StoreCodeProposalFixture(func(p *wasmtypes.StoreCodeProposal) {
		p.RunAs = addr.String()
		p.WASMByteCode = wasmCode
		checksum := sha256.Sum256(wasmCode)
		p.CodeHash = checksum[:]
	})

	// when stored
	storedProposal, err := govKeeper.SubmitProposal(chain.GetContext(), src, false)
	suite.Require().NoError(err)

	// and proposal execute
	handler := govKeeper.Router().GetRoute(storedProposal.ProposalRoute())
	err = handler(chain.GetContext(), storedProposal.GetContent())
	suite.Require().NoError(err)
}

func (chain *TestChain) InstantiateRLContract(suite *suite.Suite, quotas string) sdk.AccAddress {
	furyaApp := chain.GetFuryaApp()
	transferModule := furyaApp.AccountKeeper.GetModuleAddress(transfertypes.ModuleName)
	govModule := furyaApp.AccountKeeper.GetModuleAddress(govtypes.ModuleName)

	initMsgBz := []byte(fmt.Sprintf(`{
           "gov_module":  "%s",
           "ibc_module":"%s",
           "paths": [%s]
        }`,
		govModule, transferModule, quotas))

	contractKeeper := wasmkeeper.NewDefaultPermissionKeeper(furyaApp.WasmKeeper)
	codeID := uint64(1)
	creator := furyaApp.AccountKeeper.GetModuleAddress(govtypes.ModuleName)
	addr, _, err := contractKeeper.Instantiate(chain.GetContext(), codeID, creator, creator, initMsgBz, "rate limiting contract", nil)
	suite.Require().NoError(err)
	return addr
}

func (chain *TestChain) StoreContractCodeDirect(suite *suite.Suite, path string) uint64 {
	furyaApp := chain.GetFuryaApp()
	govKeeper := wasmkeeper.NewGovPermissionKeeper(furyaApp.WasmKeeper)
	creator := furyaApp.AccountKeeper.GetModuleAddress(govtypes.ModuleName)

	wasmCode, err := os.ReadFile(path)
	suite.Require().NoError(err)
	accessEveryone := wasmtypes.AccessConfig{Permission: wasmtypes.AccessTypeEverybody}
	codeID, _, err := govKeeper.Create(chain.GetContext(), creator, wasmCode, &accessEveryone)
	suite.Require().NoError(err)
	return codeID
}

func (chain *TestChain) InstantiateContract(suite *suite.Suite, msg string, codeID uint64) sdk.AccAddress {
	furyaApp := chain.GetFuryaApp()
	contractKeeper := wasmkeeper.NewDefaultPermissionKeeper(furyaApp.WasmKeeper)
	creator := furyaApp.AccountKeeper.GetModuleAddress(govtypes.ModuleName)
	addr, _, err := contractKeeper.Instantiate(chain.GetContext(), codeID, creator, creator, []byte(msg), "contract", nil)
	suite.Require().NoError(err)
	return addr
}

func (chain *TestChain) QueryContract(suite *suite.Suite, contract sdk.AccAddress, key []byte) string {
	furyaApp := chain.GetFuryaApp()
	state, err := furyaApp.WasmKeeper.QuerySmart(chain.GetContext(), contract, key)
	suite.Require().NoError(err)
	return string(state)
}

func (chain *TestChain) QueryContractJson(suite *suite.Suite, contract sdk.AccAddress, key []byte) gjson.Result {
	furyaApp := chain.GetFuryaApp()
	state, err := furyaApp.WasmKeeper.QuerySmart(chain.GetContext(), contract, key)
	suite.Require().NoError(err)
	suite.Require().True(gjson.Valid(string(state)))
	json := gjson.Parse(string(state))
	suite.Require().NoError(err)
	return json
}

func (chain *TestChain) RegisterRateLimitingContract(addr []byte) {
	addrStr, err := sdk.Bech32ifyAddressBytes("fury", addr)
	require.NoError(chain.T, err)
	params, err := types.NewParams(addrStr)
	require.NoError(chain.T, err)
	furyaApp := chain.GetFuryaApp()
	paramSpace, ok := furyaApp.AppKeepers.ParamsKeeper.GetSubspace(types.ModuleName)
	require.True(chain.T, ok)
	paramSpace.SetParamSet(chain.GetContext(), &params)
}
