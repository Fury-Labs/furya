package wasmbinding

import (
	"fmt"
	"sync"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"

	gammv2types "github.com/fury-labs/furya/v20/x/gamm/v2types"

	concentratedliquidityquery "github.com/fury-labs/furya/v20/x/concentrated-liquidity/client/queryproto"
	cosmwasmpooltypes "github.com/fury-labs/furya/v20/x/cosmwasmpool/client/queryproto"
	downtimequerytypes "github.com/fury-labs/furya/v20/x/downtime-detector/client/queryproto"
	gammtypes "github.com/fury-labs/furya/v20/x/gamm/types"
	incentivestypes "github.com/fury-labs/furya/v20/x/incentives/types"
	lockuptypes "github.com/fury-labs/furya/v20/x/lockup/types"
	minttypes "github.com/fury-labs/furya/v20/x/mint/types"
	poolincentivestypes "github.com/fury-labs/furya/v20/x/pool-incentives/types"
	poolmanagerqueryproto "github.com/fury-labs/furya/v20/x/poolmanager/client/queryproto"
	superfluidtypes "github.com/fury-labs/furya/v20/x/superfluid/types"
	tokenfactorytypes "github.com/fury-labs/furya/v20/x/tokenfactory/types"
	twapquerytypes "github.com/fury-labs/furya/v20/x/twap/client/queryproto"
	txfeestypes "github.com/fury-labs/furya/v20/x/txfees/types"
	epochtypes "github.com/osmosis-labs/osmosis/x/epochs/types"
)

// stargateWhitelist keeps whitelist and its deterministic
// response binding for stargate queries.
// CONTRACT: since results of queries go into blocks, queries being added here should always be
// deterministic or can cause non-determinism in the state machine.
//
// The query can be multi-thread, so we have to use
// thread safe sync.Map.
var stargateWhitelist sync.Map

// Note: When adding a migration here, we should also add it to the Async ICQ params in the upgrade.
// In the future we may want to find a better way to keep these in sync

//nolint:staticcheck
func init() {
	// ibc queries
	setWhitelistedQuery("/ibc.applications.transfer.v1.Query/DenomTrace", &ibctransfertypes.QueryDenomTraceResponse{})

	// cosmos-sdk queries

	// auth
	setWhitelistedQuery("/cosmos.auth.v1beta1.Query/Account", &authtypes.QueryAccountResponse{})
	setWhitelistedQuery("/cosmos.auth.v1beta1.Query/Params", &authtypes.QueryParamsResponse{})
	setWhitelistedQuery("/cosmos.auth.v1beta1.Query/ModuleAccounts", &authtypes.QueryModuleAccountsResponse{})

	// bank
	setWhitelistedQuery("/cosmos.bank.v1beta1.Query/Balance", &banktypes.QueryBalanceResponse{})
	setWhitelistedQuery("/cosmos.bank.v1beta1.Query/DenomMetadata", &banktypes.QueryDenomsMetadataResponse{})
	setWhitelistedQuery("/cosmos.bank.v1beta1.Query/Params", &banktypes.QueryParamsResponse{})
	setWhitelistedQuery("/cosmos.bank.v1beta1.Query/SupplyOf", &banktypes.QuerySupplyOfResponse{})

	// distribution
	setWhitelistedQuery("/cosmos.distribution.v1beta1.Query/Params", &distributiontypes.QueryParamsResponse{})
	setWhitelistedQuery("/cosmos.distribution.v1beta1.Query/DelegatorWithdrawAddress", &distributiontypes.QueryDelegatorWithdrawAddressResponse{})
	setWhitelistedQuery("/cosmos.distribution.v1beta1.Query/ValidatorCommission", &distributiontypes.QueryValidatorCommissionResponse{})

	// gov
	setWhitelistedQuery("/cosmos.gov.v1beta1.Query/Deposit", &govtypes.QueryDepositResponse{})
	setWhitelistedQuery("/cosmos.gov.v1beta1.Query/Params", &govtypes.QueryParamsResponse{})
	setWhitelistedQuery("/cosmos.gov.v1beta1.Query/Vote", &govtypes.QueryVoteResponse{})

	// slashing
	setWhitelistedQuery("/cosmos.slashing.v1beta1.Query/Params", &slashingtypes.QueryParamsResponse{})
	setWhitelistedQuery("/cosmos.slashing.v1beta1.Query/SigningInfo", &slashingtypes.QuerySigningInfoResponse{})

	// staking
	setWhitelistedQuery("/cosmos.staking.v1beta1.Query/Delegation", &stakingtypes.QueryDelegationResponse{})
	setWhitelistedQuery("/cosmos.staking.v1beta1.Query/Params", &stakingtypes.QueryParamsResponse{})
	setWhitelistedQuery("/cosmos.staking.v1beta1.Query/Validator", &stakingtypes.QueryValidatorResponse{})

	// furya queries
	// cosmwasm pool
	setWhitelistedQuery("/furya.cosmwasmpool.v1beta1.Query/Pools", &cosmwasmpooltypes.PoolsResponse{})
	setWhitelistedQuery("/furya.cosmwasmpool.v1beta1.Query/Params", &cosmwasmpooltypes.ParamsResponse{})
	setWhitelistedQuery("/furya.cosmwasmpool.v1beta1.Query/ContractInfoByPoolId", &cosmwasmpooltypes.ContractInfoByPoolIdResponse{})

	// epochs
	setWhitelistedQuery("/furya.epochs.v1beta1.Query/EpochInfos", &epochtypes.QueryEpochsInfoResponse{})
	setWhitelistedQuery("/furya.epochs.v1beta1.Query/CurrentEpoch", &epochtypes.QueryCurrentEpochResponse{})

	// gamm
	setWhitelistedQuery("/furya.gamm.v1beta1.Query/NumPools", &gammtypes.QueryNumPoolsResponse{}) // ==> use x/poolmanager
	setWhitelistedQuery("/furya.gamm.v1beta1.Query/TotalLiquidity", &gammtypes.QueryTotalLiquidityResponse{})
	setWhitelistedQuery("/furya.gamm.v1beta1.Query/Pool", &gammtypes.QueryPoolResponse{}) // ==> use x/poolmanager
	setWhitelistedQuery("/furya.gamm.v1beta1.Query/PoolParams", &gammtypes.QueryPoolParamsResponse{})
	setWhitelistedQuery("/furya.gamm.v1beta1.Query/TotalPoolLiquidity", &gammtypes.QueryTotalPoolLiquidityResponse{}) // ==> use x/poolmanager
	setWhitelistedQuery("/furya.gamm.v1beta1.Query/TotalShares", &gammtypes.QueryTotalSharesResponse{})
	setWhitelistedQuery("/furya.gamm.v1beta1.Query/CalcJoinPoolShares", &gammtypes.QueryCalcJoinPoolSharesResponse{})
	setWhitelistedQuery("/furya.gamm.v1beta1.Query/CalcExitPoolCoinsFromShares", &gammtypes.QueryCalcExitPoolCoinsFromSharesResponse{})
	setWhitelistedQuery("/furya.gamm.v1beta1.Query/CalcJoinPoolNoSwapShares", &gammtypes.QueryCalcJoinPoolNoSwapSharesResponse{})
	setWhitelistedQuery("/furya.gamm.v1beta1.Query/PoolType", &gammtypes.QueryPoolTypeResponse{})
	setWhitelistedQuery("/furya.gamm.v2.Query/SpotPrice", &gammv2types.QuerySpotPriceResponse{})
	setWhitelistedQuery("/furya.gamm.v1beta1.Query/EstimateSwapExactAmountIn", &gammtypes.QuerySwapExactAmountInResponse{})   // ==> use x/poolmanager
	setWhitelistedQuery("/furya.gamm.v1beta1.Query/EstimateSwapExactAmountOut", &gammtypes.QuerySwapExactAmountOutResponse{}) // ==> use x/poolmanager

	// incentives
	setWhitelistedQuery("/furya.incentives.Query/ModuleToDistributeCoins", &incentivestypes.ModuleToDistributeCoinsResponse{})
	setWhitelistedQuery("/furya.incentives.Query/LockableDurations", &incentivestypes.QueryLockableDurationsResponse{})

	// lockup
	setWhitelistedQuery("/furya.lockup.Query/ModuleBalance", &lockuptypes.ModuleBalanceResponse{})
	setWhitelistedQuery("/furya.lockup.Query/ModuleLockedAmount", &lockuptypes.ModuleLockedAmountResponse{})
	// Warning: it iterates over every single lock account has, which means this query can have unbounded gas
	setWhitelistedQuery("/furya.lockup.Query/AccountLockedCoins", &lockuptypes.AccountLockedCoinsResponse{})
	setWhitelistedQuery("/furya.lockup.Query/AccountUnlockableCoins", &lockuptypes.AccountUnlockableCoinsResponse{})
	setWhitelistedQuery("/furya.lockup.Query/AccountUnlockingCoins", &lockuptypes.AccountUnlockingCoinsResponse{})
	setWhitelistedQuery("/furya.lockup.Query/LockedDenom", &lockuptypes.LockedDenomResponse{})
	setWhitelistedQuery("/furya.lockup.Query/LockedByID", &lockuptypes.LockedResponse{})
	setWhitelistedQuery("/furya.lockup.Query/NextLockID", &lockuptypes.NextLockIDResponse{})
	setWhitelistedQuery("/furya.lockup.Query/LockRewardReceiver", &lockuptypes.LockRewardReceiverResponse{})

	// mint
	setWhitelistedQuery("/furya.mint.v1beta1.Query/EpochProvisions", &minttypes.QueryEpochProvisionsResponse{})
	setWhitelistedQuery("/furya.mint.v1beta1.Query/Params", &minttypes.QueryParamsResponse{})

	// pool-incentives
	setWhitelistedQuery("/furya.poolincentives.v1beta1.Query/GaugeIds", &poolincentivestypes.QueryGaugeIdsResponse{})

	// superfluid
	setWhitelistedQuery("/furya.superfluid.Query/Params", &superfluidtypes.QueryParamsResponse{})
	setWhitelistedQuery("/furya.superfluid.Query/AssetType", &superfluidtypes.AssetTypeResponse{})
	setWhitelistedQuery("/furya.superfluid.Query/AllAssets", &superfluidtypes.AllAssetsResponse{})
	setWhitelistedQuery("/furya.superfluid.Query/AssetMultiplier", &superfluidtypes.AssetMultiplierResponse{})

	// poolmanager
	setWhitelistedQuery("/furya.poolmanager.v1beta1.Query/NumPools", &poolmanagerqueryproto.NumPoolsResponse{})
	setWhitelistedQuery("/furya.poolmanager.v1beta1.Query/EstimateSwapExactAmountIn", &poolmanagerqueryproto.EstimateSwapExactAmountInResponse{})
	setWhitelistedQuery("/furya.poolmanager.v1beta1.Query/EstimateSwapExactAmountOut", &poolmanagerqueryproto.EstimateSwapExactAmountOutResponse{})
	setWhitelistedQuery("/furya.poolmanager.v1beta1.Query/EstimateSinglePoolSwapExactAmountIn", &poolmanagerqueryproto.EstimateSwapExactAmountInResponse{})
	setWhitelistedQuery("/furya.poolmanager.v1beta1.Query/EstimateSinglePoolSwapExactAmountOut", &poolmanagerqueryproto.EstimateSwapExactAmountOutResponse{})
	setWhitelistedQuery("/furya.poolmanager.v1beta1.Query/Pool", &poolmanagerqueryproto.PoolResponse{})
	setWhitelistedQuery("/furya.poolmanager.v1beta1.Query/SpotPrice", &poolmanagerqueryproto.SpotPriceResponse{})
	setWhitelistedQuery("/furya.poolmanager.v1beta1.Query/TotalPoolLiquidity", &poolmanagerqueryproto.TotalPoolLiquidityResponse{})
	setWhitelistedQuery("/furya.poolmanager.v1beta1.Query/Params", &poolmanagerqueryproto.ParamsResponse{})
	setWhitelistedQuery("/furya.poolmanager.v1beta1.Query/TradingPairTakerFee", &poolmanagerqueryproto.TradingPairTakerFeeResponse{})

	// txfees
	setWhitelistedQuery("/furya.txfees.v1beta1.Query/FeeTokens", &txfeestypes.QueryFeeTokensResponse{})
	setWhitelistedQuery("/furya.txfees.v1beta1.Query/DenomSpotPrice", &txfeestypes.QueryDenomSpotPriceResponse{})
	setWhitelistedQuery("/furya.txfees.v1beta1.Query/DenomPoolId", &txfeestypes.QueryDenomPoolIdResponse{})
	setWhitelistedQuery("/furya.txfees.v1beta1.Query/BaseDenom", &txfeestypes.QueryBaseDenomResponse{})

	// tokenfactory
	setWhitelistedQuery("/furya.tokenfactory.v1beta1.Query/Params", &tokenfactorytypes.QueryParamsResponse{})
	setWhitelistedQuery("/furya.tokenfactory.v1beta1.Query/DenomAuthorityMetadata", &tokenfactorytypes.QueryDenomAuthorityMetadataResponse{})
	// Does not include denoms_from_creator, TBD if this is the index we want contracts to use instead of admin

	// twap
	setWhitelistedQuery("/furya.twap.v1beta1.Query/ArithmeticTwap", &twapquerytypes.ArithmeticTwapResponse{})
	setWhitelistedQuery("/furya.twap.v1beta1.Query/ArithmeticTwapToNow", &twapquerytypes.ArithmeticTwapToNowResponse{})
	setWhitelistedQuery("/furya.twap.v1beta1.Query/GeometricTwap", &twapquerytypes.GeometricTwapResponse{})
	setWhitelistedQuery("/furya.twap.v1beta1.Query/GeometricTwapToNow", &twapquerytypes.GeometricTwapToNowResponse{})
	setWhitelistedQuery("/furya.twap.v1beta1.Query/Params", &twapquerytypes.ParamsResponse{})

	// downtime-detector
	setWhitelistedQuery("/furya.downtimedetector.v1beta1.Query/RecoveredSinceDowntimeOfLength", &downtimequerytypes.RecoveredSinceDowntimeOfLengthResponse{})

	// concentrated-liquidity
	setWhitelistedQuery("/furya.concentratedliquidity.v1beta1.Query/UserPositions", &concentratedliquidityquery.UserPositionsResponse{})
	setWhitelistedQuery("/furya.concentratedliquidity.v1beta1.Query/LiquidityPerTickRange", &concentratedliquidityquery.LiquidityPerTickRangeResponse{})
	setWhitelistedQuery("/furya.concentratedliquidity.v1beta1.Query/ClaimableSpreadRewards", &concentratedliquidityquery.ClaimableSpreadRewardsResponse{})
	setWhitelistedQuery("/furya.concentratedliquidity.v1beta1.Query/ClaimableIncentives", &concentratedliquidityquery.ClaimableIncentivesResponse{})
	setWhitelistedQuery("/furya.concentratedliquidity.v1beta1.Query/PositionById", &concentratedliquidityquery.PositionByIdResponse{})
	setWhitelistedQuery("/furya.concentratedliquidity.v1beta1.Query/Params", &concentratedliquidityquery.ParamsResponse{})
	setWhitelistedQuery("/furya.concentratedliquidity.v1beta1.Query/PoolAccumulatorRewards", &concentratedliquidityquery.PoolAccumulatorRewardsResponse{})
	setWhitelistedQuery("/furya.concentratedliquidity.v1beta1.Query/IncentiveRecords", &concentratedliquidityquery.IncentiveRecordsResponse{})
	setWhitelistedQuery("/furya.concentratedliquidity.v1beta1.Query/TickAccumulatorTrackers", &concentratedliquidityquery.TickAccumulatorTrackersResponse{})
	setWhitelistedQuery("/furya.concentratedliquidity.v1beta1.Query/CFMMPoolIdLinkFromConcentratedPoolId", &concentratedliquidityquery.CFMMPoolIdLinkFromConcentratedPoolIdResponse{})
}

// GetWhitelistedQuery returns the whitelisted query at the provided path.
// If the query does not exist, or it was setup wrong by the chain, this returns an error.
func GetWhitelistedQuery(queryPath string) (codec.ProtoMarshaler, error) {
	protoResponseAny, isWhitelisted := stargateWhitelist.Load(queryPath)
	if !isWhitelisted {
		return nil, wasmvmtypes.UnsupportedRequest{Kind: fmt.Sprintf("'%s' path is not allowed from the contract", queryPath)}
	}
	protoResponseType, ok := protoResponseAny.(codec.ProtoMarshaler)
	if !ok {
		return nil, wasmvmtypes.Unknown{}
	}
	return protoResponseType, nil
}

func setWhitelistedQuery(queryPath string, protoType codec.ProtoMarshaler) {
	stargateWhitelist.Store(queryPath, protoType)
}

func GetStargateWhitelistedPaths() (keys []string) {
	// Iterate over the map and collect the keys
	stargateWhitelist.Range(func(key, value interface{}) bool {
		keyStr, ok := key.(string)
		if !ok {
			panic("key is not a string")
		}
		keys = append(keys, keyStr)
		return true
	})

	return keys
}
